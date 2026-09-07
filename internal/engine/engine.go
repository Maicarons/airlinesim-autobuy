// Package engine orchestrates the monitoring pipeline for multiple servers.
package engine

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/Maicarons/airlinesim-autobuy/internal/auth"
	"github.com/Maicarons/airlinesim-autobuy/internal/client"
	"github.com/Maicarons/airlinesim-autobuy/internal/collector"
	"github.com/Maicarons/airlinesim-autobuy/internal/config"
	"github.com/Maicarons/airlinesim-autobuy/internal/executor"
	"github.com/Maicarons/airlinesim-autobuy/internal/notifier"
	"github.com/Maicarons/airlinesim-autobuy/internal/parser"
	"github.com/Maicarons/airlinesim-autobuy/internal/rules"
)

// ServerStatus tracks per-server monitoring state.
type ServerStatus struct {
	Host        string    `json:"host"`
	ScanCount   int64     `json:"scan_count"`
	FoundCount  int64     `json:"found_count"`
	BoughtCount int64     `json:"bought_count"`
	FailedCount int64     `json:"failed_count"`
	LastScan    time.Time `json:"last_scan,omitempty"`
	LastError   string    `json:"last_error,omitempty"`
}

// Status represents the current engine status.
type Status struct {
	Running    bool           `json:"running"`
	StartTime  time.Time      `json:"start_time,omitempty"`
	Servers    []ServerStatus `json:"servers"`
	ScanCount  int64          `json:"scan_count"`
	FoundCount int64          `json:"found_count"`
	BoughtCount int64         `json:"bought_count"`
	FailedCount int64         `json:"failed_count"`
	LastError  string         `json:"last_error,omitempty"`
}

// serverRunner manages monitoring for a single game server.
type serverRunner struct {
	host           string
	hostBaseURL    string
	companies      []string
	defaultCompany string
	session        *auth.Session
	collector      *collector.Collector
	httpClient     *client.Client
	executor       *executor.Executor
	status         ServerStatus
}

// Engine orchestrates the entire monitoring pipeline.
type Engine struct {
	cfgStore  *config.Store
	notifier  *notifier.Notifier
	parser    *parser.Parser
	rulesEng  *rules.Engine
	servers   []*serverRunner

	mu                 sync.RWMutex
	status             Status
	cancel             context.CancelFunc
	seen               map[string]bool
	rulePurchaseCount  map[string]int // tracks purchases per rule name for MaxCount
}

// New creates a new monitoring engine for all configured servers.
func New(
	cfgStore *config.Store,
	notif *notifier.Notifier,
) *Engine {
	cfg := cfgStore.Get()
	eng := &Engine{
		cfgStore: cfgStore,
		notifier: notif,
		parser:   parser.New(),
		rulesEng: rules.New(cfgStore.GetRules()),
		seen:     make(map[string]bool),
		rulePurchaseCount: make(map[string]int),
		status:   Status{Servers: make([]ServerStatus, 0)},
	}

// Create a runner for each server that has rules targeting it
		for _, sv := range cfg.Servers {
			// Check if any enabled rule targets this server
			hasRule := false
			for _, rule := range cfg.Rules {
				if rule.Enabled && (rule.ServerID < 0 || (rule.ServerID < len(cfg.Servers) && cfg.Servers[rule.ServerID].BaseURL == sv.BaseURL)) {
					hasRule = true
					break
				}
			}
			if !hasRule {
				slog.Info("skipping server, no rules target it", "host", sv.Host)
				continue
			}

			// Find which auth to use for this server
		// Look for the first rule targeting this server
		authIdx := 0 // default to first auth
		for _, rule := range cfg.Rules {
			if rule.ServerID >= 0 && rule.ServerID < len(cfg.Servers) {
				svDefault := cfg.Servers[rule.ServerID]
				if svDefault.BaseURL == sv.BaseURL {
					if rule.AuthID >= 0 && rule.AuthID < len(cfg.Auths) {
						authIdx = rule.AuthID
					}
					break
				}
			}
		}
		// Ensure auth index is valid
		if authIdx >= len(cfg.Auths) {
			authIdx = 0
		}

		au := cfg.Auths[authIdx]
		sessionFile := au.SessionFile
		if sessionFile == "" {
			sessionFile = fmt.Sprintf("session_%s.json", sv.Host)
		}
		session, err := auth.New(sv.BaseURL, au.Username, au.Password, sessionFile)
		if err != nil {
			slog.Warn("failed to create session for server", "host", sv.Host, "error", err)
			continue
		}
		cl := client.New(session.Client(), cfg.Monitor.Duration(), time.Duration(cfg.Monitor.Jitter)*time.Second)
		col := collector.New(cl.InnerClient(), sv.BaseURL)
		exec := executor.New(session.Client())

eng.servers = append(eng.servers, &serverRunner{
				host:           sv.Host,
				hostBaseURL:    sv.BaseURL,
				companies:      sv.Companies,
				defaultCompany: "",
				session:        session,
				collector:      col,
				httpClient:     cl,
				executor:       exec,
				status:         ServerStatus{Host: sv.Host},
			})
		eng.status.Servers = append(eng.status.Servers, ServerStatus{Host: sv.Host})
	}

	slog.Info("engine initialized", "servers", len(eng.servers))
	return eng
}

// Start begins monitoring all servers.
func (e *Engine) Start() error {
	e.mu.Lock()
	if e.status.Running {
		e.mu.Unlock()
		return fmt.Errorf("engine is already running")
	}
	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel
	e.status = Status{
		Running:   true,
		StartTime: time.Now(),
		Servers:   make([]ServerStatus, len(e.servers)),
	}
	for i, sr := range e.servers {
		e.status.Servers[i] = ServerStatus{Host: sr.host}
	}
	e.mu.Unlock()

	// Build startup notification
	cfg := e.cfgStore.Get()
	enabledRules := 0
	for _, r := range cfg.Rules {
		if r.Enabled {
			enabledRules++
		}
	}
	serverList := make([]string, len(e.servers))
	for i, sr := range e.servers {
		serverList[i] = sr.host
	}

	startupMsg := fmt.Sprintf("🚀 Engine started — %d server(s): %s | %d/%d rule(s) enabled",
		len(e.servers),
		strings.Join(serverList, ", "),
		enabledRules,
		len(cfg.Rules))
	e.notifier.NotifyInfo(startupMsg)

	slog.Info("engine starting", "servers", len(e.servers))

	// Start a goroutine for each server
	for _, sr := range e.servers {
		go e.runServerLoop(ctx, sr)
	}
	return nil
}

// Stop gracefully stops monitoring all servers.
func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.status.Running {
		return
	}
	if e.cancel != nil {
		e.cancel()
	}
	e.status.Running = false
	slog.Info("engine stopped")
}

// Status returns the current engine status.
func (e *Engine) Status() Status {
	e.mu.RLock()
	defer e.mu.RUnlock()
	s := e.status
	// Aggregate per-server stats
	s.ScanCount = 0
	s.FoundCount = 0
	s.BoughtCount = 0
	s.FailedCount = 0
	for i, sr := range e.servers {
		e.status.Servers[i] = sr.status
		s.ScanCount += sr.status.ScanCount
		s.FoundCount += sr.status.FoundCount
		s.BoughtCount += sr.status.BoughtCount
		s.FailedCount += sr.status.FailedCount
		if sr.status.LastError != "" {
			s.LastError = sr.status.LastError
		}
	}
	return s
}

// ReloadRules reloads rules from the config store.
func (e *Engine) ReloadRules() {
	cfg := e.cfgStore.Get()
	e.rulesEng.UpdateRules(cfg.Rules)
	slog.Info("rules reloaded", "count", len(cfg.Rules))
}

// GetParser returns the engine's parser for use by the web UI.
func (e *Engine) GetParser() *parser.Parser {
	return e.parser
}

// GetRulesEngine returns the engine's rules engine for use by the web UI.
func (e *Engine) GetRulesEngine() *rules.Engine {
	return e.rulesEng
}

// runServerLoop runs the monitoring loop for a single server.
func (e *Engine) runServerLoop(ctx context.Context, sr *serverRunner) {
	cfg := e.cfgStore.Get()
	ticker := time.NewTicker(cfg.Monitor.Duration())
	defer ticker.Stop()

	slog.Info("starting server monitor", "host", sr.host)

// Login
		if err := sr.session.Login(); err != nil {
			slog.Error("server login failed", "host", sr.host, "error", err)
			sr.status.LastError = err.Error()
			e.mu.Lock()
			e.status.LastError = err.Error()
			e.mu.Unlock()
			return
		}

// Select default company
			if len(sr.companies) > 0 {
				sr.defaultCompany = sr.companies[0]
				if err := sr.session.SelectCompany(sr.defaultCompany); err != nil {
					slog.Error("company selection failed", "host", sr.host, "company", sr.defaultCompany, "error", err)
					sr.status.LastError = err.Error()
					e.mu.Lock()
					e.status.LastError = err.Error()
					e.mu.Unlock()
					return
				}
			}

		// Discover market URL
	if err := sr.collector.DiscoverMarketURL(); err != nil {
		slog.Warn("market URL discovery failed", "host", sr.host, "error", err)
	}

	// Run immediately, then on ticker
	e.runServerCycle(sr)

	for {
		select {
		case <-ctx.Done():
			slog.Info("server monitor stopped", "host", sr.host)
			return
		case <-ticker.C:
			e.runServerCycle(sr)
		}
	}
}

// runServerCycle performs a single monitoring cycle for one server.
func (e *Engine) runServerCycle(sr *serverRunner) {
	slog.Debug("running server cycle", "host", sr.host)

	// Check session health
	if !sr.session.HealthCheck() {
		slog.Info("session expired, re-authenticating", "host", sr.host)
		if err := sr.session.Login(); err != nil {
			e.notifier.NotifyError(fmt.Errorf("re-authentication failed for %s: %w", sr.host, err))
			sr.status.LastError = err.Error()
			return
		}
	}

// Build filter params from the first enabled rule
		cfg := e.cfgStore.Get()
		var filterParams *collector.FilterParams
		for _, rule := range cfg.Rules {
			if !rule.Enabled {
				continue
			}
			// Match by server_id: 0-based index into Servers, or -1 for all
			if rule.ServerID >= 0 && rule.ServerID < len(cfg.Servers) {
				svConfig := cfg.Servers[rule.ServerID]
				// Match by BaseURL to handle skipped servers
				if svConfig.BaseURL != sr.hostBaseURL {
					continue
				}
			}
			filterParams = &collector.FilterParams{
				FamilyID: rule.Match.FamilyID,
				TypeID:   rule.Match.TypeID,
				SortBy:   rule.Match.SortBy,
			}
			break
		}

	// Fetch market page
	page, err := sr.collector.Fetch(filterParams)
	if err != nil {
		slog.Warn("failed to fetch market page", "host", sr.host, "error", err)
		sr.status.LastError = err.Error()
		return
	}

	// Parse the page
	parseResult, err := e.parser.Parse(page.Body)
	if err != nil {
		slog.Warn("failed to parse market page", "host", sr.host, "error", err)
		return
	}

	// Update stats
	sr.status.ScanCount++
	sr.status.LastScan = page.FetchedAt

	// Process each aircraft
	for i := range parseResult.Offers {
		ac := &parseResult.Offers[i]
		ac.SeenAt = page.FetchedAt

		id := ac.GenerateID()
		e.mu.Lock()
		if e.seen[id] {
			e.mu.Unlock()
			continue
		}
		e.seen[id] = true
		e.mu.Unlock()

		matches := e.rulesEng.Evaluate(ac)
		if len(matches) == 0 {
			continue
		}
		sr.status.FoundCount++

		best := matches[0]
		e.notifier.NotifyAircraftFound(ac, best.Rule.Name)

if best.ShouldBuy {
					// Check MaxCount limit
					if best.Rule.Action.MaxCount > 0 {
						e.mu.Lock()
						count := e.rulePurchaseCount[best.Rule.Name]
						e.mu.Unlock()
						if count >= best.Rule.Action.MaxCount {
							slog.Debug("purchase limit reached for rule",
								"rule", best.Rule.Name,
								"max_count", best.Rule.Action.MaxCount,
							)
							continue
						}
					}

					// Switch to the rule's company if specified and different from current
					ruleCompany := best.Rule.CompanyName
					switchedCompany := false
					if ruleCompany != "" && ruleCompany != sr.session.CompanyName() {
						slog.Info("switching company for purchase",
							"rule", best.Rule.Name,
							"from", sr.session.CompanyName(),
							"to", ruleCompany,
						)
						if err := sr.session.SelectCompany(ruleCompany); err != nil {
							slog.Error("failed to switch company for purchase", "error", err)
							e.notifier.NotifyPurchaseFailed(ac.Type, fmt.Sprintf("company switch failed: %v", err))
							sr.status.FailedCount++
							continue
						}
						switchedCompany = true
					}

					req := &executor.PurchaseRequest{
						AircraftURL: ac.URL,
						Price:       ac.Price,
						RuleName:    best.Rule.Name,
						IsAuction:   ac.OfferType == "auction",
						BidAmount:   best.Rule.Action.MaxBidIncrement,
						MinBalance:  cfg.Monitor.MinBalance,
					}
					result := sr.executor.Execute(req)
					if result.Success {
						e.mu.Lock()
						e.rulePurchaseCount[best.Rule.Name]++
						e.mu.Unlock()
						sr.status.BoughtCount++
						e.notifier.NotifyPurchaseMade(ac.Type, ac.Price, best.Rule.Name)
					} else {
						sr.status.FailedCount++
						e.notifier.NotifyPurchaseFailed(ac.Type, result.Message)
					}

					// Switch back to default company after purchase
					if switchedCompany && sr.defaultCompany != "" {
						if err := sr.session.SelectCompany(sr.defaultCompany); err != nil {
							slog.Warn("failed to switch back to default company", "company", sr.defaultCompany, "error", err)
						}
					}
				}
	}
	slog.Debug("server cycle complete", "host", sr.host, "offers", len(parseResult.Offers))
}
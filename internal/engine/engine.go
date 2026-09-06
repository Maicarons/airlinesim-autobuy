// Package engine orchestrates the monitoring pipeline for multiple servers.
package engine

import (
	"context"
	"fmt"
	"log/slog"
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
	host      string
	session   *auth.Session
	collector *collector.Collector
	httpClient *client.Client
	status    ServerStatus
}

// Engine orchestrates the entire monitoring pipeline.
type Engine struct {
	cfgStore  *config.Store
	notifier  *notifier.Notifier
	parser    *parser.Parser
	rulesEng  *rules.Engine
	executor  *executor.Executor
	servers   []*serverRunner

	mu       sync.RWMutex
	status   Status
	cancel   context.CancelFunc
	seen     map[string]bool
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
		executor: executor.New(nil),
		seen:     make(map[string]bool),
		status:   Status{Servers: make([]ServerStatus, 0)},
	}

	// Create a runner for each server
	for _, sv := range cfg.Servers {
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
		session, err := auth.New(sv.BaseURL, au.Username, au.Password, cfg.Auths[0].SessionFile)
		if err != nil {
			slog.Warn("failed to create session for server", "host", sv.Host, "error", err)
			continue
		}
		cl := client.New(session.Client(), cfg.Monitor.Duration(), time.Duration(cfg.Monitor.Jitter)*time.Second)
		col := collector.New(cl.InnerClient(), sv.BaseURL)

		eng.servers = append(eng.servers, &serverRunner{
			host:       sv.Host,
			session:    session,
			collector:  col,
			httpClient: cl,
			status:     ServerStatus{Host: sv.Host},
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
		if rule.Enabled && (rule.ServerID < 0 || rule.ServerID < len(e.servers) && e.servers[rule.ServerID] == sr) {
			filterParams = &collector.FilterParams{
				FamilyID: rule.Match.FamilyID,
				TypeID:   rule.Match.TypeID,
				SortBy:   rule.Match.SortBy,
			}
			break
		}
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
			req := &executor.PurchaseRequest{
				AircraftURL: ac.URL,
				Price:       ac.Price,
				RuleName:    best.Rule.Name,
				IsAuction:   ac.OfferType == "auction",
				BidAmount:   best.Rule.Action.MaxBidIncrement,
				MinBalance:  cfg.Monitor.MinBalance,
			}
			result := e.executor.Execute(req)
			if result.Success {
				sr.status.BoughtCount++
				e.notifier.NotifyPurchaseMade(ac.Type, ac.Price, best.Rule.Name)
			} else {
				sr.status.FailedCount++
				e.notifier.NotifyPurchaseFailed(ac.Type, result.Message)
			}
		}
	}
	slog.Debug("server cycle complete", "host", sr.host, "offers", len(parseResult.Offers))
}
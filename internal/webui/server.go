// Package webui provides the embedded web management interface.
package webui

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/Maicarons/airlinesim-autobuy/internal/config"
	"github.com/Maicarons/airlinesim-autobuy/internal/engine"
	"github.com/Maicarons/airlinesim-autobuy/internal/marketdata"
	"github.com/Maicarons/airlinesim-autobuy/internal/notifier"
)

//go:embed frontend/dist/*
var frontendFiles embed.FS

// Server is the HTTP server for the web management interface.
type Server struct {
	cfgStore *config.Store
	engine   *engine.Engine
	notifier *notifier.Notifier
	host     string
	port     int
	router   *chi.Mux
}

// New creates a new WebUI server.
func New(cfgStore *config.Store, eng *engine.Engine, notif *notifier.Notifier, host string, port int) *Server {
	s := &Server{
		cfgStore: cfgStore,
		engine:   eng,
		notifier: notif,
		host:     host,
		port:     port,
		router:   chi.NewRouter(),
	}

	s.setupRoutes()
	return s
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	slog.Info("web UI listening", "addr", addr)

	server := &http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return server.ListenAndServe()
}

// Stop gracefully stops the server.
func (s *Server) Stop() {}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(corsMiddleware)

	// API routes
	s.router.Route("/api", func(r chi.Router) {
		// Status
		r.Get("/status", s.handleStatus)

		// Control
		r.Post("/control/start", s.handleStart)
		r.Post("/control/stop", s.handleStop)

		// Reload rules/config without restart
		r.Post("/reload", s.handleReload)

		// Config
		r.Get("/config", s.handleGetConfig)
		r.Put("/config", s.handleUpdateConfig)

		// Rules
		r.Get("/rules", s.handleGetRules)
		r.Post("/rules", s.handleCreateRule)
		r.Get("/rules/{id}", s.handleGetRule)
		r.Put("/rules/{id}", s.handleUpdateRule)
		r.Delete("/rules/{id}", s.handleDeleteRule)
		r.Patch("/rules/{id}/toggle", s.handleToggleRule)
		r.Put("/rules/reorder", s.handleReorderRules)

		// Aircraft market data
		r.Get("/aircraft-data", s.handleAircraftData)
	})

	// Serve embedded Vue SPA — all non-API routes go to index.html
	spaFS := s.getSPAFileSystem()
	fileServer := http.FileServer(spaFS)

	s.router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Strip leading slash
		cleanPath := strings.TrimPrefix(r.URL.Path, "/")

		// Try to serve the exact file
		if _, err := spaFS.Open(cleanPath); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for Vue Router routes
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}

// getSPAFileSystem returns an http.FileSystem for the embedded frontend dist,
// with the frontend/dist/ prefix stripped.
func (s *Server) getSPAFileSystem() http.FileSystem {
	sub, err := fs.Sub(frontendFiles, "frontend/dist")
	if err != nil {
		panic(fmt.Errorf("failed to get frontend dist sub-filesystem: %w", err))
	}
	return http.FS(sub)
}

// handleReload reloads rules from config without restarting the engine.
func (s *Server) handleReload(w http.ResponseWriter, r *http.Request) {
	s.engine.ReloadRules()
	writeJSON(w, http.StatusOK, map[string]string{"status": "reloaded"})
}

// handleAircraftData returns the aircraft families and types for the UI dropdowns.
func (s *Server) handleAircraftData(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, marketdata.AircraftData)
}

// ---- API Handlers ----

// handleStatus returns the current engine status.
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.engine.Status())
}

// handleStart starts the monitoring engine.
func (s *Server) handleStart(w http.ResponseWriter, r *http.Request) {
	if err := s.engine.Start(); err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "started"})
}

// handleStop stops the monitoring engine.
func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	s.engine.Stop()
	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}

// handleGetConfig returns the full configuration.
func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.cfgStore.Get())
}

// handleUpdateConfig updates the global configuration (excluding rules).
func (s *Server) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var incoming config.Config
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

// Get current config to preserve fields not sent by the client
		cfg := s.cfgStore.Get()

		// Update server config (multi-server)
		if len(incoming.Servers) > 0 {
			cfg.Servers = incoming.Servers
		} else if incoming.Server.BaseURL != "" {
			// Legacy single-server fallback
			cfg.Servers = []config.ServerConfig{incoming.Server}
		}

// Update auth config (multi-auth)
		if len(incoming.Auths) > 0 {
			cfg.Auths = incoming.Auths
		} else if incoming.Auth.Username != "" {
			cfg.Auths = []config.AuthConfig{incoming.Auth}
		}

// Update monitor config
		if incoming.Monitor.Interval > 0 {
			cfg.Monitor.Interval = incoming.Monitor.Interval
		}
		if incoming.Monitor.Jitter > 0 {
			cfg.Monitor.Jitter = incoming.Monitor.Jitter
		}
		if incoming.Monitor.RequestTimeout > 0 {
			cfg.Monitor.RequestTimeout = incoming.Monitor.RequestTimeout
		}
		if incoming.Monitor.MinBalance > 0 {
			cfg.Monitor.MinBalance = incoming.Monitor.MinBalance
		}

	// Update notifier config (bool fields need special handling)
	cfg.Notifier.Console = incoming.Notifier.Console
	if incoming.Notifier.DiscordWebhook != "" {
		cfg.Notifier.DiscordWebhook = incoming.Notifier.DiscordWebhook
	}

	// Update webui config
	if incoming.WebUI.Host != "" {
		cfg.WebUI.Host = incoming.WebUI.Host
	}
	if incoming.WebUI.Port > 0 {
		cfg.WebUI.Port = incoming.WebUI.Port
	}

	if err := s.cfgStore.Update(cfg); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, cfg)
}

// handleGetRules returns the list of rules.
func (s *Server) handleGetRules(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.cfgStore.GetRules())
}

// handleCreateRule creates a new rule.
func (s *Server) handleCreateRule(w http.ResponseWriter, r *http.Request) {
	var rule config.RuleConfig
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	rules := s.cfgStore.GetRules()
	rules = append(rules, rule)

	if err := s.cfgStore.UpdateRules(rules); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	s.engine.ReloadRules()
	writeJSON(w, http.StatusCreated, rule)
}

// handleGetRule returns a single rule by index.
func (s *Server) handleGetRule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid rule ID"})
		return
	}

	rules := s.cfgStore.GetRules()
	if id < 0 || id >= len(rules) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "rule not found"})
		return
	}

	writeJSON(w, http.StatusOK, rules[id])
}

// handleUpdateRule updates a rule by index.
func (s *Server) handleUpdateRule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid rule ID"})
		return
	}

	var rule config.RuleConfig
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	rules := s.cfgStore.GetRules()
	if id < 0 || id >= len(rules) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "rule not found"})
		return
	}

	rules[id] = rule
	if err := s.cfgStore.UpdateRules(rules); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	s.engine.ReloadRules()
	writeJSON(w, http.StatusOK, rule)
}

// handleDeleteRule deletes a rule by index.
func (s *Server) handleDeleteRule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid rule ID"})
		return
	}

	rules := s.cfgStore.GetRules()
	if id < 0 || id >= len(rules) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "rule not found"})
		return
	}

	rules = append(rules[:id], rules[id+1:]...)
	if err := s.cfgStore.UpdateRules(rules); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	s.engine.ReloadRules()
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// handleToggleRule toggles a rule's enabled state.
func (s *Server) handleToggleRule(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid rule ID"})
		return
	}

	rules := s.cfgStore.GetRules()
	if id < 0 || id >= len(rules) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "rule not found"})
		return
	}

	rules[id].Enabled = !rules[id].Enabled
	if err := s.cfgStore.UpdateRules(rules); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	s.engine.ReloadRules()
	writeJSON(w, http.StatusOK, rules[id])
}

// handleReorderRules reorders the rules list.
func (s *Server) handleReorderRules(w http.ResponseWriter, r *http.Request) {
	var order []int
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON, expected array of indices"})
		return
	}

	rules := s.cfgStore.GetRules()
	if len(order) != len(rules) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "order length mismatch"})
		return
	}

	reordered := make([]config.RuleConfig, len(rules))
	for i, idx := range order {
		if idx < 0 || idx >= len(rules) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid index in order"})
			return
		}
		reordered[i] = rules[idx]
	}

	if err := s.cfgStore.UpdateRules(reordered); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	s.engine.ReloadRules()
	writeJSON(w, http.StatusOK, reordered)
}

// ---- Helpers ----

// writeJSON sends a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// corsMiddleware adds CORS headers for local development.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://localhost:9090" ||
			origin == "http://127.0.0.1:5173" || origin == "http://127.0.0.1:9090" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
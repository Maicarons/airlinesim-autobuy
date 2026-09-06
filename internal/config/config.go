// Package config manages application configuration loading, saving, and hot-reloading.
package config

import (
	"log/slog"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Store is a thread-safe configuration store that supports hot-reloading.
type Store struct {
	mu       sync.RWMutex
	config   *Config
	filePath string
	watcher  *Watcher
}

// New creates a new config store from a YAML file path.
func New(filePath string) (*Store, error) {
	s := &Store{
		filePath: filePath,
	}

	if err := s.Load(); err != nil {
		return nil, err
	}

	return s, nil
}

// Load reads and parses the YAML configuration file.
func (s *Store) Load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	// Migrate old single-server config to multi-server
	if cfg.Server.BaseURL != "" && len(cfg.Servers) == 0 {
		cfg.Servers = []ServerConfig{cfg.Server}
		cfg.Server = ServerConfig{}
	}
	if len(cfg.Servers) == 0 {
		cfg.Servers = []ServerConfig{
			{Host: "free1", BaseURL: "https://free1.airlinesim.aero"},
		}
	}

	// Migrate old single-auth config to multi-auth
	if cfg.Auth.Username != "" && len(cfg.Auths) == 0 {
		cfg.Auths = []AuthConfig{cfg.Auth}
		cfg.Auth = AuthConfig{}
	}
	if len(cfg.Auths) == 0 {
		cfg.Auths = []AuthConfig{{}}
	}

	s.mu.Lock()
	s.config = &cfg
	s.mu.Unlock()

	slog.Info("configuration loaded", "path", s.filePath, "servers", len(cfg.Servers), "auths", len(cfg.Auths))
	return nil
}

// Save writes the current configuration back to the YAML file.
func (s *Store) Save() error {
	s.mu.RLock()
	cfg := s.config
	s.mu.RUnlock()

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return err
	}

	slog.Info("configuration saved", "path", s.filePath)
	return nil
}

// Get returns a copy of the current configuration.
func (s *Store) Get() *Config {
	s.mu.RLock()
	defer s.mu.RUnlock()

// Return a deep copy to prevent concurrent modification
		cfg := *s.config
		cfg.Servers = make([]ServerConfig, len(s.config.Servers))
		copy(cfg.Servers, s.config.Servers)
		cfg.Rules = make([]RuleConfig, len(s.config.Rules))
		copy(cfg.Rules, s.config.Rules)
		return &cfg
}

// Update replaces the configuration and persists it.
func (s *Store) Update(cfg *Config) error {
	s.mu.Lock()
	s.config = cfg
	s.mu.Unlock()

	return s.Save()
}

// GetRules returns a copy of the current rules list.
func (s *Store) GetRules() []RuleConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rules := make([]RuleConfig, len(s.config.Rules))
	copy(rules, s.config.Rules)
	return rules
}

// UpdateRules replaces the rules list and persists the configuration.
func (s *Store) UpdateRules(rules []RuleConfig) error {
	s.mu.Lock()
	s.config.Rules = rules
	s.mu.Unlock()

	return s.Save()
}

// StartWatcher begins watching the config file for changes (hot-reload).
func (s *Store) StartWatcher(onChange func()) error {
	w, err := NewWatcher(s.filePath, func() {
		if err := s.Load(); err != nil {
			slog.Error("failed to reload config", "error", err)
			return
		}
		slog.Info("configuration hot-reloaded")
		if onChange != nil {
			onChange()
		}
	})
	if err != nil {
		return err
	}
	s.watcher = w
	return nil
}

// StopWatcher stops the file watcher.
func (s *Store) StopWatcher() {
	if s.watcher != nil {
		s.watcher.Stop()
	}
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Servers: []ServerConfig{
			{Host: "free1", BaseURL: "https://free1.airlinesim.aero"},
		},
		Auth: AuthConfig{
			SessionFile: "session.json",
		},
		Monitor: MonitorConfig{
			Interval:       30,
			Jitter:         10,
			RequestTimeout: 30,
			MinBalance:     1000000,
		},
		Notifier: NotifierConfig{
			Console: true,
		},
		WebUI: WebUIConfig{
			Enabled: true,
			Host:    "0.0.0.0",
			Port:    9090,
		},
		Rules: []RuleConfig{
			{
				Name:     "示例规则-经济型窄体机",
				Enabled:  false,
				Priority: 10,
				Match: MatchConfig{
					Types:        []string{"B737-800", "A320-200"},
					PriceRange:   PriceRange{Min: 500000, Max: 5000000},
					MaxAge:       15,
					MaxCycles:    30000,
					ConditionMin: 70,
					OfferTypes:   []string{"auction", "immediate"},
					Financing:    []string{"cash", "credit", "lease"},
				},
				Action: ActionConfig{
					AutoBuy:         true,
					Snatch:          false,
					MaxBidIncrement: 100000,
				},
			},
		},
	}
}
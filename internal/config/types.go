// Package config defines the configuration structure for the autobuy application.
package config

import "time"

// Config represents the top-level application configuration.
type Config struct {
	Servers  []ServerConfig `json:"servers" yaml:"servers"`
	Auths    []AuthConfig   `json:"auths" yaml:"auths"`
	Auth     AuthConfig     `json:"auth,omitempty" yaml:"auth,omitempty"` // Deprecated: use Auths
	Monitor  MonitorConfig  `json:"monitor" yaml:"monitor"`
	Notifier NotifierConfig `json:"notifier" yaml:"notifier"`
	WebUI    WebUIConfig    `json:"webui" yaml:"webui"`
	Rules    []RuleConfig   `json:"rules" yaml:"rules"`
	Server   ServerConfig   `json:"server,omitempty" yaml:"server,omitempty"` // Deprecated
}

// ServerConfig defines an AirlineSim game server connection.
type ServerConfig struct {
	Host    string `json:"host" yaml:"host"`
	BaseURL string `json:"base_url" yaml:"base_url"`
}

// AuthConfig defines authentication credentials and session settings.
type AuthConfig struct {
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	SessionFile string `json:"session_file" yaml:"session_file"`
}

// MonitorConfig defines the market monitoring behavior.
type MonitorConfig struct {
	Interval       int     `json:"interval" yaml:"interval"`
	Jitter         int     `json:"jitter" yaml:"jitter"`
	RequestTimeout int     `json:"request_timeout" yaml:"request_timeout"`
	MinBalance     float64 `json:"min_balance" yaml:"min_balance"`
}

// Duration returns the polling interval as a time.Duration.
func (m MonitorConfig) Duration() time.Duration {
	return time.Duration(m.Interval) * time.Second
}

// NotifierConfig defines notification channels.
type NotifierConfig struct {
	Console        bool   `json:"console" yaml:"console"`
	DiscordWebhook string `json:"discord_webhook,omitempty" yaml:"discord_webhook,omitempty"`
}

// WebUIConfig defines the embedded web management interface.
type WebUIConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Host    string `json:"host" yaml:"host"`
	Port    int    `json:"port" yaml:"port"`
}

// RuleConfig defines a single aircraft purchase rule.
type RuleConfig struct {
	Name     string       `json:"name" yaml:"name"`
	Enabled  bool         `json:"enabled" yaml:"enabled"`
	Priority int          `json:"priority" yaml:"priority"`
	ServerID int          `json:"server_id" yaml:"server_id"` // Index into Servers, -1 = all
	AuthID   int          `json:"auth_id" yaml:"auth_id"`     // Index into Auths, -1 = default
	Match    MatchConfig  `json:"match" yaml:"match"`
	Action   ActionConfig `json:"action" yaml:"action"`
}

// MatchConfig defines the conditions for matching aircraft.
type MatchConfig struct {
	FamilyID     string     `json:"family_id" yaml:"family_id"`
	TypeID       string     `json:"type_id" yaml:"type_id"`
	Types        []string   `json:"types" yaml:"types"`
	PriceRange   PriceRange `json:"price_range" yaml:"price_range"`
	MaxAge       int        `json:"max_age" yaml:"max_age"`
	MaxCycles    int        `json:"max_cycles" yaml:"max_cycles"`
	ConditionMin float64    `json:"condition_min" yaml:"condition_min"`
	OfferTypes   []string   `json:"offer_types" yaml:"offer_types"`
	Financing    []string   `json:"financing" yaml:"financing"`
	SortBy       string     `json:"sort_by" yaml:"sort_by"`
}

// PriceRange defines a minimum and maximum price.
type PriceRange struct {
	Min float64 `json:"min" yaml:"min"`
	Max float64 `json:"max" yaml:"max"`
}

// ActionConfig defines what to do when an aircraft matches.
type ActionConfig struct {
	AutoBuy         bool    `json:"auto_buy" yaml:"auto_buy"`
	Snatch          bool    `json:"snatch" yaml:"snatch"`
	MaxBidIncrement float64 `json:"max_bid_increment" yaml:"max_bid_increment"`
}
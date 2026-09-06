// Package collector handles fetching the aircraft market page data.
package collector

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// MarketPage represents the fetched market page content.
type MarketPage struct {
	Body      []byte
	URL       string
	FetchedAt time.Time
}

// FilterParams defines the market page filter parameters.
type FilterParams struct {
	FamilyID string // Wicket aircraft family ID
	TypeID   string // Wicket aircraft type ID
	SortBy   string // Wicket sort option
}

// Collector fetches aircraft market pages from the game server.
type Collector struct {
	client    *http.Client
	baseURL   string
	marketURL string // Discovered Wicket market page URL
}

// New creates a new market collector.
func New(client *http.Client, baseURL string) *Collector {
	return &Collector{
		client:  client,
		baseURL: baseURL,
	}
}

// DiscoverMarketURL navigates to find the current aircraft market page URL.
func (c *Collector) DiscoverMarketURL() error {
	slog.Info("discovering aircraft market URL")

	// The market page is accessible via /app/aircraft/market with a valid session
	// It redirects (302) to the actual Wicket page URL
	url := fmt.Sprintf("%s/app/aircraft/market", c.baseURL)
	resp, err := c.client.Get(url)
	if err != nil {
		slog.Warn("failed to access market URL", "error", err)
		c.marketURL = url
		return nil
	}
	resp.Body.Close()

	// The final URL after redirect is the Wicket page URL
	finalURL := resp.Request.URL.String()
	if finalURL != "" {
		c.marketURL = finalURL
		slog.Info("discovered market URL", "url", finalURL)
	} else {
		c.marketURL = url
	}

	return nil
}

// Fetch retrieves the aircraft market page with optional filters.
func (c *Collector) Fetch(params *FilterParams) (*MarketPage, error) {
	if c.marketURL == "" {
		if err := c.DiscoverMarketURL(); err != nil {
			return nil, fmt.Errorf("no market URL available: %w", err)
		}
	}

	// Build the filtered URL using Wicket form submission pattern
	fetchURL := c.buildFilteredURL(params)

	slog.Debug("fetching market page", "url", fetchURL)

	resp, err := c.client.Get(fetchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch market page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("market page returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read market page: %w", err)
	}

	return &MarketPage{
		Body:      body,
		URL:       fetchURL,
		FetchedAt: time.Now(),
	}, nil
}

// buildFilteredURL constructs the Wicket market URL with filter parameters.
// Wicket uses AJAX-style URL parameters for form submissions.
// Pattern: ./market?5-2.0-tab-panel-filter~aircraftFamily&tab:panel:filter-aircraftFamily=VALUE
func (c *Collector) buildFilteredURL(params *FilterParams) string {
	if params == nil {
		return c.marketURL
	}

	baseURL := c.marketURL

	// Add family filter
	if params.FamilyID != "" {
		baseURL = addWicketParam(baseURL, "tab:panel:filter-aircraftFamily", params.FamilyID)
	}

	// Add type filter
	if params.TypeID != "" {
		baseURL = addWicketParam(baseURL, "tab:panel:filter-aircraftType", params.TypeID)
	}

	// Add sort
	if params.SortBy != "" {
		baseURL = addWicketParam(baseURL, "tab:panel:sorting", params.SortBy)
	}

	return baseURL
}

// addWicketParam adds a Wicket-style query parameter to the URL.
// Wicket uses a special URL format with the parameter name and value.
func addWicketParam(url, name, value string) string {
	separator := "&"
	return url + separator + name + "=" + value
}

// MarketURL returns the discovered market URL.
func (c *Collector) MarketURL() string {
	return c.marketURL
}
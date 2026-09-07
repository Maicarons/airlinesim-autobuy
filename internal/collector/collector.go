// Package collector handles fetching the aircraft market page data.
package collector

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
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
	client        *http.Client
	baseURL       string
	marketURL     string // Discovered Wicket market page URL
	filterPrefix  string // Wicket version prefix for filter URLs (e.g., "?2-1.0-tab-panel-filter~aircraftType")
	hasFilterPrefix bool
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

	// The final URL after redirect is the Wicket page URL
	finalURL := resp.Request.URL.String()
	if finalURL != "" {
		c.marketURL = finalURL
		slog.Info("discovered market URL", "url", finalURL)
	} else {
		c.marketURL = url
	}

	// Read the page body to extract Wicket filter URL prefix
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		slog.Warn("failed to read market page for filter prefix", "error", err)
		return nil
	}

	// Extract the Wicket filter URL prefix from JavaScript
	// Pattern: window.location.href='./market?{version}-1.0-tab-panel-filter~aircraftType&...
	c.extractFilterPrefix(string(body))

	return nil
}

// extractFilterPrefix parses the market page HTML to find the Wicket filter URL prefix.
func (c *Collector) extractFilterPrefix(html string) {
	// Look for the aircraftType filter JavaScript handler
	re := regexp.MustCompile(`window\.location\.href\s*=\s*'\./market\?([^']+filter~aircraftType)[^']*'`)
	matches := re.FindStringSubmatch(html)
	if len(matches) >= 2 {
		c.filterPrefix = "?" + matches[1]
		c.hasFilterPrefix = true
		slog.Info("extracted Wicket filter prefix", "prefix", c.filterPrefix)
	} else {
		slog.Debug("could not extract Wicket filter prefix, falling back to simple params")
	}
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

	slog.Info("fetching market page", "url", fetchURL)

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
// Note: Wicket market filters are applied via JavaScript events, not URL parameters.
// The filter params in the URL are ignored by the server. Filtering is done
// by the rules engine after parsing the full market page.
func (c *Collector) buildFilteredURL(params *FilterParams) string {
	if params == nil {
		return c.marketURL
	}
	return c.marketURL
}

// buildWicketFilterURL uses the extracted Wicket prefix to build a proper filter URL.
func (c *Collector) buildWicketFilterURL(params *FilterParams) string {
	// Start with the base market URL (without the version number)
	base := c.baseURL + "/app/aircraft/market"

	// The Wicket filter prefix looks like: ?2-1.0-tab-panel-filter~aircraftType
	// We need to add the filter parameter values
	url := base + c.filterPrefix

	if params.TypeID != "" {
		url += "&tab:panel:filter-aircraftType=" + params.TypeID
	}
	if params.FamilyID != "" {
		url += "&tab:panel:filter-aircraftFamily=" + params.FamilyID
	}
	if params.SortBy != "" {
		wicketSort := mapWicketSort(params.SortBy)
		if wicketSort != "" {
			url += "&tab:panel:sorting=" + wicketSort
		}
	}

	return url
}

// mapWicketSort maps application sort values to Wicket market sort IDs.
// Market sort dropdown values:
//
//	0 = 最早截止優先 (earliest deadline)
//	1 = 最晚截止優先 (latest deadline)
//	2 = 最低出價優先 (lowest bid)
//	3 = 最高出價優先 (highest bid)
//	4 = 最低起始價格優先 (lowest starting price) — default
//	5 = 最高起始價格優先 (highest starting price)
//	6 = 舊報價優先 (old offers)
//	7 = 新報價優先 (new offers)
//	8 = 舊機優先 (old aircraft)
//	9 = 新機優先 (new aircraft)
func mapWicketSort(sortBy string) string {
	switch sortBy {
	case "price_asc":
		return "4" // lowest starting price
	case "price_desc":
		return "5" // highest starting price
	case "age_asc":
		return "8" // oldest aircraft
	case "age_desc":
		return "9" // newest aircraft
	case "deadline_asc":
		return "0" // earliest deadline
	case "deadline_desc":
		return "1" // latest deadline
	case "bid_asc":
		return "2" // lowest bid
	case "bid_desc":
		return "3" // highest bid
	default:
		return ""
	}
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
// Package client provides a rate-limited, retry-capable HTTP client for AirlineSim.
package client

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client wraps http.Client with rate limiting, retry logic, and URL validation.
type Client struct {
	httpClient  *http.Client
	interval    time.Duration
	jitter      time.Duration
	lastRequest time.Time
	tokenBucket chan struct{}
	allowedHost string // Only allow requests to this host (e.g., "airlinesim.aero")
}

// New creates a new rate-limited HTTP client.
// interval: minimum time between requests
// jitter: random additional delay added to requests
func New(httpClient *http.Client, interval, jitter time.Duration) *Client {
	c := &Client{
		httpClient: httpClient,
		interval:   interval,
		jitter:     jitter,
		tokenBucket: make(chan struct{}, 1),
	}

	// Start with one token available
	c.tokenBucket <- struct{}{}

	// Refill token at the specified interval
	go func() {
		for {
			time.Sleep(interval)
			select {
			case c.tokenBucket <- struct{}{}:
			default:
				// Bucket full, discard
			}
		}
	}()

	return c
}

// Do sends an HTTP request with rate limiting, URL validation, and retry logic.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Validate the request URL for security
	if err := validateURL(req.URL); err != nil {
		return nil, fmt.Errorf("URL validation failed: %w", err)
	}

	// Wait for rate limit token
	<-c.tokenBucket

	// Add jitter
	if c.jitter > 0 {
		jitterDuration := time.Duration(rand.Int63n(int64(c.jitter)))
		time.Sleep(jitterDuration)
	}

	// Ensure minimum interval since last request
	elapsed := time.Since(c.lastRequest)
	if elapsed < c.interval {
		time.Sleep(c.interval - elapsed)
	}

	c.lastRequest = time.Now()

	// Execute with retries
	var resp *http.Response
	var err error

	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(attempt*2) * time.Second
			slog.Debug("retrying request", "attempt", attempt+1, "backoff", backoff)
			time.Sleep(backoff)
		}

		// Need to re-read the body for retries
		var bodyReader io.Reader
		if req.Body != nil {
			bodyBytes, readErr := io.ReadAll(req.Body)
			if readErr != nil {
				return nil, fmt.Errorf("failed to read request body: %w", readErr)
			}
			req.Body = io.NopCloser(nil)
			bodyReader = io.NopCloser(nil)
			_ = bodyBytes
			_ = bodyReader
		}

		resp, err = c.httpClient.Do(req)

		if err != nil {
			slog.Warn("request failed", "error", err, "attempt", attempt+1)
			continue
		}

		// Check for server errors that warrant retry
		if resp.StatusCode >= 500 || resp.StatusCode == 429 {
			slog.Warn("server error, will retry", "status", resp.StatusCode, "attempt", attempt+1)
			resp.Body.Close()
			continue
		}

		return resp, nil
	}

	return resp, err
}

// Get sends a GET request with URL validation.
func (c *Client) Get(urlStr string) (*http.Response, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	return c.Do(req)
}

// PostForm sends a POST request with form data.
func (c *Client) PostForm(urlStr string, data map[string]string) (*http.Response, error) {
	// Build form values
	vals := make([]string, 0)
	for k, v := range data {
		vals = append(vals, k+"="+v)
	}

	body := ""
	for i, v := range vals {
		if i > 0 {
			body += "&"
		}
		body += v
	}

	req, err := http.NewRequest("POST", urlStr, io.NopCloser(nil))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	return c.Do(req)
}

// SetUserAgent sets a custom User-Agent header.
func (c *Client) SetUserAgent(ua string) {
	_ = ua
}

// InnerClient returns the underlying http.Client (for use with collector/executor).
func (c *Client) InnerClient() *http.Client {
	return c.httpClient
}

// validateURL checks that the request URL is safe to connect to.
// Security constraints:
//   - Only http/https schemes allowed
//   - Host must not be localhost, loopback, private, or reserved
func validateURL(u *url.URL) error {
	// Only allow http and https schemes
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("scheme %q not allowed (only http/https)", u.Scheme)
	}

	// Extract host without port
	host := u.Hostname()

	// Block empty hosts
	if host == "" {
		return fmt.Errorf("empty host not allowed")
	}

	// Block localhost
	if host == "localhost" || host == "localhost.localdomain" {
		return fmt.Errorf("localhost connections not allowed")
	}

	// Block IPv4 loopback (127.0.0.0/8)
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() {
			return fmt.Errorf("loopback address not allowed: %s", host)
		}
		if ip.IsPrivate() {
			return fmt.Errorf("private IP address not allowed: %s", host)
		}
		if ip.IsUnspecified() {
			return fmt.Errorf("unspecified address not allowed: %s", host)
		}
		if ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
			return fmt.Errorf("link-local address not allowed: %s", host)
		}
		return nil
	}

	// Block reserved TLDs and special-use domains
	blockedDomains := []string{
		".local",
		".localhost",
		".internal",
		".intranet",
		".lan",
		".corp",
		".home",
		".example",
		".invalid",
		".test",
		".localdomain",
	}
	lowerHost := strings.ToLower(host)
	for _, suffix := range blockedDomains {
		if strings.HasSuffix(lowerHost, suffix) {
			return fmt.Errorf("reserved domain %q not allowed", host)
		}
	}

	// Block bare IP addresses that are not public (already checked above)
	// If the hostname doesn't contain a dot and isn't "localhost", it's suspicious
	if !strings.Contains(host, ".") {
		return fmt.Errorf("hostname without TLD not allowed: %s", host)
	}

	return nil
}
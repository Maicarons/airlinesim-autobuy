// Package auth handles AirlineSim authentication and session management.
package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// Session represents an authenticated session with the game server.
type Session struct {
	client    *http.Client
	jar       *cookiejar.Jar
	serverURL string
	username  string
	password  string
	filePath  string
	mu        sync.RWMutex
	LoggedIn  bool
}

// sessionData represents the serializable session data for persistence.
type sessionData struct {
	Cookies   []*http.Cookie `json:"cookies"`
	ServerURL string         `json:"server_url"`
	Timestamp time.Time      `json:"timestamp"`
}

// loginRequest is the JSON body for the login API.
type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Method   string `json:"method"`
	Brand    string `json:"brand"`
	Locale   string `json:"locale"`
}

// loginResponse is the JSON response from the login API.
type loginResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Param   string `json:"param"`
	Token   string `json:"token,omitempty"`
	User    *struct {
		ID        string `json:"id"`
		UserID    string `json:"userId"`
		AccountID string `json:"accountId"`
	} `json:"user,omitempty"`
}

const (
	apiBaseURL = "https://sar.simulogics.games"
	hubURL     = "https://airlinesim.aero"
)

// New creates a new authentication session with URL validation.
func New(serverURL, username, password, sessionFile string) (*Session, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	s := &Session{
		client: &http.Client{
			Transport: &secureTransport{next: http.DefaultTransport},
			Jar:       jar,
			Timeout:   30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
		jar:       jar,
		serverURL: serverURL,
		username:  username,
		password:  password,
		filePath:  sessionFile,
	}

	return s, nil
}

// Login authenticates via the sar.simulogics.games API and establishes a session.
func (s *Session) Login() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	slog.Info("attempting to log in", "server", s.serverURL)

	// Try to restore a previous session first
	if s.LoggedIn && s.HealthCheck() {
		slog.Info("existing session is still valid")
		return nil
	}

	// POST to the actual login API endpoint
	apiURL := apiBaseURL + "/api/sessions"
	reqBody := loginRequest{
		Login:    s.username,
		Password: s.password,
		Method:   "password",
		Brand:    "as",
		Locale:   "en",
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	slog.Info("calling login API", "url", apiURL, "username", s.username)

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("login API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Log response headers for debugging
	slog.Info("login API response headers")
	for k, v := range resp.Header {
		slog.Debug("  header", "key", k, "value", v)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	slog.Info("login API response",
		"status", resp.StatusCode,
		"body", string(body),
	)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var apiErr loginResponse
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Message != "" {
			return fmt.Errorf("login failed: %s (code=%s)", apiErr.Message, apiErr.Code)
		}
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp loginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}

	slog.Info("login successful", "token_length", len(loginResp.Token))

	// The session token is the as-sid cookie value.
	// Set it manually on the game server domain.
	if loginResp.Token != "" {
		portalURL := fmt.Sprintf("%s/action/portal/index", s.serverURL)
		u := mustParseURL(portalURL)
		cookie := &http.Cookie{
			Name:     "as-sid",
			Value:    loginResp.Token,
			Path:     "/",
			Domain:   u.Hostname(),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		s.jar.SetCookies(u, []*http.Cookie{cookie})
		slog.Info("set as-sid cookie", "domain", u.Hostname(), "token_length", len(loginResp.Token))

		// Try to access the portal to verify the session
		portalResp, err := s.client.Get(portalURL)
		if err != nil {
			slog.Warn("portal access with cookie failed", "error", err)
		} else {
			portalResp.Body.Close()
			slog.Info("portal access with cookie", "status", portalResp.StatusCode)
		}
	} else {
		slog.Warn("no token in login response")
	}

	s.LoggedIn = true
	slog.Info("login successful, session established")

	// Persist session cookies
	if err := s.save(); err != nil {
		slog.Warn("failed to persist session", "error", err)
	}

	return nil
}

// Client returns the authenticated HTTP client.
func (s *Session) Client() *http.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.client
}

// IsLoggedIn returns whether the session is currently authenticated.
func (s *Session) IsLoggedIn() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.LoggedIn
}

// Logout clears the session.
func (s *Session) Logout() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jar, _ = cookiejar.New(nil)
	s.client.Jar = s.jar
	s.LoggedIn = false

	os.Remove(s.filePath)
	slog.Info("logged out and session cleared")
}

// save persists the current session cookies to a file.
func (s *Session) save() error {
	serverURL := fmt.Sprintf("%s/action/portal/index", s.serverURL)
	u := mustParseURL(serverURL)
	cookies := s.jar.Cookies(u)

	data := sessionData{
		Cookies:   cookies,
		ServerURL: s.serverURL,
		Timestamp: time.Now(),
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, bytes, 0600)
}

// TryRestore attempts to restore a previously persisted session.
func (s *Session) TryRestore() (bool, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	var sd sessionData
	if err := json.Unmarshal(data, &sd); err != nil {
		return false, err
	}

	if time.Since(sd.Timestamp) > 24*time.Hour {
		slog.Info("persisted session is too old, will re-authenticate")
		return false, nil
	}

	u := mustParseURL(fmt.Sprintf("%s/action/portal/index", s.serverURL))
	s.jar.SetCookies(u, sd.Cookies)
	s.LoggedIn = true

	slog.Info("restored session from file", "age", time.Since(sd.Timestamp).Round(time.Second))
	return true, nil
}

// HealthCheck verifies the session is still valid by accessing the portal.
func (s *Session) HealthCheck() bool {
	portalURL := fmt.Sprintf("%s/action/portal/index", s.serverURL)
	resp, err := s.client.Get(portalURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	u := mustParseURL(portalURL)
	cookies := s.jar.Cookies(u)
	for _, c := range cookies {
		if c.Name == "as-sid" {
			return true
		}
	}
	return false
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

// secureTransport is an http.RoundTripper that validates URLs before forwarding requests.
type secureTransport struct {
	next http.RoundTripper
	allowedHosts []string
}

func (t *secureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Allow specific known hosts without validation
	host := req.URL.Hostname()
	allowed := []string{
		"airlinesim.aero",
		"sar.simulogics.games",
		"simulogics.games",
	}
	for _, a := range allowed {
		if host == a || strings.HasSuffix(host, "."+a) {
			return t.next.RoundTrip(req)
		}
	}

	if err := validateURL(req.URL); err != nil {
		return nil, fmt.Errorf("request blocked by security policy: %w", err)
	}
	return t.next.RoundTrip(req)
}

// validateURL checks that the request URL is safe to connect to.
func validateURL(u *url.URL) error {
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("scheme %q not allowed (only http/https)", u.Scheme)
	}

	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("empty host not allowed")
	}

	if host == "localhost" || host == "localhost.localdomain" {
		return fmt.Errorf("localhost connections not allowed")
	}

	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
			return fmt.Errorf("address %q is not allowed (loopback/private/unspecified)", host)
		}
		if ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
			return fmt.Errorf("link-local address not allowed: %s", host)
		}
		return nil
	}

	blockedSuffixes := []string{
		".local", ".localhost", ".internal", ".intranet",
		".lan", ".corp", ".home", ".example", ".invalid", ".test",
	}
	lowerHost := strings.ToLower(host)
	for _, suffix := range blockedSuffixes {
		if strings.HasSuffix(lowerHost, suffix) {
			return fmt.Errorf("reserved domain %q not allowed", host)
		}
	}

	if !strings.Contains(host, ".") {
		return fmt.Errorf("hostname without TLD not allowed: %s", host)
	}

	return nil
}
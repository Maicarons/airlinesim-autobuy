# Architecture

This document describes the architecture, component design, data flow, and concurrency model of AirlineSim Autobuy.

## Overview

AirlineSim Autobuy follows a **pipeline architecture** where data flows through a series of processing stages, each handled by a dedicated package. The application is written in Go and uses a single-binary deployment model with an embedded Vue 3 frontend.

```
┌─────────────────────────────────────────────────────────────┐
│                     Application Entrypoint                  │
│                       cmd/autobuy/main.go                   │
└──────────────────────────┬──────────────────────────────────┘
                           │
          ┌────────────────┼────────────────┐
          ▼                ▼                ▼
┌─────────────────┐ ┌────────────┐ ┌────────────────┐
│   Config Store   │ │   Engine   │ │   Web UI       │
│  internal/config │ │  internal/ │ │  internal/webui│
│                  │ │  engine    │ │                │
└────────┬────────┘ │            │ │  ┌──────────┐  │
         │          │ ┌────────┐ │ │  │  Vue SPA │  │
         │          │ │ Server │ │ │  │ (embed)  │  │
         │          │ │ Runner │ │ │  └──────────┘  │
         │          │ └───┬────┘ │ └────────────────┘
         │          │     │      │
         │          │     ▼      │
         │          │  ┌─────────────────────────────┐
         │          │  │      Pipeline per Server     │
         │          │  │                              │
         │          │  │  ┌──────────┐  ┌─────────┐  │
         │          │  │  │ Auth     │  │ Client  │  │
         │          │  │  │ (session)│  │(rate-ltd)│  │
         │          │  │  └────┬─────┘  └────┬────┘  │
         │          │  │       │              │       │
         │          │  │  ┌────▼──────────────▼────┐  │
         │          │  │  │      Collector          │  │
         │          │  │  │   (market fetcher)      │  │
         │          │  │  └───────────┬────────────┘  │
         │          │  │              │               │
         │          │  │  ┌───────────▼────────────┐  │
         │          │  │  │      Parser            │  │
         │          │  │  │   (HTML → Aircraft)    │  │
         │          │  │  └───────────┬────────────┘  │
         │          │  │              │               │
         │          │  │  ┌───────────▼────────────┐  │
         │          │  │  │    Rules Engine         │  │
         │          │  │  │   (match & score)      │  │
         │          │  │  └───────────┬────────────┘  │
         │          │  │              │               │
         │          │  │  ┌───────────▼────────────┐  │
         │          │  │  │     Executor            │  │
         │          │  │  │   (purchase/bid)       │  │
         │          │  │  └───────────┬────────────┘  │
         │          │  │              │               │
         │          │  │  ┌───────────▼────────────┐  │
         │          │  │  │     Notifier            │  │
         │          │  │  │   (console/discord)    │  │
         │          │  │  └────────────────────────┘  │
         │          │  └─────────────────────────────┘
         ▼          ▼
┌──────────────────────────────────────────────────────┐
│                    Configuration File                 │
│                     configs/config.yaml               │
└──────────────────────────────────────────────────────┘
```

## Component Overview

### Entry Point (`cmd/autobuy/main.go`)

The main function orchestrates startup:

1. Initializes structured logging
2. Ensures a configuration file exists (creates default if missing)
3. Loads the configuration into a thread-safe store
4. Creates the notification system
5. Initializes the engine with the config store and notifier
6. Optionally starts the web UI server
7. Waits for OS signals (SIGINT, SIGTERM) for graceful shutdown

### Config Store (`internal/config`)

The config store manages the application configuration with thread-safe access:

- **`Store`** -- Thread-safe wrapper around the configuration struct with `sync.RWMutex`
- **`Load()`** -- Reads and parses the YAML configuration file
- **`Save()`** -- Writes the current configuration to disk
- **`Get()`** -- Returns a deep copy of the configuration
- **`Update()`** -- Replaces the configuration and persists it
- **`StartWatcher()`** -- Begins polling the config file for changes (hot-reload)
- **`DefaultConfig()`** -- Returns a configuration with sensible defaults

The config watcher polls the file every 5 seconds for modification time changes. When a change is detected, the configuration is reloaded and the rules engine is updated.

### Auth (`internal/auth`)

Handles AirlineSim authentication and session management:

- **`Session`** -- Manages authenticated HTTP sessions with cookie persistence
- **`Login()`** -- Authenticates via the `sar.simulogics.games` API using password-based login
- **`HealthCheck()`** -- Verifies the session is still valid by accessing the game portal
- **`TryRestore()`** -- Attempts to restore a previously persisted session from disk
- **`Logout()`** -- Clears the session and removes the persisted session file
- **`Client()`** -- Returns the authenticated HTTP client with session cookies

The session is persisted to a JSON file for reuse across restarts. Sessions expire after 24 hours or when the server invalidates the token.

### Client (`internal/client`)

Provides a rate-limited, retry-capable HTTP client with security validation:

- **Token bucket rate limiter** -- Ensures minimum interval between requests
- **Jitter injection** -- Adds random delay to avoid pattern detection
- **Automatic retry** -- Retries up to 3 times on server errors (5xx, 429)
- **URL validation** -- Blocks requests to localhost, private IPs, and reserved domains
- **User-Agent masking** -- Uses a browser-like User-Agent header

### Collector (`internal/collector`)

Fetches aircraft market pages from the game server:

- **`DiscoverMarketURL()`** -- Navigates to find the current Wicket market page URL (follows redirects)
- **`Fetch()`** -- Retrieves the market page with optional filter parameters
- **`buildFilteredURL()`** -- Constructs Wicket-style URLs with family, type, and sort filters

### Parser (`internal/parser`)

Parses Wicket market page HTML into structured aircraft data:

- Uses `goquery` (jQuery-like selector engine) for HTML parsing
- Extracts aircraft type, price, age, condition, cycles, location, and offer type
- Handles multiple number formats (US and European decimal separators)
- Generates unique IDs for each aircraft offer
- Detects auction vs. immediate purchase offers

### Rules Engine (`internal/rules`)

Evaluates aircraft offers against configured purchase rules:

- **`Evaluate()`** -- Checks an aircraft against all enabled rules, returns sorted matches
- **`matchRule()`** -- Applies match conditions: type, price, age, cycles, condition, offer type
- **`calculateScore()`** -- Computes a desirability score based on price, condition, age, offer type, and rule priority
- **`BestMatch()`** -- Returns the single best match for an aircraft

### Executor (`internal/executor`)

Handles the actual purchase or bid execution:

- **`Execute()`** -- Routes to `buyImmediately()` or `placeBid()` based on offer type
- **`buyImmediately()`** -- Navigates to the aircraft detail page and submits the purchase form
- **`placeBid()`** -- Submits a bid on an auction listing
- **`VerifyPrice()`** -- Checks the current price before making a purchase

### Notifier (`internal/notifier`)

Sends notifications through configured channels:

- **Console notifications** -- Structured slog output with emoji indicators
- **Discord webhook** -- Sends rich embeds to a Discord channel (webhook URL configured)
- **Event types** -- Aircraft found, purchase made, purchase failed, bid placed, error, info

### Engine (`internal/engine`)

Orchestrates the entire monitoring pipeline:

- **`New()`** -- Creates server runners for each configured server
- **`Start()`** -- Begins monitoring loops for all servers in separate goroutines
- **`Stop()`** -- Gracefully stops all monitoring loops
- **`Status()`** -- Returns aggregated status across all servers
- **`ReloadRules()`** -- Updates the rules engine with the latest configuration

### Web UI (`internal/webui`)

Serves the embedded Vue 3 frontend and REST API:

- Uses `chi` router for HTTP routing
- Serves embedded SPA files with fallback routing for Vue Router
- CORS middleware for local development
- REST API endpoints for engine control, configuration, rules management, and aircraft data

## Data Flow

### Monitoring Cycle

For each server, the monitoring cycle follows this flow:

```
1. Session Health Check
   │
   ├── Healthy → Continue
   └── Expired → Re-authenticate → Fail → Log error, skip cycle
                                        │
                                        ▼
                                   Return early
   
2. Build Filter Parameters
   │
   └── Read first enabled rule for this server
       └── Extract family_id, type_id, sort_by
   
3. Fetch Market Page
   │
   ├── Success → Continue
   └── Failure → Log error, skip cycle
   
4. Parse HTML
   │
   ├── Success → Continue
   └── Failure → Log error, skip cycle
   
5. Process Each Aircraft
   │
   ├── Already seen? → Skip (deduplication)
   │
   ├── Check rules → No match → Skip
   │
   ├── Match found → Notify "aircraft found"
   │
   └── Auto-buy enabled?
       ├── No → Done
       └── Yes → Execute purchase
                  ├── Success → Notify "purchase made"
                  └── Failure → Notify "purchase failed"
   
6. Wait for next interval (interval + jitter)
```

### Data Structures

#### AircraftOffer

```go
type AircraftOffer struct {
    ID        string    // Unique identifier (type + timestamp)
    Type      string    // Aircraft type name
    Family    string    // Aircraft family name
    Age       int       // Age in years
    Cycles    int       // Flight cycles
    Condition float64   // Condition percentage (0-100)
    Owner     string    // Current owner
    Price     float64   // Base price (AS$)
    LeaseRate float64   // Weekly lease rate (AS$)
    LeaseDep  float64   // Lease deposit (AS$)
    DownPmt   float64   // Down payment (AS$)
    Install   float64   // Weekly installment (AS$)
    HasBid    bool      // Whether there are existing bids
    OfferType string    // "auction" or "immediate"
    Financing []string  // Available payment methods
    Location  string    // Airport code
    BidCount  int       // Number of existing bids
    URL       string    // Aircraft detail page URL
    SeenAt    time.Time // When the offer was first seen
}
```

#### MatchResult

```go
type MatchResult struct {
    Aircraft  *AircraftOffer
    Rule      *RuleConfig
    Score     float64
    ShouldBuy bool
}
```

## Concurrency Model

### Goroutine Per Server

The engine starts one goroutine per configured server, each running an independent monitoring loop:

```go
for _, sr := range e.servers {
    go e.runServerLoop(ctx, sr)
}
```

Each goroutine:
1. Has its own HTTP client with session cookies
2. Runs its own ticker for the polling interval
3. Independently fetches, parses, and processes market data
4. Shares the same rules engine and config store (read-only access)

### Thread-Safe Shared State

| Component | Protection | Access Pattern |
|---|---|---|
| Config Store | `sync.RWMutex` | Read-heavy, infrequent writes |
| Engine Status | `sync.RWMutex` | Updated by server loops, read by web UI |
| Seen Aircraft | `sync.Mutex` | Updated by server loops, read for deduplication |

### Graceful Shutdown

The engine uses a context-based cancellation pattern:

1. `main()` creates a cancellable context
2. Each server goroutine receives the context
3. On shutdown, the context is cancelled, causing all goroutines to exit
4. The engine waits for goroutines to clean up

## Security Architecture

### URL Validation

Both the client and auth packages implement URL validation to prevent SSRF attacks:

- Only `http://` and `https://` schemes allowed
- Localhost, loopback, and private IP addresses blocked
- Reserved TLDs (`.local`, `.internal`, `.test`, etc.) blocked
- Hostnames without TLDs blocked
- Only connections to `airlinesim.aero`, `sar.simulogics.games`, and `simulogics.games` domains are allowed without validation

### Session Security

- Session cookies are stored in memory during runtime
- Persisted session files use `0600` permissions (owner read/write only)
- Session tokens are never exposed in logs (length is logged, not the value)

## Directory Structure

```
airlinesim-autobuy/
├── cmd/autobuy/           # Application entry point
├── internal/
│   ├── auth/              # Authentication and session management
│   ├── client/            # Rate-limited HTTP client
│   ├── collector/         # Market page fetcher
│   ├── config/            # Configuration management (file, store, watcher)
│   ├── engine/            # Pipeline orchestration
│   ├── executor/          # Purchase execution
│   ├── marketdata/        # Aircraft reference data
│   ├── notifier/          # Notification system
│   ├── parser/            # HTML market page parser
│   ├── rules/             # Rules engine
│   └── webui/             # Web UI server (Vue SPA + REST API)
│       └── frontend/      # Vue 3 + Vite frontend source
├── configs/               # YAML configuration files
├── docs/                  # VitePress documentation
├── scripts/               # Utility scripts
├── Makefile               # Build targets
├── go.mod                 # Go module definition
└── go.sum                 # Go module checksums
```

## See Also

- [API Reference](/reference/api) for REST API documentation
- [Configuration Reference](/reference/config) for all configuration fields
- [Contributing Guide](/reference/contributing) for development setup
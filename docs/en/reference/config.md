# Configuration Reference

Complete reference for all configuration fields, types, defaults, and descriptions.

## Top-Level Structure

```yaml
servers:     # Game server connections (optional)
auths:       # Authentication credentials (optional)
auth:        # Single auth (deprecated)
monitor:     # Monitoring behavior (optional)
notifier:    # Notification channels (optional)
webui:       # Web management interface (optional)
rules:       # Purchase rules (optional)
server:      # Single server (deprecated)
```

## Servers

List of game servers to monitor.

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
```

### ServerConfig

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `host` | string | Yes | -- | Short label for the server (used in logs and UI) |
| `base_url` | string | Yes | -- | Full URL of the game server (e.g., `https://free1.airlinesim.aero`) |

**Default:** If no servers are configured, a single server is created:
```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
```

## Auths

List of authentication credentials.

```yaml
auths:
  - username: "player@example.com"
    password: "s3cret"
    session_file: "session.json"
```

### AuthConfig

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `username` | string | Yes | -- | AirlineSim account email or username |
| `password` | string | Yes | -- | Account password |
| `session_file` | string | No | `"session.json"` | File path to persist session cookies for reuse across restarts |

**Note:** The `auth` field (singular) is deprecated. Use `auths` (plural) instead.

## Monitor

Market scanning behavior.

```yaml
monitor:
  interval: 30
  jitter: 10
  request_timeout: 30
  min_balance: 1000000
```

### MonitorConfig

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `interval` | int | No | `30` | Polling interval in seconds between market scans. Minimum recommended: 15 seconds |
| `jitter` | int | No | `10` | Random delay in seconds added to each interval. Actual delay = `interval + random(0, jitter)` |
| `request_timeout` | int | No | `30` | HTTP request timeout in seconds for market page fetches |
| `min_balance` | float | No | `1000000` | Minimum account balance (AS$) to retain after a purchase. Acts as a safety warning threshold |

## Notifier

Notification channel configuration.

```yaml
notifier:
  console: true
  discord_webhook: "https://discord.com/api/webhooks/..."
```

### NotifierConfig

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `console` | bool | No | `true` | Enable structured console logging with emoji indicators |
| `discord_webhook` | string | No | `""` | Discord webhook URL for sending notifications to a Discord channel. Leave empty to disable |

## WebUI

Embedded web management interface configuration.

```yaml
webui:
  enabled: true
  host: 0.0.0.0
  port: 9090
```

### WebUIConfig

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `enabled` | bool | No | `true` | Enable or disable the web UI HTTP server |
| `host` | string | No | `"0.0.0.0"` | Bind address for the HTTP server. Use `"127.0.0.1"` for local-only access |
| `port` | int | No | `9090` | HTTP server port number |

## Rules

Purchase rules define which aircraft to monitor and purchase.

```yaml
rules:
  - name: "Example Rule"
    enabled: true
    priority: 10
    server_id: 0
    auth_id: 0
    match:
      family_id: "1200300"
      type_id: "16"
      types:
        - "Airbus A320-200 heavy"
      price_range:
        min: 500000
        max: 5000000
      max_age: 15
      max_cycles: 30000
      condition_min: 70.0
      offer_types:
        - auction
        - immediate
      financing:
        - cash
        - credit
        - lease
      sort_by: "price_asc"
    action:
      auto_buy: true
      snatch: false
      max_bid_increment: 100000
```

### RuleConfig

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `name` | string | Yes | -- | Human-readable name for the rule. Displayed in logs and UI |
| `enabled` | bool | No | `true` | Enable or disable the rule. Disabled rules are skipped during evaluation |
| `priority` | int | No | `10` | Priority score added to matching aircraft. Higher priority = higher desirability score |
| `server_id` | int | No | `-1` | Index of the target server in the `servers` list. `-1` means all servers |
| `auth_id` | int | No | `-1` | Index of the authentication account in the `auths` list. `-1` means the first auth |
| `match` | MatchConfig | Yes | -- | Aircraft match conditions |
| `action` | ActionConfig | Yes | -- | Purchase action settings |

### MatchConfig

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `family_id` | string | No | `""` | Wicket aircraft family ID for server-side filtering. Empty string means no filter |
| `type_id` | string | No | `""` | Wicket aircraft type ID for server-side filtering. Empty string means no filter |
| `types` | array of strings | No | `[]` | List of aircraft type names to match. Case-insensitive. Empty array means no filter |
| `price_range` | PriceRange | No | `{min: 0, max: 0}` | Price range in AS$. `0` means unbounded |
| `max_age` | int | No | `0` | Maximum aircraft age in years. `0` means no limit |
| `max_cycles` | int | No | `0` | Maximum flight cycles. `0` means no limit |
| `condition_min` | float | No | `0` | Minimum condition percentage (0-100). `0` means no limit |
| `offer_types` | array of strings | No | `[]` | Allowed offer types. Valid values: `"auction"`, `"immediate"`. Empty array means all types |
| `financing` | array of strings | No | `[]` | Required payment methods. Valid values: `"cash"`, `"credit"`, `"lease"`. Empty array means all methods |
| `sort_by` | string | No | `""` | Server-side sort parameter. Common values: `"price_asc"`, `"price_desc"`, `"age_asc"`, `"age_desc"` |

### PriceRange

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `min` | float | No | `0` | Minimum price in AS$. `0` or negative means no minimum |
| `max` | float | No | `0` | Maximum price in AS$. `0` or negative means no maximum |

### ActionConfig

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `auto_buy` | bool | No | `false` | Automatically purchase matching aircraft. When `false`, matching is logged but no purchase is attempted |
| `snatch` | bool | No | `false` | Enable snatch mode for immediate purchase with aggressive pricing |
| `max_bid_increment` | float | No | `0` | Maximum amount (AS$) to increment the current bid for auction items. `0` means no increment |

## Market Data

The following aircraft families and types are available for use in rules.

### Aircraft Families

| ID | Name |
|---|---|
| `7200100` | 1900 Airliner |
| `3000100` | 208 Caravan |
| `3000200` | 408 SkyCourier |
| `6200100` | 410 |
| `2200410` | 737 MAX |
| `2200350` | 737-300/400/500 |
| `2200400` | 737-600/700/800/900 |
| `2200510` | 747-8 |
| `2200600` | 767-200/300/400 |
| `2200650` | 777-200/300 |
| `2200700` | 787 |
| `1200250` | A220 |
| `1200299` | A318 / A319 |
| `1200310` | A319 / A320 / A321 NEO |
| `1200300` | A320 / A321 |
| `1200400` | A330 |
| `1200410` | A330 NEO |
| `1200600` | A350 |
| `1200800` | A380 |
| `7100100` | AN-28 |
| `1400600` | AN140 |
| `1400650` | AN148 |
| `3200100` | ARJ21 |
| `1600100` | ATR 42 |
| `1600200` | ATR 72 |
| `3200200` | C919 |
| `2400200` | CRJ Series |
| `3800400` | Dash 8 |
| `3800200` | DHC-6 |
| `7300100` | DO228NG |
| `4200300` | EMB 170/175/190/195 |
| `4200400` | EMB 175/190/195 E2 |
| `4200200` | ERJ 135/140/145 |
| `2600100` | Islander / Trislander |
| `7005100` | PC-12 |
| `7005200` | PC-24 |
| `7900100` | Superjet |
| `8000300` | TU-204/214 |
| `1500200` | Xian Y-7 |

### Aircraft Types (Partial)

The full list of aircraft types is available through the `/api/aircraft-data` endpoint and includes all variants of the families above. Key types include:

| ID | Name |
|---|---|
| `16` | Airbus A320-200 heavy |
| `17` | Airbus A320-200 heavy (enhanced) |
| `24` | Airbus A320-200 medium |
| `28` | Airbus A320neo heavy |
| `34` | Airbus A321-200 heavy |
| `42` | Airbus A321LR |
| `43` | Airbus A321neo heavy |
| `49` | Airbus A330-200 |
| `53` | Airbus A330-300 |
| `62` | Airbus A350-1000 |
| `63` | Airbus A350-900 |
| `65` | Airbus A380-800 |
| `83` | Boeing 737-700 BGW |
| `93` | Boeing 737-8 |
| `97` | Boeing 737-800 HGW |
| `105` | Boeing 737-9 |
| `113` | Boeing 747-8F |
| `114` | Boeing 747-8I |
| `115` | Boeing 767-300ER |
| `119` | Boeing 777-200ER |
| `122` | Boeing 777-300ER |
| `123` | Boeing 787-10 |
| `124` | Boeing 787-8 |
| `125` | Boeing 787-9 |

## Default Configuration

The default configuration is created when no configuration file exists:

```yaml
servers:
  - host: "free1"
    base_url: "https://free1.airlinesim.aero"
auth:
  username: ""
  password: ""
  session_file: "session.json"
monitor:
  interval: 30
  jitter: 10
  request_timeout: 30
  min_balance: 1000000
notifier:
  console: true
webui:
  enabled: true
  host: "0.0.0.0"
  port: 9090
rules:
  - name: "示例规则-经济型窄体机"
    enabled: false
    priority: 10
    match:
      types:
        - "B737-800"
        - "A320-200"
      price_range:
        min: 500000
        max: 5000000
      max_age: 15
      max_cycles: 30000
      condition_min: 70
      offer_types:
        - auction
        - immediate
      financing:
        - cash
        - credit
        - lease
    action:
      auto_buy: true
      snatch: false
      max_bid_increment: 100000
```

## See Also

- [Configuration Guide](/guide/configuration) for configuration walkthrough
- [Purchase Rules](/guide/rules) for detailed rule configuration
- [API Reference](/reference/api) for programmatic configuration
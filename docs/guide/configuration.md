# Configuration

AirlineSim Autobuy uses a YAML configuration file to define servers, authentication, monitoring behavior, notifications, web UI settings, and purchase rules. The default location is `configs/config.yaml`.

## Configuration File Structure

The configuration file has the following top-level sections:

```yaml
servers:      # Game server connections (multi-server)
auths:        # Authentication credentials (multi-account)
auth:         # Single auth (deprecated, use auths)
monitor:      # Monitoring behavior
notifier:     # Notification channels
webui:        # Web management interface
rules:        # Purchase rules
server:       # Single server (deprecated, use servers)
```

## Servers

Define one or more AirlineSim game servers to monitor.

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
  - host: free2
    base_url: https://free2.airlinesim.aero
```

| Field | Type | Required | Description |
|---|---|---|---|
| `host` | string | Yes | Short label for the server (used in logs and UI) |
| `base_url` | string | Yes | Full URL of the game server |

### Default

If no servers are configured, the application creates a default entry:

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
```

### Legacy Single Server

You can use the deprecated `server` field instead of `servers`:

```yaml
server:
  host: free1
  base_url: https://free1.airlinesim.aero
```

If both `server` and `servers` are present, the `servers` list takes precedence.

## Auths (Authentication)

Configure credentials for one or more AirlineSim game accounts.

```yaml
auths:
  - username: "player1@example.com"
    password: "s3cret123"
    session_file: "session_player1.json"
  - username: "player2@example.com"
    password: "s3cret456"
    session_file: "session_player2.json"
```

| Field | Type | Required | Description |
|---|---|---|---|
| `username` | string | Yes | AirlineSim account email or username |
| `password` | string | Yes | Account password |
| `session_file` | string | No | Path to persist session cookies (default: `session.json`) |

### How Auths Map to Servers

The authentication account used for a server is determined by the rules that target that server. If no rule specifies an `auth_id`, the first auth entry is used.

### Legacy Single Auth

You can use the deprecated `auth` field instead of `auths`:

```yaml
auth:
  username: "your@email.com"
  password: "your_password"
  session_file: "session.json"
```

If both `auth` and `auths` are present, the `auths` list takes precedence.

### Session Persistence

After a successful login, session cookies are saved to the file specified by `session_file`. On subsequent starts, the application attempts to restore the session from this file, avoiding repeated login requests. Sessions expire after 24 hours or when the server invalidates them.

## Monitor

Configure the market scanning behavior.

```yaml
monitor:
  interval: 30
  jitter: 10
  request_timeout: 30
  min_balance: 1000000
```

| Field | Type | Default | Description |
|---|---|---|---|
| `interval` | int | `30` | Polling interval in seconds between market scans |
| `jitter` | int | `10` | Random delay in seconds added to each interval to avoid pattern detection |
| `request_timeout` | int | `30` | HTTP request timeout in seconds |
| `min_balance` | float | `1000000` | Minimum account balance (AS$) to retain after a purchase. Acts as a safety net |

### Interval and Jitter

The actual time between scans is `interval + random(0, jitter)`. For example, with `interval: 30` and `jitter: 10`, the actual delay is between 30 and 40 seconds.

### Min Balance

The `min_balance` field is a safety threshold. If the purchase price exceeds this value, the application logs a warning. This is a soft check -- the actual balance is not queried from the game server, so use this as a conservative ceiling.

## Notifier

Configure notification channels for events.

```yaml
notifier:
  console: true
  discord_webhook: "https://discord.com/api/webhooks/..."
  # dingtalk_webhook: "https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN"
  # dingtalk_secret: "YOUR_SIGNING_SECRET"
```

| Field | Type | Default | Description |
|---|---|---|---|
| `console` | bool | `true` | Enable structured console logging with emoji indicators |
| `discord_webhook` | string | empty | Discord webhook URL for sending notifications to a Discord channel |
| `dingtalk_webhook` | string | empty | DingTalk custom robot webhook URL for sending notifications to a DingTalk group |
| `dingtalk_secret` | string | empty | HMAC-SHA256 signing secret for the DingTalk robot (optional, required if the robot has signature verification enabled) |

### Console Notifications

When `console: true`, events are logged with emoji prefixes:

| Event | Emoji | Example |
|---|---|---|
| Aircraft Found | 🔍 | `🔍 Aircraft found matching rule 'My Rule': A320-200 heavy` |
| Purchase Made | 🛒 | `🛒 PURCHASED A320-200 heavy for AS$ 2500000 (rule: My Rule)` |
| Purchase Failed | ❌ | `❌ Failed to purchase A320-200 heavy: bid failed` |
| Error | ⚠️ | `⚠️ Error: re-authentication failed` |
| Bid Placed | 💰 | `💰 Bid of AS$ 100000 placed successfully` |
| Info | ℹ️ | `ℹ️ Engine started` |

### Discord Notifications

To enable Discord notifications, create a webhook in your Discord server settings and add the URL to the config:

```yaml
notifier:
  console: true
  discord_webhook: "https://discord.com/api/webhooks/123456789/abcdef"
```

Discord messages include the same event information as console notifications, formatted as rich embeds.

### DingTalk Notifications

To enable DingTalk notifications, create a custom robot in a DingTalk group and add the webhook URL to the config:

```yaml
notifier:
  console: true
  dingtalk_webhook: "https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN"
  # dingtalk_secret: "YOUR_SIGNING_SECRET"  # Required if signature verification is enabled
```

If the robot has "Signature Verification" enabled, set the `dingtalk_secret` field with the signing secret. Messages are sent as plain text with the same emoji prefixes as console notifications.

## WebUI

Configure the embedded web management interface.

```yaml
webui:
  enabled: true
  host: 0.0.0.0
  port: 9090
```

| Field | Type | Default | Description |
|---|---|---|---|
| `enabled` | bool | `true` | Enable or disable the web UI server |
| `host` | string | `"0.0.0.0"` | Bind address. Use `"127.0.0.1"` for local-only access |
| `port` | int | `9090` | HTTP server port |

### Security Notes

- For local-only access, bind to `127.0.0.1` instead of `0.0.0.0`
- The web UI does not have built-in authentication. Consider placing a reverse proxy (e.g., Nginx, Caddy) in front of it for production use
- The web UI is served over plain HTTP. For internet-facing deployments, use a TLS-terminating reverse proxy

## Full Example Configuration

Here is a complete configuration file with multiple servers, multiple accounts, and sample rules:

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
  - host: free2
    base_url: https://free2.airlinesim.aero

auths:
  - username: "player1@example.com"
    password: "s3cret123"
    session_file: "session_p1.json"
  - username: "player2@example.com"
    password: "s3cret456"
    session_file: "session_p2.json"

monitor:
  interval: 30
  jitter: 10
  request_timeout: 30
  min_balance: 500000

notifier:
  console: true
  discord_webhook: "https://discord.com/api/webhooks/abc123/def456"
  # dingtalk_webhook: "https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN"
  # dingtalk_secret: "YOUR_SIGNING_SECRET"

webui:
  enabled: true
  host: 127.0.0.1
  port: 9090

rules:
  - name: "Narrow-body targets"
    enabled: true
    priority: 10
    server_id: 0
    auth_id: 0
    match:
      types:
        - "Airbus A320-200 heavy"
        - "Boeing 737-800 HGW"
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
      sort_by: price_asc
    action:
      auto_buy: true
      snatch: false
      max_bid_increment: 100000

  - name: "Wide-body lease monitoring"
    enabled: true
    priority: 5
    server_id: 1
    auth_id: 1
    match:
      types:
        - "Boeing 787-9"
      price_range:
        max: 15000000
      max_age: 10
    action:
      auto_buy: false
      max_bid_increment: 200000
```

## Hot Reloading

The application watches the configuration file for changes using a polling watcher (checks every 5 seconds). When the file is modified:

1. The configuration is re-read from disk
2. The rules engine is updated with any new rules
3. Changes to server settings, auth, or monitor parameters require a restart

> **Note:** Hot reloading updates rules in real time. Changes to `servers`, `auths`, `monitor`, or `webui` sections require restarting the application to take effect.

## Default Configuration

If no configuration file exists at startup, the application creates one with sensible defaults:

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
auth:
  username: ""
  password: ""
  session_file: session.json
monitor:
  interval: 30
  jitter: 10
  request_timeout: 30
  min_balance: 1000000
notifier:
  console: true
  # discord_webhook: ""
  # dingtalk_webhook: ""
  # dingtalk_secret: ""
webui:
  enabled: true
  host: 0.0.0.0
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

## Environment Variables

The application does not currently support environment variable overrides. All configuration must be provided through the YAML file.

## See Also

- [Full Configuration Reference](/reference/config) for all fields, types, and defaults
- [Purchase Rules](/guide/rules) for detailed rule configuration
- [Deployment Guide](/guide/deployment) for production setup
# Web UI

AirlineSim Autobuy includes a full-featured web management interface built with Vue 3 and Vite. The frontend is embedded in the Go binary, so no separate web server is required.

## Accessing the Web UI

By default, the web UI is available at `http://localhost:9090`.

To change the bind address or port, modify the `webui` section in your configuration:

```yaml
webui:
  enabled: true
  host: 0.0.0.0      # Bind to all interfaces
  port: 9090          # HTTP port
```

### Local-Only Access

For security, bind to localhost when you do not need remote access:

```yaml
webui:
  host: 127.0.0.1
```

## Dashboard

The dashboard is the main page displayed at `http://localhost:9090`. It provides a real-time overview of the monitoring engine.

### Engine Status

The top of the dashboard shows the engine status:

- **Running** -- Engine is actively scanning markets
- **Stopped** -- Engine is idle
- **Start/Stop Button** -- Toggle the engine on and off

### Server Cards

Each configured server is displayed as a card showing:

| Metric | Description |
|---|---|
| Host | Server label (e.g., `free1`) |
| Scan Count | Total number of market scans performed |
| Found Count | Number of aircraft that matched rules |
| Bought Count | Number of successful purchases |
| Failed Count | Number of failed purchase attempts |
| Last Scan | Timestamp of the most recent scan |
| Last Error | Most recent error message, if any |

### Statistics Summary

Aggregate statistics across all servers:

- Total scans
- Total found aircraft
- Total purchases
- Total failures

### Log Feed

A real-time log feed displays recent events, including:

- Aircraft found matching rules
- Purchase confirmations
- Purchase failures
- Errors and warnings
- Informational messages

Each log entry includes a timestamp, event type, and message details.

## Rules Page

Navigate to the Rules page to manage purchase rules through the web interface.

### Rules List

The rules page displays all configured rules in a table with columns:

| Column | Description |
|---|---|
| # | Rule index (position in the list) |
| Name | Rule name |
| Enabled | Toggle switch to enable/disable the rule |
| Priority | Rule priority value |
| Match | Summary of match conditions (types, price range, age, etc.) |
| Action | Summary of action settings (auto-buy, snatch, max bid) |
| Server | Which server the rule targets |
| Actions | Edit, toggle, delete buttons |

### Creating a Rule

Click the **Add Rule** button to open the rule editor. Fill in the following fields:

**Basic Settings:**
- **Name** -- A descriptive name for the rule
- **Enabled** -- Enable the rule immediately
- **Priority** -- Priority score (higher = more important)
- **Server ID** -- Target server index (or `-1` for all)
- **Auth ID** -- Target auth index (or `-1` for default)

**Match Conditions:**
- **Family** -- Select an aircraft family from the dropdown
- **Type** -- Select specific aircraft types from the dropdown
- **Price Min/Max** -- Price range in AS$
- **Max Age** -- Maximum aircraft age in years
- **Max Cycles** -- Maximum flight cycles
- **Min Condition** -- Minimum condition percentage
- **Offer Types** -- Checkboxes for auction and immediate
- **Financing** -- Checkboxes for cash, credit, lease
- **Sort By** -- Server-side sort order

**Action Settings:**
- **Auto Buy** -- Enable automatic purchase
- **Snatch** -- Enable snatch mode
- **Max Bid Increment** -- Maximum bid increment for auctions

### Editing a Rule

Click the **Edit** button on any rule to modify its settings. The rule editor is pre-filled with the current values.

### Enabling/Disabling a Rule

Click the **Toggle** button or the enabled switch to quickly enable or disable a rule without opening the editor.

### Deleting a Rule

Click the **Delete** button to remove a rule. The deletion is immediate and the configuration file is updated.

### Reordering Rules

Use the drag handle or the **Reorder** feature to change rule priority order. Rules are evaluated in list order.

## Settings Page

The Settings page allows you to view and modify global configuration settings.

### Server Settings

- Add, edit, or remove game server entries
- Each server requires a host label and base URL

### Authentication Settings

- Add, edit, or remove authentication accounts
- Each account requires username, password, and optional session file path

### Monitor Settings

| Setting | Description |
|---|---|
| Interval | Polling interval in seconds |
| Jitter | Random delay in seconds |
| Request Timeout | HTTP request timeout in seconds |
| Min Balance | Minimum account balance to retain (AS$) |

### Notifier Settings

| Setting | Description |
|---|---|
| Console | Enable/disable console logging |
| Discord Webhook | Discord webhook URL for notifications |

### Web UI Settings

| Setting | Description |
|---|---|
| Host | Bind address for the web UI |
| Port | HTTP port for the web UI |

### Saving Settings

Changes to settings are saved to the configuration file immediately. Some changes (like server and auth modifications) require a restart of the application to take effect.

## Logs Page

The Logs page provides a dedicated view of the application log stream.

### Log Stream

Real-time streaming log entries from the running application. Each entry includes:

- **Timestamp** -- When the event occurred
- **Level** -- INFO, WARN, ERROR
- **Message** -- The log message
- **Details** -- Additional structured data

### Log Levels

| Level | Color | Description |
|---|---|---|
| INFO | Blue | Normal operation events |
| WARN | Yellow | Warnings that do not require immediate action |
| ERROR | Red | Errors that may require attention |

### Filtering

You can filter logs by:

- **Level** -- Show only specific log levels
- **Search** -- Text search across log messages
- **Time range** -- Show logs within a specific time window

## API Endpoints

The web UI interacts with the backend through a REST API. These endpoints are also available for external integrations.

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/status` | Get engine status and statistics |
| `POST` | `/api/control/start` | Start the monitoring engine |
| `POST` | `/api/control/stop` | Stop the monitoring engine |
| `POST` | `/api/reload` | Reload rules without restart |
| `GET` | `/api/config` | Get full configuration |
| `PUT` | `/api/config` | Update configuration |
| `GET` | `/api/rules` | Get all rules |
| `POST` | `/api/rules` | Create a new rule |
| `GET` | `/api/rules/{id}` | Get a specific rule |
| `PUT` | `/api/rules/{id}` | Update a specific rule |
| `DELETE` | `/api/rules/{id}` | Delete a rule |
| `PATCH` | `/api/rules/{id}/toggle` | Toggle a rule's enabled state |
| `PUT` | `/api/rules/reorder` | Reorder rules |
| `GET` | `/api/aircraft-data` | Get aircraft families and types |

For detailed API documentation, see the [API Reference](/reference/api).

## See Also

- [API Reference](/reference/api) for complete endpoint documentation
- [Configuration Guide](/guide/configuration) for web UI configuration options
- [Deployment Guide](/guide/deployment) for reverse proxy setup
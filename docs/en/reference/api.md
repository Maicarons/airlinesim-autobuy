# API Reference

The web UI exposes a REST API for engine control, configuration management, rules management, and data retrieval. All endpoints return JSON responses.

## Base URL

All API endpoints are prefixed with `/api`:

```
http://localhost:9090/api
```

## Authentication

The API does not have built-in authentication. For production deployments, place a reverse proxy with authentication in front of the web UI.

## Content Type

All request bodies must be `application/json`. All responses are `application/json`.

## Common Response Codes

| Code | Description |
|---|---|
| `200` | Success |
| `201` | Created (resource created successfully) |
| `400` | Bad request (invalid input) |
| `404` | Not found (resource does not exist) |
| `409` | Conflict (e.g., engine already running) |
| `500` | Internal server error |

## Endpoints

### Status

#### GET /api/status

Returns the current engine status and per-server statistics.

**Response:**

```json
{
  "running": true,
  "start_time": "2024-01-15T10:30:00Z",
  "servers": [
    {
      "host": "free1",
      "scan_count": 42,
      "found_count": 5,
      "bought_count": 2,
      "failed_count": 0,
      "last_scan": "2024-01-15T11:00:00Z",
      "last_error": ""
    }
  ],
  "scan_count": 42,
  "found_count": 5,
  "bought_count": 2,
  "failed_count": 0,
  "last_error": ""
}
```

**Fields:**

| Field | Type | Description |
|---|---|---|
| `running` | bool | Whether the engine is currently running |
| `start_time` | string (ISO 8601) | When the engine was last started |
| `servers` | array | Per-server statistics |
| `servers[].host` | string | Server label |
| `servers[].scan_count` | int | Number of market scans performed |
| `servers[].found_count` | int | Number of aircraft that matched rules |
| `servers[].bought_count` | int | Number of successful purchases |
| `servers[].failed_count` | int | Number of failed purchase attempts |
| `servers[].last_scan` | string (ISO 8601) | Timestamp of the most recent scan |
| `servers[].last_error` | string | Most recent error message |
| `scan_count` | int | Aggregate scan count across all servers |
| `found_count` | int | Aggregate found count across all servers |
| `bought_count` | int | Aggregate purchase count across all servers |
| `failed_count` | int | Aggregate failure count across all servers |
| `last_error` | string | Most recent error across all servers |

---

### Control

#### POST /api/control/start

Starts the monitoring engine.

**Response (success):**

```json
{
  "status": "started"
}
```

**Response (already running):**

```json
{
  "error": "engine is already running"
}
```

Status code: `409 Conflict`

#### POST /api/control/stop

Stops the monitoring engine.

**Response:**

```json
{
  "status": "stopped"
}
```

---

### Reload

#### POST /api/reload

Reloads rules from the configuration file without restarting the engine. Useful after editing the config file manually.

**Response:**

```json
{
  "status": "reloaded"
}
```

---

### Configuration

#### GET /api/config

Returns the full current configuration.

**Response:**

```json
{
  "servers": [
    {
      "host": "free1",
      "base_url": "https://free1.airlinesim.aero"
    }
  ],
  "auths": [
    {
      "username": "user@example.com",
      "password": "****",
      "session_file": "session.json"
    }
  ],
  "monitor": {
    "interval": 30,
    "jitter": 10,
    "request_timeout": 30,
    "min_balance": 1000000
  },
  "notifier": {
    "console": true,
    "discord_webhook": ""
  },
  "webui": {
    "enabled": true,
    "host": "0.0.0.0",
    "port": 9090
  },
  "rules": [
    {
      "name": "Example Rule",
      "enabled": true,
      "priority": 10,
      "server_id": 0,
      "auth_id": 0,
      "match": {
        "family_id": "",
        "type_id": "",
        "types": ["Airbus A320-200 heavy"],
        "price_range": { "min": 0, "max": 5000000 },
        "max_age": 15,
        "max_cycles": 30000,
        "condition_min": 70,
        "offer_types": ["auction", "immediate"],
        "financing": ["cash", "credit", "lease"],
        "sort_by": "price_asc"
      },
      "action": {
        "auto_buy": false,
        "snatch": false,
        "max_bid_increment": 100000
      }
    }
  ]
}
```

#### PUT /api/config

Updates the global configuration. Partial updates are supported -- only the fields that are present and non-zero will be updated.

**Request Body:**

```json
{
  "monitor": {
    "interval": 60,
    "jitter": 15
  },
  "notifier": {
    "console": true,
    "discord_webhook": "https://discord.com/api/webhooks/..."
  }
}
```

**Update Rules:**

- **Servers:** Replaced entirely if `servers` array is non-empty
- **Auths:** Replaced entirely if `auths` array is non-empty
- **Monitor:** Individual fields updated if non-zero
- **Notifier:** `console` is always updated; `discord_webhook` updated if non-empty
- **WebUI:** `host` updated if non-empty; `port` updated if non-zero
- **Rules:** Not updated through this endpoint (use `/api/rules` endpoints)

**Response:** Returns the full updated configuration object.

---

### Rules

#### GET /api/rules

Returns all configured rules.

**Response:**

```json
[
  {
    "name": "Example Rule",
    "enabled": true,
    "priority": 10,
    "server_id": 0,
    "auth_id": 0,
    "match": {
      "family_id": "",
      "type_id": "",
      "types": ["Airbus A320-200 heavy"],
      "price_range": { "min": 0, "max": 5000000 },
      "max_age": 15,
      "max_cycles": 30000,
      "condition_min": 70,
      "offer_types": ["auction", "immediate"],
      "financing": ["cash", "credit", "lease"],
      "sort_by": "price_asc"
    },
    "action": {
      "auto_buy": false,
      "snatch": false,
      "max_bid_increment": 100000
    }
  }
]
```

#### POST /api/rules

Creates a new rule. The rule is appended to the end of the rules list.

**Request Body:**

```json
{
  "name": "New Rule",
  "enabled": true,
  "priority": 10,
  "server_id": -1,
  "auth_id": -1,
  "match": {
    "types": ["Boeing 787-9"],
    "price_range": { "max": 15000000 },
    "max_age": 10
  },
  "action": {
    "auto_buy": false,
    "max_bid_increment": 200000
  }
}
```

**Response:** Returns the created rule. Status code: `201 Created`

#### GET /api/rules/{id}

Returns a single rule by its index (0-based).

**Response:**

```json
{
  "name": "Example Rule",
  "enabled": true,
  ...
}
```

**Error Response:**

```json
{
  "error": "rule not found"
}
```

Status code: `404 Not Found`

#### PUT /api/rules/{id}

Updates a rule by its index. The entire rule object is replaced.

**Request Body:**

```json
{
  "name": "Updated Rule",
  "enabled": true,
  "priority": 20,
  ...
}
```

**Response:** Returns the updated rule.

#### DELETE /api/rules/{id}

Deletes a rule by its index.

**Response:**

```json
{
  "status": "deleted"
}
```

#### PATCH /api/rules/{id}/toggle

Toggles the `enabled` state of a rule.

**Response:** Returns the rule with the toggled `enabled` value.

#### PUT /api/rules/reorder

Reorders the rules list. The request body is an array of indices representing the new order.

**Request Body:**

```json
[2, 0, 1]
```

This moves the rule at index 2 to position 0, index 0 to position 1, and index 1 to position 2.

**Response:** Returns the reordered rules array.

---

### Aircraft Data

#### GET /api/aircraft-data

Returns the available aircraft families and types for use in UI dropdowns.

**Response:**

```json
{
  "families": [
    {
      "id": "1200300",
      "name": "A320 / A321"
    },
    {
      "id": "2200400",
      "name": "737-600/700/800/900"
    }
  ],
  "types": [
    {
      "id": "16",
      "name": "Airbus A320-200 heavy"
    },
    {
      "id": "97",
      "name": "Boeing 737-800 HGW"
    }
  ]
}
```

## SPA Fallback

All non-API routes serve the Vue 3 single-page application. The SPA has a fallback handler that returns `index.html` for any unmatched route, enabling Vue Router to handle client-side routing.

## CORS

The API includes CORS middleware that allows requests from the following origins:

- `http://localhost:5173` (Vite dev server)
- `http://localhost:9090`
- `http://127.0.0.1:5173`
- `http://127.0.0.1:9090`

Allowed methods: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS`

## Error Responses

All error responses follow this format:

```json
{
  "error": "description of the error"
}
```

## See Also

- [Web UI Guide](/guide/webui) for using the web interface
- [Architecture](/reference/architecture) for understanding the system design
- [Configuration Reference](/reference/config) for configuration field details
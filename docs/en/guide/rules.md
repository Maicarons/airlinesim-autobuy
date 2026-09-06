# Purchase Rules

Purchase rules define which aircraft the autobuy system should monitor, evaluate, and potentially purchase. Rules are evaluated in order of priority, and matching aircraft are scored for desirability.

## Rule Structure

Each rule is a YAML object with the following structure:

```yaml
rules:
  - name: "My Rule"
    enabled: true
    priority: 10
    server_id: 0
    auth_id: 0
    match:
      # ... match conditions ...
    action:
      # ... action settings ...
```

### Top-Level Fields

| Field | Type | Default | Description |
|---|---|---|---|
| `name` | string | Required | Human-readable name for the rule, displayed in logs and UI |
| `enabled` | bool | `true` | Enable or disable the rule. Disabled rules are skipped during evaluation |
| `priority` | int | `10` | Priority score added to matched aircraft. Higher priority rules produce higher scores |
| `server_id` | int | `-1` | Index of the server in the `servers` list. `-1` means all servers |
| `auth_id` | int | `-1` | Index of the auth entry in the `auths` list. `-1` means the first auth |
| `match` | object | Required | Match conditions for filtering aircraft |
| `action` | object | Required | Action settings when an aircraft matches |

## Match Conditions

The `match` section defines the criteria for filtering aircraft listings.

### Family

Filter by aircraft family ID. The family ID is a numeric identifier used by the AirlineSim Wicket API.

```yaml
match:
  family_id: "1200300"  # A320 / A321 family
```

| Field | Type | Description |
|---|---|---|
| `family_id` | string | Wicket aircraft family ID. Use the Web UI's aircraft data dropdown to find IDs |

Available family IDs include:

| Family ID | Name |
|---|---|
| `1200300` | A320 / A321 |
| `1200310` | A319 / A320 / A321 NEO |
| `2200400` | 737-600/700/800/900 |
| `2200410` | 737 MAX |
| `1200400` | A330 |
| `1200600` | A350 |
| `2200700` | 787 |
| `4200300` | EMB 170/175/190/195 |

See the [market data reference](/reference/config#market-data) for the complete list.

### Type

Filter by specific aircraft type name(s).

```yaml
match:
  type_id: "16"                    # Single type by Wicket ID
  types:                           # Multiple types by name
    - "Airbus A320-200 heavy"
    - "Boeing 737-800 HGW"
```

| Field | Type | Description |
|---|---|---|
| `type_id` | string | Wicket aircraft type ID. Used for server-side filtering |
| `types` | array of strings | List of aircraft type names to match. Case-insensitive matching |

You can use `types` to match multiple aircraft variants. The check is case-insensitive, so `"a320-200 heavy"` matches `"Airbus A320-200 heavy"`.

### Price

Filter by price range.

```yaml
match:
  price_range:
    min: 500000
    max: 5000000
```

| Field | Type | Description |
|---|---|---|
| `price_range.min` | float | Minimum price in AS$. `0` or negative means no minimum |
| `price_range.max` | float | Maximum price in AS$. `0` or negative means no maximum |

Set either field to `0` to leave it unbounded.

### Age

Filter by aircraft age in years.

```yaml
match:
  max_age: 15
```

| Field | Type | Description |
|---|---|---|
| `max_age` | int | Maximum age in years. `0` means no limit |

### Cycles

Filter by flight cycles (takeoff/landing cycles).

```yaml
match:
  max_cycles: 30000
```

| Field | Type | Description |
|---|---|---|
| `max_cycles` | int | Maximum flight cycles. `0` means no limit |

### Condition

Filter by minimum aircraft condition percentage.

```yaml
match:
  condition_min: 70
```

| Field | Type | Description |
|---|---|---|
| `condition_min` | float | Minimum condition percentage (0-100). `0` means no limit |

### Offer Types

Filter by the type of market offer.

```yaml
match:
  offer_types:
    - auction
    - immediate
```

| Field | Type | Description |
|---|---|---|
| `offer_types` | array of strings | Allowed offer types. Empty array means all types |

Valid values:

| Value | Description |
|---|---|
| `auction` | Aircraft listed for auction (bidding required) |
| `immediate` | Aircraft available for immediate purchase (buy now) |

### Financing

Filter by available financing/payment methods.

```yaml
match:
  financing:
    - cash
    - credit
    - lease
```

| Field | Type | Description |
|---|---|---|
| `financing` | array of strings | Required payment methods. Empty array means all methods |

Valid values:

| Value | Description |
|---|---|
| `cash` | Full cash purchase |
| `credit` | Purchase with credit (financing) |
| `lease` | Operating lease |

### Sort By

Set the server-side sorting for market page results.

```yaml
match:
  sort_by: price_asc
```

| Field | Type | Description |
|---|---|---|
| `sort_by` | string | Wicket sort parameter. Passed to the server as-is |

Common values: `price_asc`, `price_desc`, `age_asc`, `age_desc`.

## Action Settings

The `action` section defines what happens when an aircraft matches the rule.

```yaml
action:
  auto_buy: true
  snatch: false
  max_bid_increment: 100000
```

### Auto Buy

```yaml
action:
  auto_buy: true
```

| Field | Type | Default | Description |
|---|---|---|---|
| `auto_buy` | bool | `false` | Automatically purchase matching aircraft. When `false`, matching is logged but no purchase is attempted |

Set to `false` for monitoring-only rules. This lets you verify your match conditions are correct before enabling automated purchases.

### Snatch

```yaml
action:
  snatch: false
```

| Field | Type | Default | Description |
|---|---|---|---|
| `snatch` | bool | `false` | Enable snatch mode for immediate purchase with aggressive pricing |

Snatch mode is designed for highly desirable listings. When enabled, the system attempts to purchase immediately at the listed price without waiting for auction cycles.

### Max Bid Increment

```yaml
action:
  max_bid_increment: 100000
```

| Field | Type | Default | Description |
|---|---|---|---|
| `max_bid_increment` | float | `0` | Maximum amount (AS$) to increment the current bid for auction items |

For auction listings, this controls how much above the current price the system is willing to bid.

## Server and Auth Binding

### Server ID

The `server_id` field binds a rule to a specific server by its index in the `servers` list:

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
  - host: free2
    base_url: https://free2.airlinesim.aero

rules:
  - name: "Free1 only"
    server_id: 0    # Targets free1

  - name: "Free2 only"
    server_id: 1    # Targets free2

  - name: "All servers"
    server_id: -1   # Targets all servers
```

### Auth ID

The `auth_id` field binds a rule to a specific authentication account:

```yaml
auths:
  - username: "account1@example.com"
    password: "..."
  - username: "account2@example.com"
    password: "..."

rules:
  - name: "Account 1 rules"
    auth_id: 0

  - name: "Account 2 rules"
    auth_id: 1
```

## Rule Evaluation and Scoring

When the engine scans the market, it evaluates each discovered aircraft against all enabled rules. The scoring system determines which rule to use.

### Scoring Formula

```
Score = (Price Score) + (Condition Score) + (Age Score) + (Immediate Bonus) + (Priority)
```

| Component | Weight | Calculation |
|---|---|---|
| Price Score | 50 points max | `(1 - price / max_price) * 50`. Lower price = higher score |
| Condition Score | 20 points max | `(condition - min_condition) / (100 - min_condition) * 20`. Better condition = higher score |
| Age Score | 20 points max | `(1 - age / max_age) * 20`. Younger aircraft = higher score |
| Immediate Bonus | 10 points | `+10` if the offer is immediate purchase (not auction) |
| Priority | Variable | Rule's `priority` value added directly |

### Evaluation Order

1. All enabled rules are evaluated against each aircraft
2. Rules that do not match are discarded
3. Matching results are sorted by score descending
4. The highest-scoring match determines whether to purchase
5. Purchase is triggered if `auto_buy` is `true` on the best-matching rule

## Example Rules

### Monitoring Only (No Auto-Purchase)

```yaml
rules:
  - name: "Monitor A320-200"
    enabled: true
    priority: 10
    match:
      types:
        - "Airbus A320-200 heavy"
      price_range:
        max: 5000000
      max_age: 20
      max_cycles: 50000
      condition_min: 50
      offer_types:
        - auction
        - immediate
    action:
      auto_buy: false
      max_bid_increment: 100000
```

### Aggressive Auction Sniping

```yaml
rules:
  - name: "Snatch good deals"
    enabled: true
    priority: 20
    match:
      types:
        - "Airbus A320neo heavy"
        - "Boeing 737 MAX 8"
      price_range:
        max: 8000000
      max_age: 5
      condition_min: 85
    action:
      auto_buy: true
      snatch: true
      max_bid_increment: 250000
```

### Lease-Only Monitoring

```yaml
rules:
  - name: "Lease opportunities"
    enabled: true
    priority: 5
    match:
      family_id: "1200300"
      price_range:
        max: 10000000
      financing:
        - lease
    action:
      auto_buy: false
```

### Multi-Server, Multi-Account

```yaml
rules:
  - name: "Free1 narrow-body"
    enabled: true
    server_id: 0
    auth_id: 0
    match:
      types:
        - "Airbus A320-200 heavy"
      price_range:
        max: 3000000
    action:
      auto_buy: true

  - name: "Free2 wide-body"
    enabled: true
    server_id: 1
    auth_id: 1
    match:
      types:
        - "Boeing 787-9"
      price_range:
        max: 15000000
    action:
      auto_buy: false
```

## Best Practices

1. **Start with monitoring-only rules** -- Set `auto_buy: false` until you verify the matching behavior
2. **Use specific types** -- Prefer the `types` list over generic `family_id` matching to avoid unexpected purchases
3. **Set conservative price limits** -- Always define a `price_range.max` to prevent overspending
4. **Use priority for preference** -- Higher priority rules take precedence when multiple rules match the same aircraft
5. **Test with different scenarios** -- Create rules for different market conditions and monitor which ones trigger

## See Also

- [Configuration Guide](/guide/configuration) for full config file structure
- [Web UI Rules Page](/guide/webui#rules-page) for managing rules through the interface
- [Configuration Reference](/reference/config) for complete field documentation
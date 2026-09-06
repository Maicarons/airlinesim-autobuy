# 配置说明

AirlineSim Autobuy 使用 YAML 格式的配置文件，默认路径为 `configs/config.yaml`。程序启动时自动加载，并支持热重载——修改配置文件后无需重启即可生效。

## YAML 结构概览

```yaml
servers:      # 游戏服务器列表（多服务器）
auths:        # 认证账户列表（多账户）
auth:         # （已弃用）单账户配置
monitor:      # 监控设置
notifier:     # 通知设置
webui:        # Web 管理界面设置
rules:        # 购买规则列表
```

## 服务器配置（多服务器）

`servers` 字段是一个数组，每个元素定义一个 AirlineSim 游戏服务器。你可以同时监控多个服务器。

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
  - host: free2
    base_url: https://free2.airlinesim.aero
  - host: free3
    base_url: https://free3.airlinesim.aero
```

### 字段说明

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `host` | string | 是 | 服务器标识名，用于显示和日志 |
| `base_url` | string | 是 | 服务器基础 URL |

### 默认服务器

如果未配置 `servers`，程序会自动使用以下默认值：

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
```

## 认证配置（多账户）

`auths` 字段是一个数组，每个元素定义一个游戏账户的认证信息。每个规则可以指定使用哪个账户来执行购买。

```yaml
auths:
  - username: "player1@example.com"
    password: "password1"
    session_file: session1.json
  - username: "player2@example.com"
    password: "password2"
    session_file: session2.json
```

### 字段说明

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `username` | string | 是 | AirlineSim 账户邮箱或用户名 |
| `password` | string | 是 | 账户密码 |
| `session_file` | string | 否 | 会话持久化文件路径，默认 `session.json` |

### 会话管理

认证系统通过 `sar.simulogics.games` API 进行登录，获取 `as-sid` 令牌并持久化到本地文件。后续启动时会尝试恢复之前的会话，如果会话过期则自动重新登录。

### 单账户兼容

`auth` 字段（已弃用）仍然支持单账户配置，当不设置 `auths` 时可以使用：

```yaml
auth:
  username: "your_username"
  password: "your_password"
  session_file: session.json
```

## 监控设置

`monitor` 字段控制市场数据采集的行为。

```yaml
monitor:
  interval: 30
  jitter: 10
  request_timeout: 30
  min_balance: 1000000
```

### 字段说明

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `interval` | int | `30` | 轮询间隔（秒），每次市场扫描之间的等待时间 |
| `jitter` | int | `10` | 随机延迟（秒），在每次请求前添加随机延迟，避免被检测为自动化工具 |
| `request_timeout` | int | `30` | HTTP 请求超时时间（秒） |
| `min_balance` | float | `1000000` | 最低余额保护（AS$），当购买价格接近此值时发出警告 |

### 关于轮询间隔

- 建议 `interval` 设置为 30-60 秒，避免对游戏服务器造成过大压力
- `jitter` 应设置为 `interval` 的 1/3 到 1/2，使请求时间分布更自然
- 实际请求间隔 = `interval` + 随机 `[0, jitter]` 秒

## 通知设置

`notifier` 字段配置事件通知的投递方式。

```yaml
notifier:
  console: true
  discord_webhook: "https://discord.com/api/webhooks/..."
```

### 字段说明

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `console` | bool | `true` | 是否在控制台输出通知 |
| `discord_webhook` | string | 空 | Discord Webhook URL，设置后会将通知发送到 Discord 频道 |

### 通知事件类型

| 事件类型 | 说明 |
|----------|------|
| 发现飞机 | 检测到符合规则的飞机上架 |
| 购买成功 | 成功完成购买 |
| 购买失败 | 购买操作失败 |
| 出价成功 | 成功提交竞拍出价 |
| 错误 | 系统运行错误 |
| 信息 | 一般信息通知 |

## WebUI 设置

`webui` 字段配置内嵌的 Web 管理界面。

```yaml
webui:
  enabled: true
  host: 0.0.0.0
  port: 9090
```

### 字段说明

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `enabled` | bool | `true` | 是否启用 Web 界面 |
| `host` | string | `0.0.0.0` | 监听地址，`0.0.0.0` 表示所有网络接口 |
| `port` | int | `9090` | 监听端口 |

### 安全建议

- 如果只需要本地访问，将 `host` 设置为 `127.0.0.1`
- 在生产环境中，建议在 WebUI 前添加反向代理（如 Nginx）并配置认证

## 规则配置

`rules` 是一个数组，每个元素定义一条购买规则。详细说明请参阅[购买规则](./rules)。

```yaml
rules:
  - name: "示例规则"
    enabled: true
    priority: 10
    server_id: 0
    auth_id: 0
    match:
      family_id: "1200300"
      type_id: ""
      types:
        - "Airbus A320-200 heavy"
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
      sort_by: price_asc
    action:
      auto_buy: false
      snatch: false
      max_bid_increment: 100000
```

## 配置热重载

程序内置了文件监控器（`config.Watcher`），每 5 秒检查配置文件是否被修改。检测到变化后，会自动重新加载配置并更新规则引擎。

这意味着你可以在不重启程序的情况下，修改规则或调整监控参数——修改保存后，程序会立即应用新配置。

## 完整配置示例

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
  - host: free2
    base_url: https://free2.airlinesim.aero

auths:
  - username: "main@example.com"
    password: "main_password"
    session_file: session_main.json
  - username: "alt@example.com"
    password: "alt_password"
    session_file: session_alt.json

monitor:
  interval: 45
  jitter: 15
  request_timeout: 30
  min_balance: 2000000

notifier:
  console: true
  discord_webhook: "https://discord.com/api/webhooks/..."

webui:
  enabled: true
  host: 127.0.0.1
  port: 9090

rules:
  - name: "经济型窄体机"
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
        max: 3000000
      max_age: 12
      max_cycles: 25000
      condition_min: 75
      offer_types:
        - auction
        - immediate
      financing:
        - cash
        - credit
    action:
      auto_buy: true
      snatch: false
      max_bid_increment: 50000

  - name: "宽体机监控"
    enabled: true
    priority: 5
    server_id: -1
    auth_id: 0
    match:
      family_id: "1200400"
      price_range:
        min: 5000000
        max: 20000000
      max_age: 20
      condition_min: 60
      offer_types:
        - auction
    action:
      auto_buy: false
      max_bid_increment: 200000
```
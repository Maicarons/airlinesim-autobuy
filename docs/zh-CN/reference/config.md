# 配置参考

本文档列出了 AirlineSim Autobuy 所有配置字段的完整参考信息。

## 顶级字段

```yaml
servers:    []ServerConfig    # 游戏服务器列表
auths:      []AuthConfig      # 认证账户列表
auth:       AuthConfig        # [已弃用] 单账户配置
monitor:    MonitorConfig     # 监控设置
notifier:   NotifierConfig    # 通知设置
webui:      WebUIConfig       # Web 管理界面设置
rules:      []RuleConfig      # 购买规则列表
server:     ServerConfig      # [已弃用] 单服务器配置
```

## ServerConfig

定义 AirlineSim 游戏服务器连接。

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `host` | string | 是 | — | 服务器标识名，用于显示和日志 |
| `base_url` | string | 是 | — | 服务器基础 URL，例如 `https://free1.airlinesim.aero` |

**默认值：** 如果未配置任何服务器，自动使用：

```yaml
- host: free1
  base_url: https://free1.airlinesim.aero
```

## AuthConfig

定义认证凭据和会话设置。

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `username` | string | 是 | — | AirlineSim 账户邮箱或用户名 |
| `password` | string | 是 | — | 账户密码 |
| `session_file` | string | 否 | `session.json` | 会话持久化文件路径 |

**注意：** 如果同时配置了 `auth`（已弃用）和 `auths`，优先使用 `auths`。

## MonitorConfig

定义市场监控行为。

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `interval` | int | 否 | `30` | 轮询间隔（秒）。每次市场扫描之间的等待时间。建议 30-60 秒 |
| `jitter` | int | 否 | `10` | 随机延迟（秒）。在每次请求前添加随机延迟，避免请求模式被检测。建议设置为 interval 的 1/3 到 1/2 |
| `request_timeout` | int | 否 | `30` | HTTP 请求超时时间（秒） |
| `min_balance` | float | 否 | `1000000` | 最低余额保护（AS$）。当购买价格接近此值时会发出警告 |

**轮询间隔计算：** 实际请求间隔 = `interval` + 随机值 `[0, jitter]` 秒

## NotifierConfig

定义通知渠道。

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `console` | bool | 否 | `true` | 是否在控制台输出事件通知 |
| `discord_webhook` | string | 否 | 空 | Discord Webhook URL。设置后会将通知发送到指定的 Discord 频道 |

**Discord Webhook 格式：** `https://discord.com/api/webhooks/{webhook.id}/{webhook.token}`

## WebUIConfig

定义内嵌 Web 管理界面。

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `enabled` | bool | 否 | `true` | 是否启用 Web 管理界面 |
| `host` | string | 否 | `0.0.0.0` | 监听地址。`0.0.0.0` 表示所有网络接口，`127.0.0.1` 仅本地访问 |
| `port` | int | 否 | `9090` | 监听端口 |

## RuleConfig

定义一条飞机购买规则。

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `name` | string | 是 | — | 规则名称，用于日志和通知标识 |
| `enabled` | bool | 否 | `true` | 是否启用该规则。禁用后规则不参与匹配 |
| `priority` | int | 否 | `10` | 优先级（1-100）。数值越大，匹配时评分越高 |
| `server_id` | int | 否 | `-1` | 目标服务器索引（从 0 开始）。`-1` 表示匹配所有服务器 |
| `auth_id` | int | 否 | `-1` | 执行账户索引（从 0 开始）。`-1` 使用默认账户 |
| `match` | MatchConfig | 是 | — | 匹配条件 |
| `action` | ActionConfig | 是 | — | 操作设置 |

## MatchConfig

定义飞机匹配条件。

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `family_id` | string | 否 | 空 | Wicket 飞机系列 ID，例如 `"1200300"` 表示 A320/A321 系列 |
| `type_id` | string | 否 | 空 | Wicket 飞机具体型号 ID |
| `types` | string[] | 否 | `[]` | 飞机型号名称列表，支持大小写不敏感匹配 |
| `price_range` | PriceRange | 否 | `{min: 0, max: 0}` | 价格范围（AS$）。`0` 表示不限制 |
| `max_age` | int | 否 | `0` | 最大机龄（年）。`0` 表示不限制 |
| `max_cycles` | int | 否 | `0` | 最大飞行循环数。`0` 表示不限制 |
| `condition_min` | float | 否 | `0` | 最低机身状态百分比（0-100）。`0` 表示不限制 |
| `offer_types` | string[] | 否 | `[]` | 交易类型列表。可选值：`auction`（拍卖）、`immediate`（立即购买）。空数组表示不限制 |
| `financing` | string[] | 否 | `[]` | 融资方式列表。可选值：`cash`（现金）、`credit`（贷款）、`lease`（租赁）。空数组表示不限制 |
| `sort_by` | string | 否 | 空 | 市场页面排序方式。可选值：`price_asc`、`price_desc`、`age_asc`、`age_desc` |

### PriceRange

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `min` | float | `0` | 最低价格（AS$）。`0` 表示不限制 |
| `max` | float | `0` | 最高价格（AS$）。`0` 表示不限制 |

## ActionConfig

定义匹配后的操作。

| 字段 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|--------|------|
| `auto_buy` | bool | 否 | `false` | 是否自动购买。如果为 `false`，仅发送通知不执行购买 |
| `snatch` | bool | 否 | `false` | 是否启用抢购模式。启用后使用更激进的购买策略 |
| `max_bid_increment` | float | 否 | `0` | 拍卖时最大加价幅度（AS$）。`0` 表示按当前价格出价 |

## 完整配置示例

```yaml
# 服务器配置
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
  - host: free2
    base_url: https://free2.airlinesim.aero

# 认证账户
auths:
  - username: "player1@example.com"
    password: "your_password_here"
    session_file: session1.json
  - username: "player2@example.com"
    password: "your_password_here"
    session_file: session2.json

# 监控设置
monitor:
  interval: 45           # 45 秒轮询一次
  jitter: 15             # 随机延迟 0-15 秒
  request_timeout: 30    # 请求超时 30 秒
  min_balance: 2000000   # 保留至少 200 万 AS$

# 通知设置
notifier:
  console: true
  discord_webhook: ""    # 留空表示不使用 Discord

# Web 管理界面
webui:
  enabled: true
  host: 0.0.0.0
  port: 9090

# 购买规则
rules:
  - name: "示例规则"
    enabled: true
    priority: 10
    server_id: -1        # 所有服务器
    auth_id: 0           # 使用第一个账户
    match:
      family_id: "1200300"
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
        - lease
      sort_by: price_asc
    action:
      auto_buy: false
      snatch: false
      max_bid_increment: 100000
```

## 默认配置

程序在首次运行或配置文件不存在时，会自动生成以下默认配置：

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero

auth:
  session_file: session.json

monitor:
  interval: 30
  jitter: 10
  request_timeout: 30
  min_balance: 1000000

notifier:
  console: true

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

## 已弃用字段

以下字段保留用于向后兼容，建议使用新字段替代：

| 旧字段 | 替代字段 | 说明 |
|--------|----------|------|
| `auth` | `auths` | 单账户配置，自动迁移到 `auths[0]` |
| `server` | `servers` | 单服务器配置，自动迁移到 `servers[0]` |
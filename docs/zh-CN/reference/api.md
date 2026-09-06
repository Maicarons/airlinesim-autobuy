# API 参考

Web 管理界面通过 REST API 提供完整的后端功能。所有 API 端点均以 `/api` 为前缀。

## 基础信息

- **基础路径**：`/api`
- **内容类型**：`application/json`
- **错误格式**：`{"error": "错误描述"}`

## 状态与控制

### 获取引擎状态

```
GET /api/status
```

返回引擎的当前运行状态和统计信息。

**响应示例：**

```json
{
  "running": true,
  "start_time": "2024-01-01T00:00:00Z",
  "servers": [
    {
      "host": "free1",
      "scan_count": 150,
      "found_count": 8,
      "bought_count": 3,
      "failed_count": 1,
      "last_scan": "2024-01-01T12:00:00Z",
      "last_error": ""
    }
  ],
  "scan_count": 150,
  "found_count": 8,
  "bought_count": 3,
  "failed_count": 1,
  "last_error": ""
}
```

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `running` | bool | 引擎是否正在运行 |
| `start_time` | string | 启动时间（ISO 8601） |
| `servers` | array | 各服务器的状态数组 |
| `servers[n].host` | string | 服务器标识名 |
| `servers[n].scan_count` | int | 扫描次数 |
| `servers[n].found_count` | int | 发现飞机数 |
| `servers[n].bought_count` | int | 购买成功数 |
| `servers[n].failed_count` | int | 购买失败数 |
| `servers[n].last_scan` | string | 最后扫描时间 |
| `servers[n].last_error` | string | 最后错误信息 |
| `scan_count` | int | 总扫描次数 |
| `found_count` | int | 总发现数 |
| `bought_count` | int | 总购买成功数 |
| `failed_count` | int | 总购买失败数 |
| `last_error` | string | 最后错误信息 |

### 启动引擎

```
POST /api/control/start
```

启动监控引擎。

**成功响应：**

```json
{
  "status": "started"
}
```

**错误响应（引擎已运行）：**

```json
{
  "error": "engine is already running"
}
```

### 停止引擎

```
POST /api/control/stop
```

停止监控引擎。

**响应：**

```json
{
  "status": "stopped"
}
```

### 重新加载规则

```
POST /api/reload
```

从配置文件重新加载规则，无需重启引擎。

**响应：**

```json
{
  "status": "reloaded"
}
```

## 配置管理

### 获取配置

```
GET /api/config
```

返回完整的当前配置（包括规则）。

**响应：** 完整的 `Config` 对象。

### 更新配置

```
PUT /api/config
```

更新全局配置（不包括规则）。只更新请求中提供的字段，未提供的字段保持不变。

**请求体示例：**

```json
{
  "servers": [
    {"host": "free1", "base_url": "https://free1.airlinesim.aero"}
  ],
  "monitor": {
    "interval": 60,
    "jitter": 20
  },
  "notifier": {
    "console": true,
    "discord_webhook": "https://discord.com/api/webhooks/..."
  },
  "webui": {
    "host": "0.0.0.0",
    "port": 9090
  }
}
```

## 规则管理

### 获取规则列表

```
GET /api/rules
```

返回所有规则的数组。

**响应：** `RuleConfig[]`

### 创建规则

```
POST /api/rules
```

创建一条新规则。

**请求体：**

```json
{
  "name": "新规则",
  "enabled": true,
  "priority": 10,
  "server_id": 0,
  "auth_id": 0,
  "match": {
    "types": ["Airbus A320-200 heavy"],
    "price_range": {"min": 500000, "max": 5000000},
    "max_age": 15,
    "max_cycles": 30000,
    "condition_min": 70,
    "offer_types": ["auction", "immediate"],
    "financing": ["cash", "credit", "lease"]
  },
  "action": {
    "auto_buy": true,
    "snatch": false,
    "max_bid_increment": 100000
  }
}
```

**状态码：** `201 Created`

### 获取单条规则

```
GET /api/rules/{id}
```

按索引获取规则。`{id}` 是规则的数组索引（从 0 开始）。

**成功响应：** `RuleConfig` 对象

**错误响应：**

```json
{
  "error": "rule not found"
}
```

### 更新规则

```
PUT /api/rules/{id}
```

按索引更新规则。

**请求体：** 完整的 `RuleConfig` 对象

**状态码：** `200 OK`

### 删除规则

```
DELETE /api/rules/{id}
```

按索引删除规则。

**响应：**

```json
{
  "status": "deleted"
}
```

### 切换规则启用状态

```
PATCH /api/rules/{id}/toggle
```

切换指定规则的启用/禁用状态。

**响应：** 更新后的 `RuleConfig` 对象

### 重新排序规则

```
PUT /api/rules/reorder
```

重新排序规则列表。

**请求体：** 一个整数数组，表示新的索引顺序

```json
[2, 0, 1]
```

## 飞机数据

### 获取飞机市场参考数据

```
GET /api/aircraft-data
```

返回飞机系列和型号的参考数据，用于 Web 界面的下拉选择。

**响应示例：**

```json
{
  "families": [
    {"id": "1200300", "name": "A320 / A321"},
    {"id": "2200400", "name": "737-600/700/800/900"}
  ],
  "types": [
    {"id": "16", "name": "Airbus A320-200 heavy"},
    {"id": "97", "name": "Boeing 737-800 BGW"}
  ]
}
```

## 通用错误

| 状态码 | 说明 |
|--------|------|
| `400 Bad Request` | 请求体格式错误或参数无效 |
| `404 Not Found` | 请求的资源（如规则）不存在 |
| `409 Conflict` | 状态冲突（如引擎已在运行） |
| `500 Internal Server Error` | 服务器内部错误 |

## 数据模型

### Config

| 字段 | 类型 | 说明 |
|------|------|------|
| `servers` | ServerConfig[] | 服务器列表 |
| `auths` | AuthConfig[] | 认证账户列表 |
| `monitor` | MonitorConfig | 监控设置 |
| `notifier` | NotifierConfig | 通知设置 |
| `webui` | WebUIConfig | WebUI 设置 |
| `rules` | RuleConfig[] | 规则列表 |

### ServerConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| `host` | string | 服务器标识名 |
| `base_url` | string | 服务器基础 URL |

### AuthConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| `username` | string | 用户名 |
| `password` | string | 密码 |
| `session_file` | string | 会话文件路径 |

### MonitorConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| `interval` | int | 轮询间隔（秒） |
| `jitter` | int | 随机延迟（秒） |
| `request_timeout` | int | 请求超时（秒） |
| `min_balance` | float | 最低余额保护 |

### NotifierConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| `console` | bool | 控制台输出 |
| `discord_webhook` | string | Discord Webhook URL |

### WebUIConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| `enabled` | bool | 是否启用 |
| `host` | string | 监听地址 |
| `port` | int | 监听端口 |

### RuleConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 规则名称 |
| `enabled` | bool | 是否启用 |
| `priority` | int | 优先级 |
| `server_id` | int | 目标服务器索引（-1 表示全部）|
| `auth_id` | int | 执行账户索引 |
| `match` | MatchConfig | 匹配条件 |
| `action` | ActionConfig | 操作设置 |

### MatchConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| `family_id` | string | 飞机系列 ID |
| `type_id` | string | 飞机型号 ID |
| `types` | string[] | 型号名称列表 |
| `price_range` | PriceRange | 价格范围 |
| `max_age` | int | 最大机龄（年） |
| `max_cycles` | int | 最大飞行循环 |
| `condition_min` | float | 最低状态（%） |
| `offer_types` | string[] | 交易类型 |
| `financing` | string[] | 融资方式 |
| `sort_by` | string | 排序方式 |

### PriceRange

| 字段 | 类型 | 说明 |
|------|------|------|
| `min` | float | 最低价格 |
| `max` | float | 最高价格 |

### ActionConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| `auto_buy` | bool | 是否自动购买 |
| `snatch` | bool | 是否启用抢购 |
| `max_bid_increment` | float | 最大加价幅度 |
# 购买规则

购买规则是 AirlineSim Autobuy 的核心功能。通过配置规则，你可以精确控制监控哪些飞机、在什么条件下购买、使用哪个账户进行交易。

## 规则结构

每条规则包含以下字段：

```yaml
rules:
  - name: "规则名称"           # 规则标识名
    enabled: true              # 是否启用
    priority: 10               # 优先级（数值越大优先级越高）
    server_id: 0               # 目标服务器索引
    auth_id: 0                 # 执行账户索引
    match:                     # 匹配条件
      # ... 详见下文
    action:                    # 操作设置
      # ... 详见下文
```

### 基本字段

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `name` | string | 必填 | 规则名称，用于日志和通知标识 |
| `enabled` | bool | `true` | 是否启用该规则，禁用后规则不参与匹配 |
| `priority` | int | `10` | 优先级分数，作为评分的一部分，优先级越高越优先匹配 |

## 匹配条件

`match` 字段定义了一架飞机需要满足的条件才能匹配该规则。

### 完整示例

```yaml
match:
  family_id: "1200300"          # 飞机系列 ID
  type_id: ""                   # 飞机具体型号 ID
  types:                        # 飞机型号名称列表
    - "Airbus A320-200 heavy"
    - "Boeing 737-800 HGW"
  price_range:                  # 价格范围（AS$）
    min: 500000
    max: 5000000
  max_age: 15                   # 最大机龄（年）
  max_cycles: 30000             # 最大飞行循环数
  condition_min: 70             # 最低状态（%）
  offer_types:                  # 交易类型
    - auction                   # 拍卖
    - immediate                 # 立即购买
  financing:                    # 融资方式
    - cash                      # 现金
    - credit                    # 贷款
    - lease                     # 租赁
  sort_by: price_asc            # 市场排序方式
```

### 字段说明

#### 飞机系列（family_id）

通过 Wicket 系列 ID 筛选飞机系列。例如 `"1200300"` 对应 A320/A321 系列。可在 Web 界面的"飞机数据"下拉菜单中查看所有可用系列。

常见的系列 ID：

| 系列 ID | 名称 |
|---------|------|
| `1200300` | A320 / A321 |
| `1200310` | A319 / A320 / A321 NEO |
| `2200400` | 737-600/700/800/900 |
| `2200410` | 737 MAX |
| `1200400` | A330 |
| `1200600` | A350 |
| `4200300` | EMB 170/175/190/195 |

#### 飞机型号（types）

通过飞机型号名称进行匹配，支持大小写不敏感匹配。例如：

```yaml
types:
  - "Airbus A320-200 heavy"
  - "Boeing 737-800 HGW"
  - "Embraer 190"
```

#### 价格范围（price_range）

设置价格上下限，单位为 AS$：

```yaml
price_range:
  min: 500000     # 最低价格
  max: 5000000    # 最高价格
```

将 `min` 或 `max` 设为 `0` 表示不限制。

#### 机龄（max_age）

最大机龄，单位为年。例如 `max_age: 15` 表示只匹配机龄不超过 15 年的飞机。

#### 飞行循环（max_cycles）

最大飞行循环数。飞行循环（Cycles）是飞机起降次数，通常与机身疲劳相关。

#### 状态（condition_min）

最低机身状态百分比（0-100）。例如 `condition_min: 70` 表示只匹配状态不低于 70% 的飞机。

#### 交易类型（offer_types）

限制匹配的交易类型：

| 值 | 说明 |
|----|------|
| `auction` | 拍卖——需要竞拍出价 |
| `immediate` | 立即购买——可直接按标价购买 |

#### 融资方式（financing）

限制可接受的融资方式：

| 值 | 说明 |
|----|------|
| `cash` | 现金全额支付 |
| `credit` | 贷款分期 |
| `lease` | 租赁（仅需支付押金和租金） |

#### 排序方式（sort_by）

设置市场页面的排序方式：

| 值 | 说明 |
|----|------|
| `price_asc` | 价格从低到高 |
| `price_desc` | 价格从高到低 |
| `age_asc` | 机龄从低到高 |
| `age_desc` | 机龄从高到低 |

## 操作设置

`action` 字段定义当飞机匹配规则时要执行的操作。

```yaml
action:
  auto_buy: true           # 是否自动购买
  snatch: false            # 是否启用抢购模式
  max_bid_increment: 100000  # 最大加价幅度（AS$）
```

### 字段说明

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `auto_buy` | bool | `false` | 是否自动执行购买。如果为 `false`，仅通知不执行购买 |
| `snatch` | bool | `false` | 抢购模式，启用后以更激进的策略快速下单 |
| `max_bid_increment` | float | `0` | 拍卖时最大加价幅度（AS$），`0` 表示按当前价格出价 |

### 自动购买流程

当 `auto_buy: true` 时，匹配成功后的执行流程：

1. 检查规则评分是否大于 0
2. 如果是立即购买（`immediate`）：访问飞机详情页，提交购买表单
3. 如果是拍卖（`auction`）：提交出价，出价金额为当前价格 + `max_bid_increment`
4. 购买成功后更新统计数据并发送通知
5. 购买失败时记录失败原因并发送通知

## 服务器选择

通过 `server_id` 字段将规则绑定到特定的游戏服务器。

```yaml
server_id: 0    # 索引 0 对应 servers 列表中的第一个服务器
```

- `server_id` 是 `servers` 数组的索引（从 0 开始）
- 设置为 `-1` 表示匹配所有服务器
- 例如，如果配置了 3 个服务器，`server_id: 1` 表示第二个服务器

## 账户选择

通过 `auth_id` 字段指定执行该规则时使用的认证账户。

```yaml
auth_id: 0    # 索引 0 对应 auths 列表中的第一个账户
```

- `auth_id` 是 `auths` 数组的索引（从 0 开始）
- 设置为 `-1` 使用默认账户（第一个账户）
- 多账户场景下，不同规则可以使用不同的账户执行购买

## 规则评分机制

当一架飞机匹配多条规则时，系统会计算每条规则的评分，并按评分降序选择最优匹配。

评分由以下因素组成：

| 因素 | 权重 | 说明 |
|------|------|------|
| 价格优势 | 最高 50 分 | 价格越低评分越高，计算公式：`(1 - 价格/最高价) * 50` |
| 机身状态 | 最高 20 分 | 状态越好评分越高，超出最低状态越多分越高 |
| 机龄优势 | 最高 20 分 | 机龄越短评分越高，计算公式：`(1 - 机龄/最大机龄) * 20` |
| 立即购买 | 10 分 | 立即购买类型获得额外加分 |
| 规则优先级 | 按设定值 | 规则配置的 `priority` 值直接加入评分 |

## 规则示例

### 示例 1：经济型窄体机自动购买

```yaml
- name: "经济型窄体机自动购买"
  enabled: true
  priority: 10
  server_id: -1
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
    max_bid_increment: 50000
```

### 示例 2：宽体机监控（仅通知不购买）

```yaml
- name: "宽体机监控"
  enabled: true
  priority: 5
  server_id: 0
  auth_id: 0
  match:
    types:
      - "Airbus A330-300"
      - "Boeing 787-9"
    price_range:
      max: 25000000
    max_age: 20
    condition_min: 60
  action:
    auto_buy: false
```

### 示例 3：租赁机会监控

```yaml
- name: "租赁机会"
  enabled: true
  priority: 8
  server_id: -1
  auth_id: 0
  match:
    family_id: "1200300"
    price_range:
      max: 1000000
    max_age: 25
    financing:
      - lease
  action:
    auto_buy: true
    max_bid_increment: 100000
```

## 最佳实践

1. **从仅通知开始**：首次使用将 `auto_buy` 设为 `false`，观察匹配结果是否准确
2. **合理设置优先级**：高优先级规则用于首选机型，低优先级规则作为后备
3. **设置余额保护**：通过 `min_balance` 防止因购买导致账户余额过低
4. **分散轮询间隔**：多服务器场景下，设置合理的 `jitter` 避免请求集中
5. **定期检查日志**：通过日志确认规则匹配是否符合预期
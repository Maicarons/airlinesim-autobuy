# 快速开始

本指南将帮助你从零开始安装、配置并运行 AirlineSim Autobuy。

## 前置条件

在开始之前，请确保你的系统满足以下要求：

- **Go 1.22+** — 用于编译后端二进制文件
- **Node.js 18+** — 用于构建前端界面（仅开发时需要）
- **一个 AirlineSim 游戏账户** — 用于登录游戏服务器

## 安装

### 克隆仓库

```bash
git clone https://github.com/Maicarons/airlinesim-autobuy.git
cd airlinesim-autobuy
```

### 安装前端依赖

```bash
make frontend-install
```

该命令会进入 `internal/webui/frontend` 目录并执行 `npm install`。

### 构建前端

```bash
make frontend-build
```

该命令将 Vue 3 前端编译为静态文件，并通过 Go 的 `embed` 打包到二进制文件中。

### 构建后端

```bash
make build
```

该命令会生成一个名为 `autobuy`（Windows 下为 `autobuy.exe`）的单二进制文件，包含所有后端逻辑和嵌入的前端资源。

### 一键构建

你也可以使用以下命令一次性完成前端安装、前端构建和后端构建：

```bash
make all
```

## 跨平台构建

Makefile 提供了跨平台构建支持：

```bash
make build-all
```

该命令会生成以下平台的二进制文件：

| 平台 | 文件名 |
|------|--------|
| Windows AMD64 | `autobuy-windows-amd64.exe` |
| Linux AMD64 | `autobuy-linux-amd64` |
| macOS AMD64 | `autobuy-darwin-amd64` |

## 配置

### 创建配置文件

首次运行时会自动生成默认配置文件 `configs/config.yaml`。你也可以手动创建：

```bash
cp configs/config.yaml configs/config.local.yaml
```

### 编辑配置

打开 `configs/config.yaml`，填入你的 AirlineSim 账户信息：

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero

auth:
  username: "your_username"
  password: "your_password"
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
  - name: "示例规则"
    enabled: true
    priority: 10
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
    action:
      auto_buy: false
      snatch: false
      max_bid_increment: 100000
```

::: warning 安全提示
请勿将包含密码的配置文件提交到版本控制系统中。建议将 `configs/config.yaml` 加入 `.gitignore` 或使用环境变量管理敏感信息。
:::

## 运行

### 直接运行

```bash
./autobuy
```

### 使用 Make 命令

```bash
make run
```

### 开发模式（热重载）

如果安装了 [Air](https://github.com/air-verse/air) 热重载工具：

```bash
make dev
```

## 首次启动

启动后，程序会依次执行以下操作：

1. 加载配置文件 `configs/config.yaml`
2. 初始化认证会话，登录到指定的游戏服务器
3. 发现飞机市场页面的 URL
4. 启动 Web 管理界面（默认端口 9090）
5. 开始按设定的间隔轮询市场数据
6. 对每个发现的飞机，使用规则引擎进行匹配评估
7. 如果匹配且启用了自动购买，执行购买操作

## 访问 Web 界面

打开浏览器访问 `http://localhost:9090`，你将看到：

- **仪表盘** — 引擎运行状态、扫描统计、服务器状态
- **规则管理** — 创建、编辑、启用/禁用购买规则
- **设置** — 修改全局配置（监控、通知、WebUI 等）
- **日志** — 实时查看操作日志和事件

## 下一步

- 阅读[配置说明](./configuration)了解完整的配置项
- 阅读[购买规则](./rules)学习如何编写有效的匹配规则
- 阅读[部署指南](./deployment)了解生产环境部署方案
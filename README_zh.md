<h1 align="center">✈️ AirlineSim Autobuy</h1>
<p align="center">
  <em>AirlineSim 二手飞机市场自动监控与购买工具</em>
</p>

<p align="center">
  <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go" alt="Go"></a>
  <a href="https://vuejs.org"><img src="https://img.shields.io/badge/Vue-3.4-4FC08D?logo=vue.js" alt="Vue"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-AGPL--3.0-blue" alt="License"></a>
</p>

<p align="center">
  <a href="#功能">功能</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#配置说明">配置说明</a> •
  <a href="#部署">部署</a> •
  <a href="docs/zh-CN/guide/quickstart.md">完整文档</a>
</p>

---

## 功能

- **🔍 实时监控** — 持续扫描配置服务器上的二手飞机市场，即时发现新上架
- **📋 智能规则引擎** — 按航机系列、类型、价格、机龄、状态、循环、租赁费率定义购买规则
- **🤖 自动购买** — 支持立即购买和竞拍出价，可配置抢购和余额保护
- **🌐 多服务器多账户** — 同时监控多个游戏服务器，每个服务器可独立配置认证账户
- **⚙️ Web 管理界面** — 功能完整的仪表盘，支持引擎控制、规则管理、实时日志和设置
- **🌍 多语言支持** — English、中文、한국어
- **🚀 单二进制部署** — Go 编译为单二进制文件，内嵌 Vue 前端，零依赖部署

## 快速开始

### 前置条件

- Go 1.22+
- Node.js 18+（前端开发可选）
- AirlineSim 游戏账户

### 安装

```bash
# 克隆仓库
git clone https://github.com/Maicarons/airlinesim-autobuy.git
cd airlinesim-autobuy

# 安装前端依赖并构建
make frontend-install
make frontend-build

# 构建 Go 后端
make build

# 复制并编辑配置文件
cp configs/config.example.yaml configs/config.yaml
# 编辑 configs/config.yaml 填入你的账户信息
```

### 运行

```bash
./autobuy
```

在浏览器中打开 http://localhost:9090 访问 Web 管理界面。

## 配置说明

编辑 `configs/config.yaml`：

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
auths:
  - username: "your_email@example.com"
    password: "your_password"
monitor:
  interval: 30          # 轮询间隔（秒）
  min_balance: 1000000  # 购买后保留的最低余额
```

### 购买规则

```yaml
rules:
  - name: "监控 A320-200 Heavy"
    enabled: true
    server_id: 0
    auth_id: 0
    match:
      family_id: "1200300"       # A320 / A321 系列
      type_id: "16"              # Airbus A320-200 heavy
      price_range:
        min: 0
        max: 5000000
      max_age: 20
      condition_min: 50
    action:
      auto_buy: false
      snatch: false
```

## 项目结构

```
airlinesim-autobuy/
├── cmd/autobuy/          # 程序入口
├── internal/
│   ├── auth/             # 认证与会话管理
│   ├── client/           # 限速 HTTP 客户端
│   ├── collector/        # 市场页面采集
│   ├── config/           # 配置管理
│   ├── engine/           # 管道编排引擎
│   ├── executor/         # 购买执行
│   ├── marketdata/       # 飞机参考数据
│   ├── notifier/         # 通知系统
│   ├── parser/           # HTML 解析器
│   ├── rules/            # 规则引擎
│   └── webui/            # Web 管理界面
├── configs/              # 配置文件
├── docs/                 # VitePress 文档
├── Makefile
└── README.md
```

## 文档

完整文档在 [docs](docs/) 目录：

| 语言 | 用户指南 | 开发者文档 |
|------|---------|-----------|
| English | [Quick Start](docs/en/guide/quickstart.md) | [Architecture](docs/en/reference/architecture.md) |
| 中文 | [快速开始](docs/zh-CN/guide/quickstart.md) | [架构说明](docs/zh-CN/reference/architecture.md) |
| 한국어 | [빠른 시작](docs/ko/guide/quickstart.md) | [아키텍처](docs/ko/reference/architecture.md) |

## 部署

### Docker

```bash
docker build -t airlinesim-autobuy .
docker run -d -p 9090:9090 -v ./configs:/app/configs airlinesim-autobuy
```

### systemd

```bash
sudo cp airlinesim-autobuy.service /etc/systemd/system/
sudo systemctl enable airlinesim-autobuy
sudo systemctl start airlinesim-autobuy
```

## 许可证

[AGPL-3.0](LICENSE)

## 免责声明

本工具仅用于个人学习用途。自动化操作可能违反游戏服务条款，请自行承担风险。请合理配置轮询间隔，尊重游戏服务器负载。
<h1 align="center">✈️ AirlineSim Autobuy</h1>
<p align="center">
  <em>Automated Aircraft Market Monitor & Purchaser for AirlineSim</em>
</p>

<p align="center">
  <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go" alt="Go"></a>
  <a href="https://vuejs.org"><img src="https://img.shields.io/badge/Vue-3.4-4FC08D?logo=vue.js" alt="Vue"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue" alt="License"></a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#deployment">Deployment</a> •
  <a href="docs/en/guide/quickstart.md">Documentation</a>
</p>

---

## Features

- **🔍 Real-time Monitoring** — Continuously scans the used aircraft market on configured servers, instantly detecting new listings
- **📋 Smart Rule Engine** — Define purchase rules by aircraft family, type, price, age, condition, cycles, and lease rate
- **🤖 Automated Purchasing** — Supports immediate purchase and auction bidding with configurable snatch and balance protection
- **🌐 Multi-Server & Multi-Account** — Monitor multiple game servers simultaneously with independent authentication accounts
- **⚙️ Web Management UI** — Full-featured dashboard for engine control, rule management, real-time logs, and settings
- **🌍 i18n Support** — English, 中文, 한국어
- **🚀 Single Binary** — Go-compiled single binary with embedded Vue frontend — deploy anywhere with zero dependencies

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+ (for frontend development, optional)
- An AirlineSim game account

### Installation

```bash
# Clone the repository
git clone https://github.com/Maicarons/airlinesim-autobuy.git
cd airlinesim-autobuy

# Install frontend dependencies and build
make frontend-install
make frontend-build

# Build the Go binary
make build

# Copy and edit the configuration
cp configs/config.example.yaml configs/config.yaml
# Edit configs/config.yaml with your credentials
```

### Running

```bash
./autobuy
```

Open http://localhost:9090 in your browser to access the Web UI.

## Configuration

Edit `configs/config.yaml`:

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
auths:
  - username: "your_email@example.com"
    password: "your_password"
monitor:
  interval: 30       # Polling interval in seconds
  min_balance: 1000000  # Minimum balance to retain after purchase
```

### Purchase Rules

```yaml
rules:
  - name: "A320-200 Heavy Monitor"
    enabled: true
    server_id: 0
    auth_id: 0
    match:
      family_id: "1200300"       # A320 / A321 family
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

## Project Structure

```
airlinesim-autobuy/
├── cmd/autobuy/          # Entry point
├── internal/
│   ├── auth/             # Authentication & session management
│   ├── client/           # Rate-limited HTTP client
│   ├── collector/        # Market page fetcher
│   ├── config/           # Configuration management
│   ├── engine/           # Pipeline orchestrator
│   ├── executor/         # Purchase execution
│   ├── marketdata/       # Aircraft reference data
│   ├── notifier/         # Notification system
│   ├── parser/           # HTML parser
│   ├── rules/            # Rules engine
│   └── webui/            # Web management UI
├── configs/              # Configuration files
├── docs/                 # VitePress documentation
├── Makefile
└── README.md
```

## Documentation

Full documentation is available in the [docs](docs/) directory:

| Language | User Guide | Developer Guide |
|----------|-----------|-----------------|
| English | [Quick Start](docs/en/guide/quickstart.md) | [Architecture](docs/en/reference/architecture.md) |
| 中文 | [快速开始](docs/zh-CN/guide/quickstart.md) | [架构说明](docs/zh-CN/reference/architecture.md) |
| 한국어 | [빠른 시작](docs/ko/guide/quickstart.md) | [아키텍처](docs/ko/reference/architecture.md) |

## Deployment

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

## License

[MIT](LICENSE)

## Disclaimer

This tool is for personal educational use. Automated gameplay may violate the game's Terms of Service. Use at your own risk. Please be respectful of the game servers and configure reasonable polling intervals.
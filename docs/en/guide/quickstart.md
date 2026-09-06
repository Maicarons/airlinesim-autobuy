# Quick Start

Get AirlineSim Autobuy up and running in minutes. This guide covers everything from prerequisites to your first automated scan.

## Prerequisites

Before you begin, ensure you have the following installed:

| Requirement | Version | Purpose |
|---|---|---|
| [Go](https://golang.org/dl/) | 1.22+ | Build the backend binary |
| [Node.js](https://nodejs.org/) | 18+ | Build the Vue 3 frontend |
| [Git](https://git-scm.com/) | Any | Clone the repository |
| Make | Any | Use build targets (optional) |

You will also need:

- An **AirlineSim game account** with sufficient balance for purchases
- The **game server URL** you play on (e.g., `https://free1.airlinesim.aero`)

## Installation

### Clone the Repository

```bash
git clone https://github.com/Maicarons/airlinesim-autobuy.git
cd airlinesim-autobuy
```

### Install Frontend Dependencies

The embedded web UI is built with Vue 3 and Vite. Install its dependencies:

```bash
make frontend-install
```

Or manually:

```bash
cd internal/webui/frontend
npm install
cd ../../..
```

### Build the Frontend

```bash
make frontend-build
```

This produces the production bundle at `internal/webui/frontend/dist/`, which is embedded into the Go binary at compile time.

### Build the Backend

```bash
make build
```

This compiles a single standalone binary named `autobuy` (or `autobuy.exe` on Windows).

To build for all platforms at once:

```bash
make build-all
```

The resulting binaries are placed in the project root with platform-specific suffixes.

### Full Build (Frontend + Backend)

If you prefer a single command that does everything:

```bash
make all
```

This runs `frontend-install`, `frontend-build`, and `build` in sequence.

## Configuration

### Create Your Configuration File

Copy the default configuration as a starting point:

```bash
cp configs/config.yaml configs/config.local.yaml
```

Edit the configuration file with your AirlineSim credentials:

```yaml
auth:
  username: "your_username@email.com"
  password: "your_password"
  session_file: "session.json"

monitor:
  interval: 30
  jitter: 10
  min_balance: 1000000

rules:
  - name: "My First Rule"
    enabled: true
    match:
      types:
        - "Airbus A320-200 heavy"
      price_range:
        min: 500000
        max: 5000000
      max_age: 15
      max_cycles: 30000
      condition_min: 70
    action:
      auto_buy: false
      max_bid_increment: 100000
```

> **Important:** Set `auto_buy: false` for your first rule to test monitoring without making purchases. Enable it only after you are confident the rules match correctly.

### Configuration File Location

By default, the application looks for `configs/config.yaml`. You can change this by editing the `configPath` constant in `cmd/autobuy/main.go` or by using the Web UI settings page after startup.

## Running

### Start the Application

```bash
./autobuy
```

Or using the Make target:

```bash
make run
```

On startup, you should see output similar to:

```
INFO starting airlinesim-autobuy
INFO configuration loaded path=configs/config.yaml servers=2 auths=1
INFO engine initialized servers=2
INFO starting web UI host=0.0.0.0 port=9090
INFO starting server monitor host=free1
INFO attempting to log in server=https://free1.airlinesim.aero
INFO login successful, session established
INFO discovering aircraft market URL
INFO starting server monitor host=free2
...
```

### Access the Web UI

Open your browser and navigate to [http://localhost:9090](http://localhost:9090).

You will see the dashboard showing:

- Engine status (running / stopped)
- Per-server scan statistics
- Found and purchased aircraft counts
- Recent log entries

### Development Mode (Hot Reload)

For development with automatic restart on file changes, install [Air](https://github.com/air-verse/air) and run:

```bash
make dev
```

This watches for `.go` file changes and restarts the application automatically.

## First Steps

### 1. Verify Authentication

Check the startup logs for a successful authentication message:

```
INFO login successful, session established
```

If authentication fails, verify your username and password in the config file.

### 2. Monitor the Dashboard

Open the Web UI at `http://localhost:9090`. The dashboard shows:

- **Engine Status** — whether the monitoring engine is running
- **Server Cards** — one per configured game server, showing scan count, found aircraft, and last scan time
- **Log Feed** — real-time streaming log of events

### 3. Review Matched Aircraft

When the engine finds aircraft matching your rules, they appear in the dashboard log. The notifier also prints details to the console:

```
INFO 🔍 Aircraft found matching rule 'My First Rule': Airbus A320-200 heavy
    details="Type: Airbus A320-200 heavy | Price: AS$ 2500000 | Age: 8y | Cond: 85% | Type: auction"
```

### 4. Enable Auto-Buy

Once you are satisfied with the matching behavior, edit the rule and set `auto_buy: true`. You can do this either:

- **In the config file** — change `auto_buy: false` to `auto_buy: true` and save (the file watcher reloads automatically)
- **In the Web UI** — navigate to the Rules page, find your rule, and toggle the Auto-Buy switch

### 5. Monitor Purchases

When a purchase is made, the console shows:

```
INFO 🛒 PURCHASED Airbus A320-200 heavy for AS$ 2500000 (rule: My First Rule)
```

## Stopping the Application

Press `Ctrl+C` in the terminal. The application performs a graceful shutdown:

```
^C
INFO received signal signal=interrupt
INFO engine stopped
INFO shutdown complete
```

## Next Steps

- Learn about [Configuration options](/guide/configuration) for multi-server and multi-account setups
- Explore [Purchase Rules](/guide/rules) in detail to fine-tune matching
- Read the [Deployment Guide](/guide/deployment) for production setups
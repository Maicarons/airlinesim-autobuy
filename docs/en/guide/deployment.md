# Deployment

This guide covers production deployment scenarios for AirlineSim Autobuy, including systemd service, Docker, and Docker Compose setups.

## Building for Production

### Cross-Platform Build

The Makefile provides a convenient target for building all platforms:

```bash
make build-all
```

This produces:

| File | Platform |
|---|---|
| `autobuy-windows-amd64.exe` | Windows x86_64 |
| `autobuy-linux-amd64` | Linux x86_64 |
| `autobuy-darwin-amd64` | macOS x86_64 |

### Single-Platform Build

For the current platform:

```bash
make build
```

### Full Build with Frontend

```bash
make all
```

This installs frontend dependencies, builds the Vue app, and compiles the Go binary.

## Systemd Service (Linux)

Deploy as a systemd service for automatic startup and process management.

### 1. Create a Dedicated User

```bash
sudo useradd -r -s /bin/false autobuy
```

### 2. Place Files

```bash
sudo mkdir -p /opt/autobuy
sudo cp autobuy-linux-amd64 /opt/autobuy/autobuy
sudo cp -r configs /opt/autobuy/
sudo chown -R autobuy:autobuy /opt/autobuy
sudo chmod +x /opt/autobuy/autobuy
```

### 3. Create the Service File

Create `/etc/systemd/system/autobuy.service`:

```ini
[Unit]
description=AirlineSim Autobuy - Automated Aircraft Market Monitor
after=network.target
wants=network.target

[Service]
type=simple
user=autobuy
group=autobuy
workingdirectory=/opt/autobuy
execstart=/opt/autobuy/autobuy
restart=on-failure
restartsec=10
standardoutput=journal
standarderror=journal

# Security hardening
noprivileges=true
nonmountutils=true
capabilityboundingset=
protectsystem=strict
protecthome=true
privateTmp=true

[Install]
wantedby=multi-user.target
```

### 4. Enable and Start

```bash
sudo systemctl daemon-reload
sudo systemctl enable autobuy
sudo systemctl start autobuy
```

### 5. Verify

```bash
sudo systemctl status autobuy
journalctl -u autobuy -f
```

### 6. Managing the Service

```bash
sudo systemctl restart autobuy    # Restart
sudo systemctl stop autobuy       # Stop
sudo systemctl start autobuy      # Start
```

## Docker

### Dockerfile

Create a `Dockerfile` in the project root:

```dockerfile
# Stage 1: Build frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY internal/webui/frontend/package*.json ./
RUN npm install
COPY internal/webui/frontend/ .
RUN npm run build

# Stage 2: Build backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/dist ./internal/webui/frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o autobuy ./cmd/autobuy

# Stage 3: Runtime
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /app/autobuy .
COPY configs/config.yaml ./configs/
EXPOSE 9090
VOLUME ["/app/configs", "/app/sessions"]
ENTRYPOINT ["./autobuy"]
```

### Build the Docker Image

```bash
docker build -t airlinesim-autobuy:latest .
```

### Run the Container

```bash
docker run -d \
  --name autobuy \
  -p 9090:9090 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/sessions:/app/sessions \
  airlinesim-autobuy:latest
```

### Docker Run Options

| Option | Description |
|---|---|
| `-d` | Run in detached mode |
| `--name autobuy` | Container name |
| `-p 9090:9090` | Map web UI port |
| `-v ./configs:/app/configs` | Mount configuration directory |
| `-v ./sessions:/app/sessions` | Mount session persistence directory |
| `--restart unless-stopped` | Auto-restart policy |

## Docker Compose

### docker-compose.yml

Create a `docker-compose.yml` file:

```yaml
version: "3.8"

services:
  autobuy:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: airlinesim-autobuy
    restart: unless-stopped
    ports:
      - "9090:9090"
    volumes:
      - ./configs:/app/configs
      - ./sessions:/app/sessions
      - ./logs:/app/logs
    environment:
      - TZ=UTC
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Start with Docker Compose

```bash
docker-compose up -d
```

### View Logs

```bash
docker-compose logs -f
```

### Stop

```bash
docker-compose down
```

## Reverse Proxy Setup

The web UI does not include built-in authentication. For production deployments, place a reverse proxy in front of it.

### Nginx

```nginx
server {
    listen 443 ssl;
    server_name autobuy.example.com;

    ssl_certificate /etc/ssl/certs/autobuy.crt;
    ssl_certificate_key /etc/ssl/private/autobuy.key;

    location / {
        proxy_pass http://127.0.0.1:9090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_buffering off;
        proxy_cache off;
    }
}
```

### Caddy

```caddyfile
autobuy.example.com {
    reverse_proxy 127.0.0.1:9090
}
```

Caddy automatically handles TLS certificate provisioning.

### Basic Authentication (Nginx)

```nginx
location / {
    auth_basic "Restricted";
    auth_basic_user_file /etc/nginx/.htpasswd;
    proxy_pass http://127.0.0.1:9090;
}
```

Create the `.htpasswd` file:

```bash
sudo htpasswd -c /etc/nginx/.htpasswd admin
```

## Monitoring

### Health Check

The web UI's `/api/status` endpoint can be used for health checks:

```bash
curl http://localhost:9090/api/status
```

A healthy response includes `"running": true` or `"running": false` and server statistics.

### Prometheus Integration

The application does not expose Prometheus metrics natively. Use a Prometheus Blackbox Exporter to monitor the HTTP health endpoint:

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'autobuy'
    metrics_path: /probe
    params:
      module: [http_2xx]
    static_configs:
      - targets:
        - http://localhost:9090/api/status
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: blackbox-exporter:9115
```

### Log Monitoring

Systemd journal logs can be forwarded to an external log aggregation system:

```bash
journalctl -u autobuy -f -o json
```

For Docker, use the built-in logging drivers:

```yaml
logging:
  driver: "fluentd"
  options:
    fluentd-address: "localhost:24224"
    tag: "autobuy"
```

## Performance Considerations

### Polling Interval

The `interval` setting controls how often the market is scanned. Recommended values:

| Scenario | Interval | Jitter |
|---|---|---|
| Development / Testing | 60s | 10s |
| Active monitoring | 30s | 10s |
| Aggressive monitoring | 15s | 5s |

> **Note:** Shorter intervals increase load on the game server. Use reasonable values to avoid rate limiting.

### Memory Usage

The application is lightweight, typically using 20-50 MB of RAM. The embedded web UI adds minimal overhead.

### File Descriptors

Each server connection uses one HTTP client. The default configuration with 2 servers uses approximately 10-15 file descriptors.

## Upgrading

### Binary Upgrade

1. Stop the service:
   ```bash
   sudo systemctl stop autobuy
   ```

2. Replace the binary:
   ```bash
   sudo cp autobuy-linux-amd64 /opt/autobuy/autobuy
   sudo chmod +x /opt/autobuy/autobuy
   ```

3. Start the service:
   ```bash
   sudo systemctl start autobuy
   ```

### Docker Upgrade

```bash
docker-compose pull
docker-compose up -d
```

## See Also

- [Quick Start](/guide/quickstart) for initial setup
- [Configuration Guide](/guide/configuration) for production configuration
- [FAQ](/guide/faq) for troubleshooting common issues
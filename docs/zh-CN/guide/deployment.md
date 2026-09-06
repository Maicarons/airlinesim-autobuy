# 部署指南

本指南介绍如何将 AirlineSim Autobuy 部署到生产环境。

## 编译二进制文件

在部署前，需要为目标平台编译二进制文件：

```bash
# 为 Linux AMD64 编译
make build-all

# 或单独为 Linux 编译
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o autobuy-linux-amd64 ./cmd/autobuy
```

编译完成后，将二进制文件和配置文件一同传输到目标服务器。

## systemd 部署

对于 Linux 系统，推荐使用 systemd 管理服务。

### 创建服务文件

创建 `/etc/systemd/system/airlinesim-autobuy.service`：

```ini
[Unit]
Description=AirlineSim Autobuy - Aircraft Market Monitor
After=network.target

[Service]
Type=simple
User=autobuy
Group=autobuy
WorkingDirectory=/opt/airlinesim-autobuy
ExecStart=/opt/airlinesim-autobuy/autobuy
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal
Environment=GIN_MODE=release

# 安全限制
NoNewPrivileges=true
ProtectHome=true
ProtectSystem=full
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

### 部署步骤

```bash
# 创建用户和目录
sudo useradd -r -s /bin/false autobuy
sudo mkdir -p /opt/airlinesim-autobuy/configs

# 复制文件
sudo cp autobuy-linux-amd64 /opt/airlinesim-autobuy/autobuy
sudo cp configs/config.yaml /opt/airlinesim-autobuy/configs/config.yaml

# 设置权限
sudo chown -R autobuy:autobuy /opt/airlinesim-autobuy
sudo chmod +x /opt/airlinesim-autobuy/autobuy

# 安装服务
sudo cp airlinesim-autobuy.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable airlinesim-autobuy
sudo systemctl start airlinesim-autobuy
```

### 管理命令

```bash
# 查看服务状态
sudo systemctl status airlinesim-autobuy

# 查看日志
sudo journalctl -u airlinesim-autobuy -f

# 重启服务
sudo systemctl restart airlinesim-autobuy

# 停止服务
sudo systemctl stop airlinesim-autobuy
```

## Docker 部署

### Dockerfile

在项目根目录创建 `Dockerfile`：

```dockerfile
# 构建阶段
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache nodejs npm make

WORKDIR /app
COPY . .

# 构建前端
RUN make frontend-install && make frontend-build

# 编译 Go 二进制
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o autobuy ./cmd/autobuy

# 运行阶段
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

RUN adduser -D -g '' autobuy
USER autobuy

WORKDIR /app
COPY --from=builder /app/autobuy .
COPY --from=builder /app/configs ./configs

EXPOSE 9090

ENTRYPOINT ["./autobuy"]
```

### 构建镜像

```bash
docker build -t airlinesim-autobuy:latest .
```

### 运行容器

```bash
docker run -d \
  --name airlinesim-autobuy \
  -p 9090:9090 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/sessions:/app/sessions \
  --restart unless-stopped \
  airlinesim-autobuy:latest
```

## Docker Compose 部署

### docker-compose.yml

```yaml
version: '3.8'

services:
  autobuy:
    build: .
    container_name: airlinesim-autobuy
    restart: unless-stopped
    ports:
      - "9090:9090"
    volumes:
      - ./configs:/app/configs
      - ./sessions:/app/sessions
    environment:
      - TZ=Asia/Shanghai
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:9090/api/status"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### 启动服务

```bash
# 启动
docker compose up -d

# 查看日志
docker compose logs -f

# 停止
docker compose down

# 重新构建并启动
docker compose up -d --build
```

## 反向代理配置

### Nginx

在生产环境中，建议使用 Nginx 作为反向代理，提供 SSL 终结和访问控制：

```nginx
server {
    listen 443 ssl;
    server_name autobuy.example.com;

    ssl_certificate /etc/nginx/ssl/autobuy.crt;
    ssl_certificate_key /etc/nginx/ssl/autobuy.key;

    location / {
        proxy_pass http://127.0.0.1:9090;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # 可选：添加基本认证
    auth_basic "Restricted Access";
    auth_basic_user_file /etc/nginx/.htpasswd;
}
```

### 启用基本认证

```bash
# 安装 htpasswd 工具
sudo apt-get install apache2-utils

# 创建用户
sudo htpasswd -c /etc/nginx/.htpasswd admin

# 测试配置
sudo nginx -t

# 重载配置
sudo systemctl reload nginx
```

## 监控

### 健康检查

通过 API 端点检查服务状态：

```bash
curl http://localhost:9090/api/status
```

成功响应示例：

```json
{
  "running": true,
  "start_time": "2024-01-01T00:00:00Z",
  "servers": [
    {
      "host": "free1",
      "scan_count": 100,
      "found_count": 5,
      "bought_count": 2,
      "failed_count": 0,
      "last_scan": "2024-01-01T12:00:00Z"
    }
  ],
  "scan_count": 100,
  "found_count": 5,
  "bought_count": 2,
  "failed_count": 0
}
```

### Prometheus 集成

可以通过在 Nginx 或其他反向代理层添加 Prometheus 指标导出。目前程序本身未内置 metrics 端点，但可以通过 API 响应数据自行采集。

### 日志监控

systemd 日志可以通过 `journalctl` 查看：

```bash
# 实时追踪日志
journalctl -u airlinesim-autobuy -f

# 查看最近 100 条日志
journalctl -u airlinesim-autobuy -n 100

# 查看指定时间范围
journalctl -u airlinesim-autobuy --since "1 hour ago"
```

Docker 日志：

```bash
docker logs -f airlinesim-autobuy
```

## 数据持久化

需要持久化的数据：

| 路径 | 说明 |
|------|------|
| `configs/config.yaml` | 配置文件 |
| `session.json` | 会话文件（可重建） |

建议在 Docker 部署时使用 volume 挂载这些文件，确保数据不会因容器重建而丢失。

## 更新升级

```bash
# 拉取最新代码
git pull origin main

# 重新构建
make all

# 如果使用 systemd
sudo systemctl restart airlinesim-autobuy

# 如果使用 Docker Compose
docker compose up -d --build
```
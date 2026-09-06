# 배포

AirlineSim Autobuy는 단일 Go 바이너리로 컴파일되므로 배포가 매우 간단합니다. 이 가이드에서는 다양한 환경에서의 배포 방법을 설명합니다.

## systemd (Linux)

Linux 서버에서 systemd를 사용하여 AirlineSim Autobuy를 서비스로 등록할 수 있습니다.

### 1. 바이너리 설치

```bash
# 바이너리 다운로드 또는 빌드
sudo cp autobuy-linux-amd64 /usr/local/bin/airlinesim-autobuy
sudo chmod +x /usr/local/bin/airlinesim-autobuy
```

### 2. 설정 파일 준비

```bash
sudo mkdir -p /etc/airlinesim-autobuy
sudo cp configs/config.yaml /etc/airlinesim-autobuy/config.yaml
sudo chmod 600 /etc/airlinesim-autobuy/config.yaml
```

### 3. systemd 서비스 파일 생성

`/etc/systemd/system/airlinesim-autobuy.service`:

```ini
[Unit]
Description=AirlineSim Autobuy - Automated Aircraft Market Monitor
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=airlinesim
Group=airlinesim
WorkingDirectory=/etc/airlinesim-autobuy
ExecStart=/usr/local/bin/airlinesim-autobuy
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal
Environment=HOME=/etc/airlinesim-autobuy

# 보안 설정
NoNewPrivileges=true
ProtectSystem=full
ProtectHome=true
PrivateTmp=true
CapabilityBoundingSet=
AmbientCapabilities=

[Install]
WantedBy=multi-user.target
```

### 4. 서비스 활성화 및 실행

```bash
# 전용 사용자 생성
sudo useradd -r -s /bin/false -m -d /etc/airlinesim-autobuy airlinesim

# 설정 파일 권한 설정
sudo chown -R airlinesim:airlinesim /etc/airlinesim-autobuy

# 서비스 등록 및 실행
sudo systemctl daemon-reload
sudo systemctl enable airlinesim-autobuy
sudo systemctl start airlinesim-autobuy

# 상태 확인
sudo systemctl status airlinesim-autobuy
```

### 5. 로그 확인

```bash
# 실시간 로그
sudo journalctl -u airlinesim-autobuy -f

# 최근 로그
sudo journalctl -u airlinesim-autobuy -n 100
```

## Docker

### Dockerfile

프로젝트 루트에 다음 Dockerfile을 생성합니다:

```dockerfile
# 빌드 단계
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache nodejs npm

WORKDIR /app
COPY . .

# 프론트엔드 빌드
RUN cd internal/webui/frontend && npm install && npm run build

# Go 바이너리 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/autobuy ./cmd/autobuy

# 실행 단계
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/autobuy /usr/local/bin/autobuy

EXPOSE 9090

ENTRYPOINT ["autobuy"]
```

### Docker 이미지 빌드

```bash
docker build -t airlinesim-autobuy:latest .
```

### Docker 컨테이너 실행

```bash
docker run -d \
  --name airlinesim-autobuy \
  --restart unless-stopped \
  -p 9090:9090 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/session.json:/app/session.json \
  airlinesim-autobuy:latest
```

::: tip
설정 파일과 세션 파일을 호스트에 마운트하여 설정 변경 및 세션 지속성을 보장합니다.
:::

## Docker Compose

### docker-compose.yml

```yaml
version: "3.8"

services:
  autobuy:
    image: airlinesim-autobuy:latest
    container_name: airlinesim-autobuy
    restart: unless-stopped
    ports:
      - "9090:9090"
    volumes:
      - ./configs:/app/configs
      - ./session.json:/app/session.json
      - ./logs:/var/log/autobuy
    environment:
      - TZ=UTC
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:9090/api/status"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
```

### 실행

```bash
docker-compose up -d
```

### 로그 확인

```bash
docker-compose logs -f
```

## 리버스 프록시 설정

### Nginx

```nginx
server {
    listen 80;
    server_name autobuy.example.com;

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
}
```

### Caddy

```caddy
autobuy.example.com {
    reverse_proxy 127.0.0.1:9090
}
```

## 모니터링

### 헬스체크

Web UI가 `/api/status` 엔드포인트를 제공하므로, 헬스체크에 활용할 수 있습니다:

```bash
# 상태 확인
curl http://localhost:9090/api/status

# 응답 예시
{
  "running": true,
  "start_time": "2024-01-01T00:00:00Z",
  "servers": [
    {
      "host": "free1",
      "scan_count": 150,
      "found_count": 12,
      "bought_count": 3,
      "failed_count": 1
    }
  ],
  "scan_count": 150,
  "found_count": 12,
  "bought_count": 3,
  "failed_count": 1
}
```

### Prometheus / Grafana (준비 중)

Prometheus 메트릭 엔드포인트는 현재 개발 중입니다. 향후 버전에서 `/metrics` 엔드포인트가 추가될 예정입니다.

### 모니터링 스크립트 예시

```bash
#!/bin/bash
# autobuy-healthcheck.sh

STATUS=$(curl -sf http://localhost:9090/api/status | jq -r '.running')
if [ "$STATUS" != "true" ]; then
    echo "AirlineSim Autobuy engine is not running!"
    # 알림 전송 (예: 이메일, Slack 등)
    exit 1
fi

# 구매 통계 확인
BOUGHT=$(curl -sf http://localhost:9090/api/status | jq -r '.bought_count')
echo "Total purchases made: $BOUGHT"
```

## 성능 고려 사항

### 메모리 사용량

AirlineSim Autobuy는 매우 가벼운 애플리케이션입니다:

| 구성 요소 | 예상 사용량 |
|----------|----------|
| Go 바이너리 | ~15-25 MB |
| 추가 메모리 | ~10-50 MB (처리 중인 데이터) |
| 총 사용량 | ~25-75 MB |

### 네트워크 사용량

- 각 서버 스캔당 약 50-200 KB 다운로드
- 30초 간격, 1개 서버 기준: 하루 약 5-15 MB

## 보안 권장 사항

1. **방화벽 설정**: Web UI 포트(9090)는 필요한 경우에만 외부에 노출하세요.
2. **리버스 프록시 사용**: Nginx 또는 Caddy를 앞에 두고 TLS를 적용하세요.
3. **파일 권한**: 설정 파일의 권한을 600으로 설정하세요.
4. **전용 사용자**: systemd 서비스는 전용 시스템 사용자로 실행하세요.
5. **Docker 보안**: 루트가 아닌 사용자로 컨테이너를 실행하는 것을 고려하세요.

## 관련 문서

- [빠른 시작](/ko/guide/quickstart) - 설치 및 기본 실행
- [설정 가이드](/ko/guide/configuration) - 모든 설정 옵션
- [FAQ](/ko/guide/faq) - 문제 해결
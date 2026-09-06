<h1 align="center">✈️ AirlineSim Autobuy</h1>
<p align="center">
  <em>AirlineSim 중고 항공기 시장 자동 모니터링 및 구매 도구</em>
</p>

<p align="center">
  <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go" alt="Go"></a>
  <a href="https://vuejs.org"><img src="https://img.shields.io/badge/Vue-3.4-4FC08D?logo=vue.js" alt="Vue"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-AGPL--3.0-blue" alt="License"></a>
</p>

<p align="center">
  <a href="#기능">기능</a> •
  <a href="#빠른-시작">빠른 시작</a> •
  <a href="#설정">설정</a> •
  <a href="#배포">배포</a> •
  <a href="docs/ko/guide/quickstart.md">전체 문서</a>
</p>

---

## 기능

- **🔍 실시간 모니터링** — 설정된 서버의 중고 항공기 시장을 지속적으로 스캔하여 새 매물을 즉시 감지
- **📋 스마트 규칙 엔진** — 기종, 유형, 가격, 연령, 상태, 사이클, 임대료로 구매 규칙 정의
- **🤖 자동 구매** — 즉시 구매 및 경매 입찰 지원, 스내치 및 잔액 보호 기능
- **🌐 다중 서버 및 계정** — 여러 게임 서버를 동시 모니터링, 각 서버에 독립 인증 계정 사용
- **⚙️ 웹 관리 UI** — 엔진 제어, 규칙 관리, 실시간 로그, 설정을 위한 대시보드
- **🌍 다국어 지원** — English, 中文, 한국어
- **🚀 단일 바이너리** — Go 컴파일 단일 바이너리, Vue 프론트엔드 내장, 제로 의존성 배포

## 빠른 시작

### 전제 조건

- Go 1.22+
- Node.js 18+ (프론트엔드 개발 시, 선택 사항)
- AirlineSim 게임 계정

### 설치

```bash
# 저장소 복제
git clone https://github.com/Maicarons/airlinesim-autobuy.git
cd airlinesim-autobuy

# 프론트엔드 의존성 설치 및 빌드
make frontend-install
make frontend-build

# Go 백엔드 빌드
make build

# 설정 파일 복사 및 편집
cp configs/config.example.yaml configs/config.yaml
# configs/config.yaml 파일에 계정 정보 입력
```

### 실행

```bash
./autobuy
```

브라우저에서 http://localhost:9090 을 열어 웹 UI에 접속하세요.

## 설정

`configs/config.yaml` 파일 편집:

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
auths:
  - username: "your_email@example.com"
    password: "your_password"
monitor:
  interval: 30          # 폴링 간격 (초)
  min_balance: 1000000  # 구매 후 유지할 최소 잔액
```

### 구매 규칙

```yaml
rules:
  - name: "A320-200 Heavy 모니터링"
    enabled: true
    server_id: 0
    auth_id: 0
    match:
      family_id: "1200300"       # A320 / A321 계열
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

## 프로젝트 구조

```
airlinesim-autobuy/
├── cmd/autobuy/          # 진입점
├── internal/
│   ├── auth/             # 인증 및 세션 관리
│   ├── client/           # 속도 제한 HTTP 클라이언트
│   ├── collector/        # 시장 페이지 수집
│   ├── config/           # 설정 관리
│   ├── engine/           # 파이프라인 오케스트레이터
│   ├── executor/         # 구매 실행
│   ├── marketdata/       # 항공기 참조 데이터
│   ├── notifier/         # 알림 시스템
│   ├── parser/           # HTML 파서
│   ├── rules/            # 규칙 엔진
│   └── webui/            # 웹 관리 UI
├── configs/              # 설정 파일
├── docs/                 # VitePress 문서
├── Makefile
└── README.md
```

## 문서

전체 문서는 [docs](docs/) 디렉토리에 있습니다:

| 언어 | 사용자 가이드 | 개발자 문서 |
|------|-------------|------------|
| English | [Quick Start](docs/en/guide/quickstart.md) | [Architecture](docs/en/reference/architecture.md) |
| 中文 | [快速开始](docs/zh-CN/guide/quickstart.md) | [架构说明](docs/zh-CN/reference/architecture.md) |
| 한국어 | [빠른 시작](docs/ko/guide/quickstart.md) | [아키텍처](docs/ko/reference/architecture.md) |

## 배포

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

## 라이선스

[AGPL-3.0](LICENSE)

## 면책 조항

이 도구는 개인 학습 용도로만 사용하세요. 자동화된 게임 플레이는 게임 이용약관을 위반할 수 있습니다. 사용에 따른 책임은 본인에게 있습니다. 합리적인 폴링 간격을 설정하고 게임 서버에 부하를 주지 않도록 주의하세요.
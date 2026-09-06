# 빠른 시작

AirlineSim Autobuy는 AirlineSim 게임의 중고 항공기 시장을 실시간으로 모니터링하고 사용자 정의 규칙에 따라 자동으로 항공기를 구매하는 Go 기반 도구입니다.

## 전제 조건

### 시스템 요구 사항

| 항목 | 요구 사항 |
|------|----------|
| 운영 체제 | Windows, Linux, macOS |
| Go (개발용) | Go 1.25.5 이상 |
| Node.js (프론트엔드 빌드용) | Node.js 18 이상, npm 9 이상 |
| AirlineSim 계정 | 유효한 게임 계정 (각 서버별) |

### 선택 도구

- **Air** (개발 중 핫 리로드): `go install github.com/air-verse/air@latest`
- **golangci-lint** (린터): `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

## 설치

### 1. 저장소 클론

```bash
git clone https://github.com/Maicarons/airlinesim-autobuy.git
cd airlinesim-autobuy
```

### 2. 프론트엔드 의존성 설치

```bash
make frontend-install
```

또는 직접:

```bash
cd internal/webui/frontend && npm install
```

### 3. 설정 파일 준비

`configs/config.yaml` 파일이 없으면 애플리케이션 실행 시 자동으로 기본 설정 파일이 생성됩니다. 기본 설정은 다음과 같습니다:

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
auth:
  username: ""
  password: ""
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
  - name: 예시 규칙 - 경제형 협동체
    enabled: false
    priority: 10
    match:
      types:
        - B737-800
        - A320-200
      price_range:
        min: 500000
        max: 5000000
      max_age: 15
      max_cycles: 30000
      condition_min: 70
      offer_types:
        - auction
        - immediate
      financing:
        - cash
        - credit
        - lease
    action:
      auto_buy: true
      snatch: false
      max_bid_increment: 100000
```

## 빌드

### 전체 빌드 (프론트엔드 + 백엔드)

```bash
make all
```

이 명령어는 다음을 순서대로 실행합니다:
1. 프론트엔드 npm 의존성 설치
2. Vue 3 프론트엔드 빌드
3. Go 바이너리 컴파일 (프론트엔드 내장)

### 개별 빌드

```bash
# 프론트엔드만 빌드
make frontend-build

# 백엔드만 빌드
make build

# 클린 빌드
make clean
```

### 크로스 플랫폼 빌드

```bash
make build-all
```

지원 대상:
- `autobuy-windows-amd64.exe` (Windows)
- `autobuy-linux-amd64` (Linux)
- `autobuy-darwin-amd64` (macOS)

## 설정

### AirlineSim 계정 설정

`configs/config.yaml` 파일에서 인증 정보를 설정합니다:

```yaml
auth:
  username: your-email@example.com
  password: your-password
  session_file: session.json
```

::: warning
자격 증명은 일반 텍스트로 저장되므로 파일 권한을 적절히 설정하세요.
:::

### 구매 규칙 설정

기본 규칙을 수정하거나 새 규칙을 추가하여 자동 구매 조건을 정의하세요. 자세한 내용은 [구매 규칙](/ko/guide/rules) 문서를 참조하세요.

## 실행

### 기본 실행

```bash
make run
```

또는 직접:

```bash
go run ./cmd/autobuy
```

### 개발 모드 (핫 리로드)

```bash
make dev
```

### 실행 옵션

애플리케이션은 기본적으로 `configs/config.yaml` 파일을 읽습니다. 명령줄 인수는 현재 지원되지 않으며, 모든 설정은 YAML 파일을 통해 관리됩니다.

## 첫 단계

### 1. 계정 설정 확인

`configs/config.yaml` 파일에 올바른 AirlineSim 계정 정보를 입력했는지 확인합니다.

### 2. Web UI 접속

브라우저에서 `http://localhost:9090`으로 접속합니다.

### 3. 구매 규칙 생성

Web UI의 **Rules** 페이지에서 새 규칙을 생성합니다:
- **Rule Name**: 규칙의 이름
- **Server**: 모니터링할 게임 서버 선택
- **Account**: 사용할 계정 선택
- **Aircraft Type**: 대상 항공기 유형
- **Price Range**: 가격 범위 설정
- **Max Age / Max Cycles**: 최대 기령 및 사이클
- **Min Condition**: 최소 상태
- **Auto Buy**: 자동 구매 활성화

### 4. 엔진 시작

Dashboard에서 **Start Engine** 버튼을 클릭하여 모니터링을 시작합니다.

### 5. 로그 확인

Logs 페이지에서 실시간 엔진 활동 로그를 확인할 수 있습니다.

## 기본 워크플로

```
설정 파일 작성 → 계정 인증 → 서버 연결 → 시장 스캔 시작
    → 항공기 매물 발견 → 규칙 평가 → 일치 항목 발견
    → 자동 구매 실행 → 알림 전송
```

## 다음 단계

- [설정 가이드](/ko/guide/configuration) - 모든 설정 옵션 상세 설명
- [구매 규칙](/ko/guide/rules) - 규칙 생성 및 최적화
- [웹 UI](/ko/guide/webui) - 대시보드 사용법
- [배포 가이드](/ko/guide/deployment) - 프로덕션 환경 배포
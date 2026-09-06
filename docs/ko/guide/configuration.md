# 설정

AirlineSim Autobuy는 YAML 설정 파일을 통해 모든 동작을 제어합니다. 기본 설정 파일은 `configs/config.yaml`에 위치하며, 애플리케이션 실행 시 자동으로 생성됩니다.

## 설정 파일 구조

설정 파일은 다음과 같은 최상위 섹션으로 구성됩니다:

```yaml
servers:     # 게임 서버 연결 설정 (다중 서버 지원)
auths:       # 인증 계정 설정 (다중 계정 지원)
monitor:     # 시장 모니터링 동작 설정
notifier:    # 알림 채널 설정
webui:       # 웹 관리 인터페이스 설정
rules:       # 구매 규칙 목록
```

## 서버 (Servers)

AirlineSim Autobuy는 여러 게임 서버를 동시에 모니터링할 수 있습니다.

### 단일 서버 설정

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
```

### 다중 서버 설정

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
  - host: free2
    base_url: https://free2.airlinesim.aero
  - host: free3
    base_url: https://free3.airlinesim.aero
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `host` | string | 서버 식별자 (표시용) |
| `base_url` | string | 게임 서버 기본 URL |

기본값: `host: "free1"`, `base_url: "https://free1.airlinesim.aero"`

## 인증 (Auth / Auths)

### 단일 계정 설정 (레거시)

```yaml
auth:
  username: your-email@example.com
  password: your-password
  session_file: session.json
```

### 다중 계정 설정 (권장)

```yaml
auths:
  - username: account1@example.com
    password: password1
    session_file: session1.json
  - username: account2@example.com
    password: password2
    session_file: session2.json
```

| 필드 | 타입 | 설명 |
|------|------|------|
| `username` | string | AirlineSim 계정 이메일 |
| `password` | string | 계정 비밀번호 |
| `session_file` | string | 세션 쿠키를 저장할 파일 경로 |

::: warning 보안 주의
- 비밀번호는 일반 텍스트로 저장됩니다. 파일 권한을 600으로 설정하는 것을 권장합니다.
- `.gitignore`에 설정 파일을 추가하여 자격 증명이 저장소에 커밋되지 않도록 하세요.
- `session.json` 파일에도 민감한 세션 정보가 포함될 수 있습니다.
:::

### 세션 관리

- 세션은 `sar.simulogics.games` API를 통해 인증됩니다.
- 인증 성공 시 `as-sid` 쿠키가 설정됩니다.
- 세션은 `session.json` 파일에 저장되어 재시작 시 복원됩니다.
- 세션 만료 시 자동으로 재인증이 수행됩니다.
- 세션 유효 기간은 최대 24시간입니다.

## 모니터링 설정 (Monitor)

```yaml
monitor:
  interval: 30
  jitter: 10
  request_timeout: 30
  min_balance: 1000000
```

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `interval` | int | 30 | 서버 스캔 간격 (초) |
| `jitter` | int | 10 | 무작위 지연 시간 (초) |
| `request_timeout` | int | 30 | HTTP 요청 타임아웃 (초) |
| `min_balance` | float | 1000000 | 구매 후 유지해야 할 최소 계정 잔액 (AS$) |

### Interval

시장 페이지를 스캔하는 주기입니다. 값이 너무 작으면 게임 서버에 과도한 부하를 줄 수 있습니다. 권장값: 15~60초.

### Jitter

각 요청 전에 추가되는 무작위 지연 시간입니다. 일정한 패턴을 방지하여 탐지를 회피하는 데 도움이 됩니다. 실제 간격은 `interval + random(0, jitter)` 범위입니다.

### Min Balance

자동 구매 실행 시 계정에 유지해야 할 최소 잔액입니다. 이 값보다 큰 금액의 구매는 경고를 출력하지만 현재는 차단하지 않습니다. 향후 버전에서 하드 제한으로 구현될 예정입니다.

## 알림 (Notifier)

```yaml
notifier:
  console: true
  discord_webhook: https://discord.com/api/webhooks/...
```

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `console` | bool | true | 콘솔 출력 활성화 |
| `discord_webhook` | string | "" | Discord 웹훅 URL (비활성화 시 빈 값) |

### 콘솔 알림

콘솔 알림은 구조화된 로그를 사용하여 이벤트를 출력합니다:

| 이벤트 | 아이콘 | 로그 레벨 |
|--------|--------|----------|
| 구매 성공 | 🛒 | Info |
| 구매 실패 | ❌ | Warn |
| 오류 발생 | ⚠️ | Error |
| 항공기 발견 | 🔍 | Info |
| 입찰 완료 | 💰 | Info |

### Discord 알림

Discord 웹훅을 설정하면 동일한 이벤트가 Discord 채널로 전송됩니다. 웹훅 URL은 Discord 서버 설정에서 생성할 수 있습니다.

## Web UI

```yaml
webui:
  enabled: true
  host: 0.0.0.0
  port: 9090
```

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `enabled` | bool | true | 웹 UI 활성화 |
| `host` | string | "0.0.0.0" | 바인딩할 호스트 주소 |
| `port` | int | 9090 | 리스닝 포트 |

### 보안 참고 사항

- 프로덕션 환경에서는 `host`를 `127.0.0.1`로 설정하고 리버스 프록시(Nginx, Caddy)를 사용하는 것을 권장합니다.
- Web UI는 현재 인증을 제공하지 않습니다. 필요한 경우 리버스 프록시에서 인증을 추가하세요.
- Web UI를 비활성화하려면 `enabled: false`로 설정합니다.

## 전체 설정 예시

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
  - host: free2
    base_url: https://free2.airlinesim.aero

auths:
  - username: main-account@example.com
    password: secure-password-1
    session_file: session_main.json
  - username: alt-account@example.com
    password: secure-password-2
    session_file: session_alt.json

monitor:
  interval: 45
  jitter: 15
  request_timeout: 30
  min_balance: 2000000

notifier:
  console: true
  discord_webhook: https://discord.com/api/webhooks/your-webhook-id/your-webhook-token

webui:
  enabled: true
  host: 127.0.0.1
  port: 9090

rules:
  - name: A320-200 Heavy Lease Monitor
    enabled: true
    priority: 10
    server_id: 0
    auth_id: 0
    match:
      family_id: "1200300"
      type_id: ""
      price_range:
        min: 0
        max: 5000000
      max_age: 20
      max_cycles: 50000
      condition_min: 50
      offer_types:
        - auction
        - immediate
      financing:
        - cash
        - credit
        - lease
      sort_by: price_asc
    action:
      auto_buy: false
      snatch: false
      max_bid_increment: 100000
```

## 설정 파일 핫 리로드

설정 파일은 변경 사항을 실시간으로 감지하여 자동으로 리로드합니다. 파일이 수정되면:

1. 변경 사항이 감지됩니다 (5초마다 폴링).
2. 설정이 자동으로 리로드됩니다.
3. 규칙이 엔진에 다시 로드됩니다.
4. 로그에 "configuration hot-reloaded" 메시지가 출력됩니다.

::: tip
Web UI를 통해 설정을 변경하면 자동으로 파일에 저장되고 리로드됩니다. Web UI를 사용하지 않는 경우 텍스트 편집기로 파일을 직접 수정하면 됩니다.
:::

## 레거시 호환성

이전 버전의 단일 서버/단일 계정 설정(`server`, `auth`)은 자동으로 다중 서버/다중 계정 형식으로 마이그레이션됩니다:

```yaml
# 레거시 형식 (자동 변환됨)
server:
  host: free1
  base_url: https://free1.airlinesim.aero
auth:
  username: user@example.com
  password: pass
  session_file: session.json

# 변환 후
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
auths:
  - username: user@example.com
    password: pass
    session_file: session.json
```

## 관련 문서

- [설정 참조](/ko/reference/config) - 모든 설정 필드 상세 참조
- [구매 규칙](/ko/guide/rules) - 규칙 생성 및 설정
- [배포 가이드](/ko/guide/deployment) - 프로덕션 배포
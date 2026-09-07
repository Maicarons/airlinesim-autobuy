# 설정 참조

이 문서는 AirlineSim Autobuy 설정 파일의 모든 필드를 상세히 설명합니다. 설정은 `configs/config.yaml`에 위치한 YAML 파일입니다.

## Config (최상위)

```yaml
servers:  []ServerConfig   # 게임 서버 목록
auths:    []AuthConfig     # 인증 계정 목록
auth:     AuthConfig       # (사용 중단) 단일 계정 설정
monitor:  MonitorConfig    # 모니터링 동작 설정
notifier: NotifierConfig   # 알림 채널 설정
webui:    WebUIConfig      # 웹 UI 설정
rules:    []RuleConfig     # 구매 규칙 목록
server:   ServerConfig     # (사용 중단) 단일 서버 설정
```

## ServerConfig

게임 서버 연결 정보를 정의합니다.

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
```

| 필드 | 타입 | 필수 | 기본값 | 설명 |
|------|------|------|--------|------|
| `host` | string | 예 | `"free1"` | 서버 식별자 (표시용) |
| `base_url` | string | 예 | `"https://free1.airlinesim.aero"` | 게임 서버 기본 URL |

**참고:** 서버가 하나도 설정되지 않은 경우, 기본값으로 `free1` 서버가 자동 생성됩니다.

## AuthConfig

인증 계정 정보를 정의합니다.

```yaml
auths:
  - username: account@example.com
    password: secure-password
    session_file: session.json
```

| 필드 | 타입 | 필수 | 기본값 | 설명 |
|------|------|------|--------|------|
| `username` | string | 예 | `""` | AirlineSim 계정 이메일 |
| `password` | string | 예 | `""` | 계정 비밀번호 |
| `session_file` | string | 아니오 | `"session.json"` | 세션 쿠키 저장 파일 경로 |

**참고:** `auth` (단수) 필드는 레거시 호환성용입니다. `auths` (복수) 사용을 권장합니다.

## MonitorConfig

시장 모니터링 동작을 정의합니다.

```yaml
monitor:
  interval: 30
  jitter: 10
  request_timeout: 30
  min_balance: 1000000
```

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `interval` | int | `30` | 서버 스캔 간격 (초) |
| `jitter` | int | `10` | 무작위 지연 시간 (초) |
| `request_timeout` | int | `30` | HTTP 요청 타임아웃 (초) |
| `min_balance` | float | `1000000` | 구매 후 유지할 최소 계정 잔액 (AS$) |

### 필드 상세

**`interval`**
- 각 서버의 시장 페이지를 스캔하는 주기입니다.
- 권장 범위: 15~60초
- 너무 낮은 값은 게임 서버에 과도한 부하를 줄 수 있습니다.

**`jitter`**
- 각 요청 전에 추가되는 무작위 지연 시간입니다.
- 실제 간격 = `interval + random(0, jitter)` 초
- 탐지 회피 및 트래픽 패턴 분산에 도움이 됩니다.

**`request_timeout`**
- 단일 HTTP 요청의 최대 대기 시간입니다.
- 네트워크 지연이나 서버 응답 없음 상황에서 빠른 실패를 위해 사용됩니다.

**`min_balance`**
- 자동 구매 실행 시 계정에 유지해야 할 최소 잔액입니다.
- 현재는 경고만 출력하며, 향후 하드 제한으로 구현될 예정입니다.

## NotifierConfig

알림 채널을 정의합니다.

```yaml
notifier:
  console: true
  discord_webhook: https://discord.com/api/webhooks/...
  # dingtalk_webhook: https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN
  # dingtalk_secret: "YOUR_SIGNING_SECRET"
```

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `console` | bool | `true` | 콘솔(표준 출력) 알림 활성화 |
| `discord_webhook` | string | `""` | Discord 웹훅 URL (빈 값: 비활성화) |
| `dingtalk_webhook` | string | `""` | DingTalk 커스텀 봇 웹훅 URL (빈 값: 비활성화) |
| `dingtalk_secret` | string | `""` | DingTalk 봇의 HMAC-SHA256 서명 비밀키 (봇이 서명 확인을 필요로 하는 경우 설정) |

## WebUIConfig

웹 관리 인터페이스 설정을 정의합니다.

```yaml
webui:
  enabled: true
  host: 0.0.0.0
  port: 9090
```

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `enabled` | bool | `true` | 웹 UI 활성화 |
| `host` | string | `"0.0.0.0"` | HTTP 서버 바인딩 주소 |
| `port` | int | `9090` | HTTP 서버 리스닝 포트 |

**보안 참고:**
- `host: "0.0.0.0"`은 모든 네트워크 인터페이스에서 접근 가능합니다.
- 프로덕션에서는 `host: "127.0.0.1"`로 설정하고 리버스 프록시를 사용하는 것을 권장합니다.
- Web UI는 현재 인증을 제공하지 않습니다.

## RuleConfig

구매 규칙을 정의합니다.

```yaml
rules:
  - name: 규칙 이름
    enabled: true
    priority: 10
    server_id: 0
    auth_id: 0
    match:
      # MatchConfig (아래 참조)
    action:
      # ActionConfig (아래 참조)
```

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `name` | string | `""` | 규칙 이름 (필수) |
| `enabled` | bool | `true` | 규칙 활성화 여부 |
| `priority` | int | `10` | 규칙 우선순위 (높을수록 우선) |
| `server_id` | int | `0` | 적용 서버 인덱스 (-1: 모든 서버) |
| `auth_id` | int | `0` | 사용 계정 인덱스 (-1: 기본 계정) |

**`server_id`**
- `servers` 배열의 0-based 인덱스입니다.
- `-1`로 설정하면 모든 서버에 적용됩니다.
- 배열 범위를 벗어나면 첫 번째 서버(인덱스 0)가 사용됩니다.

**`auth_id`**
- `auths` 배열의 0-based 인덱스입니다.
- `-1`로 설정하면 기본 계정이 사용됩니다.
- 배열 범위를 벗어나면 첫 번째 계정(인덱스 0)이 사용됩니다.

## MatchConfig

규칙의 일치 조건을 정의합니다.

```yaml
match:
  family_id: "1200300"
  type_id: "16"
  types:
    - "Airbus A320-200 heavy"
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
  sort_by: price_asc
```

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `family_id` | string | `""` | Wicket 항공기 패밀리 ID |
| `type_id` | string | `""` | Wicket 항공기 유형 ID |
| `types` | []string | `[]` | 항공기 유형명 목록 (대소문자 구분 없음) |
| `price_range.min` | float | `0` | 최소 가격 (0: 제한 없음) |
| `price_range.max` | float | `0` | 최대 가격 (0: 제한 없음) |
| `max_age` | int | `0` | 최대 기령 (년, 0: 제한 없음) |
| `max_cycles` | int | `0` | 최대 비행 사이클 (0: 제한 없음) |
| `condition_min` | float | `0` | 최소 상태 (0-100, 0: 제한 없음) |
| `offer_types` | []string | `[]` | 허용 제공 유형 (빈 배열: 모두 허용) |
| `financing` | []string | `[]` | 허용 금융 옵션 (빈 배열: 모두 허용) |
| `sort_by` | string | `""` | 시장 정렬 기준 |

### `offer_types` 허용 값

| 값 | 설명 |
|-----|------|
| `"auction"` | 경매 매물 |
| `"immediate"` | 즉시 구매 매물 |

### `financing` 허용 값

| 값 | 설명 |
|-----|------|
| `"cash"` | 현금 구매 |
| `"credit"` | 신용 할부 |
| `"lease"` | 리스 |

### `sort_by` 허용 값

| 값 | 설명 |
|-----|------|
| `""` (빈 문자열) | 기본 정렬 |
| `"price_asc"` | 가격 오름차순 |
| `"price_desc"` | 가격 내림차순 |
| `"age_asc"` | 기령 오름차순 |
| `"age_desc"` | 기령 내림차순 |

## ActionConfig

규칙의 실행 작업을 정의합니다.

```yaml
action:
  auto_buy: true
  snatch: false
  max_bid_increment: 100000
```

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `auto_buy` | bool | `false` | 자동 구매 활성화 |
| `snatch` | bool | `false` | 스내치 모드 활성화 (예비 필드) |
| `max_bid_increment` | float | `0` | 최대 입찰 증가액 (AS$, 0: 제한 없음) |

## PriceRange

가격 범위를 정의합니다.

```yaml
price_range:
  min: 500000
  max: 5000000
```

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `min` | float | `0` | 최소 가격 (0: 제한 없음) |
| `max` | float | `0` | 최대 가격 (0: 제한 없음) |

## 기본 설정

설정 파일이 없을 때 자동으로 생성되는 기본 설정입니다:

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
  # discord_webhook: ""
  # dingtalk_webhook: ""
  # dingtalk_secret: ""
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

## 관련 문서

- [설정 가이드](/ko/guide/configuration) - 설정 설명
- [구매 규칙](/ko/guide/rules) - 규칙 설정 가이드
- [API 참조](/ko/reference/api) - REST API 엔드포인트
# 구매 규칙

AirlineSim Autobuy의 핵심 기능은 사용자 정의 규칙에 따라 항공기를 자동으로 구매하는 것입니다. 각 규칙은 일치 조건과 실행 작업을 정의합니다.

## 규칙 구조

규칙은 다음과 같은 구조를 가집니다:

```yaml
rules:
  - name: 규칙 이름
    enabled: true
    priority: 10
    server_id: 0
    auth_id: 0
    match:
      # 일치 조건
      family_id: ""
      type_id: ""
      types: []
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
      # 실행 작업
      auto_buy: false
      snatch: false
      max_bid_increment: 100000
```

### 공통 필드

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `name` | string | 필수 | 규칙을 식별하는 고유한 이름 |
| `enabled` | bool | true | 규칙 활성화 여부 |
| `priority` | int | 10 | 규칙 우선순위 (높을수록 우선) |
| `server_id` | int | 0 | 적용할 서버 인덱스 (-1: 모든 서버) |
| `auth_id` | int | 0 | 사용할 계정 인덱스 (-1: 기본 계정) |

## 일치 조건 (Match)

### 항공기 유형 조건

```yaml
match:
  family_id: "1200300"        # A320 / A321 패밀리 ID
  type_id: ""                  # 특정 유형 ID (빈 값: 모든 유형)
  types:                       # 유형 이름 목록
    - "Airbus A320-200 heavy"
    - "Airbus A320-200 medium"
```

**비교 방식:**
- `family_id`와 `type_id`는 Wicket 필터 파라미터로 사용되어 서버 측에서 필터링됩니다.
- `types`는 클라이언트 측에서 대소문자 구분 없이 비교됩니다.
- `family_id` 또는 `type_id`가 설정되면 서버 측 필터가 적용되어 더 효율적입니다.

### 가격 조건

```yaml
match:
  price_range:
    min: 500000     # 최소 가격 (AS$)
    max: 5000000    # 최대 가격 (AS$)
```

| 필드 | 타입 | 기본값 | 설명 |
|------|------|--------|------|
| `min` | float | 0 | 최소 가격 (0: 제한 없음) |
| `max` | float | 0 | 최대 가격 (0: 제한 없음) |

### 기령 조건

```yaml
match:
  max_age: 15  # 최대 기령 (년)
```

- `0`으로 설정하면 기령 제한이 없습니다.
- 게임 내 항공기 기령은 제조 연도를 기준으로 계산됩니다.

### 사이클 조건

```yaml
match:
  max_cycles: 30000  # 최대 비행 사이클
```

- `0`으로 설정하면 사이클 제한이 없습니다.
- 사이클은 항공기의 엔진 및 기체 수명을 나타냅니다.

### 상태 조건

```yaml
match:
  condition_min: 70  # 최소 상태 (%) 0-100
```

- `0`으로 설정하면 상태 제한이 없습니다.
- 상태는 항공기의 전반적인 컨디션을 백분율로 나타냅니다.

### 제공 유형

```yaml
match:
  offer_types:
    - auction     # 경매
    - immediate   # 즉시 구매
```

선택 가능한 값:
- `auction` - 경매 매물
- `immediate` - 즉시 구매 가능 매물

빈 배열로 설정하면 모든 유형을 허용합니다.

### 금융 옵션

```yaml
match:
  financing:
    - cash     # 현금
    - credit   # 할부
    - lease    # 리스
```

선택 가능한 값:
- `cash` - 현금 구매
- `credit` - 신용 할부
- `lease` - 리스

빈 배열로 설정하면 모든 금융 옵션을 허용합니다.

### 정렬 기준

```yaml
match:
  sort_by: price_asc  # 가격 오름차순
```

선택 가능한 값:
- `price_asc` - 가격 오름차순
- `price_desc` - 가격 내림차순
- `age_asc` - 기령 오름차순
- `age_desc` - 기령 내림차순
- 빈 값 - 기본 정렬

## 실행 작업 (Action)

### 자동 구매

```yaml
action:
  auto_buy: true   # 자동 구매 활성화
```

- `true`: 규칙에 일치하는 항공기를 자동으로 구매합니다.
- `false`: 항공기를 발견만 하고 구매하지 않습니다 (모니터링 전용).

### 스내치 (Snatch)

```yaml
action:
  snatch: true    # 스내치 활성화
```

스내치 기능은 경매에서 경쟁자를 제치고 빠르게 입찰하는 모드를 의미합니다. 현재는 예비 필드로 구현되어 있으며 향후 확장될 예정입니다.

### 최대 입찰 인크리먼트

```yaml
action:
  max_bid_increment: 100000  # 최대 입찰 증가액 (AS$)
```

- 경매에서 현재 가격에 추가로 지불할 수 있는 최대 금액입니다.
- 즉시 구매에는 영향을 미치지 않습니다.
- `0`으로 설정하면 입찰 인크리먼트 제한이 없습니다.

## 서버 선택

각 규칙은 특정 서버 또는 모든 서버에 적용할 수 있습니다.

```yaml
# 특정 서버에 적용 (인덱스 0: 첫 번째 서버)
rules:
  - name: Free1 전용 규칙
    server_id: 0
```

```yaml
# 모든 서버에 적용
rules:
  - name: 모든 서버 공통 규칙
    server_id: -1
```

`server_id`는 `servers` 배열의 인덱스입니다:
- `0`: 첫 번째 서버
- `1`: 두 번째 서버
- `-1`: 모든 서버 (기본값)

## 계정 선택

각 규칙은 특정 인증 계정을 사용할 수 있습니다.

```yaml
rules:
  - name: 첫 번째 계정 사용
    auth_id: 0
```

```yaml
rules:
  - name: 두 번째 계정 사용
    auth_id: 1
```

`auth_id`는 `auths` 배열의 인덱스입니다:
- `0`: 첫 번째 계정
- `1`: 두 번째 계정
- `-1`: 기본 계정 (기본값)

## 점수 시스템

규칙 엔진은 각 항공기-규칙 쌍에 대해 점수를 계산합니다. 점수가 높을수록 더 좋은 거래로 간주됩니다.

### 점수 구성 요소

| 구성 요소 | 최대 점수 | 설명 |
|-----------|----------|------|
| 가격 점수 | 50 | 최대 가격 대비 낮을수록 높은 점수 |
| 상태 점수 | 20 | 최소 상태 대비 높을수록 높은 점수 |
| 기령 점수 | 20 | 최대 기령 대비 낮을수록 높은 점수 |
| 즉시 구매 보너스 | 10 | 즉시 구매 가능 매물에 추가 점수 |
| 우선순위 보너스 | - | 규칙의 `priority` 값이 점수에 직접 추가 |

### 점수 계산 예시

```
가격: AS$ 3,000,000 / 최대 AS$ 5,000,000 → 가격 점수: 20
상태: 85% / 최소 70% → 상태 점수: 10
기령: 5년 / 최대 15년 → 기령 점수: 13.3
즉시 구매 → 보너스: 10
우선순위: 10 → 보너스: 10
총점: 63.3
```

## 규칙 예시

### 경제형 협동체 모니터링

```yaml
- name: 저가 A320 모니터링
  enabled: true
  priority: 10
  match:
    types:
      - "Airbus A320-200 heavy"
      - "Airbus A320-200 medium"
    price_range:
      min: 500000
      max: 3000000
    max_age: 15
    max_cycles: 30000
    condition_min: 70
    offer_types:
      - auction
      - immediate
    financing:
      - cash
      - credit
  action:
    auto_buy: true
    max_bid_increment: 100000
```

### 광폭동체 장거리용

```yaml
- name: B787 우선 매수
  enabled: true
  priority: 20
  match:
    family_id: "2200700"
    price_range:
      min: 1000000
      max: 15000000
    max_age: 10
    max_cycles: 15000
    condition_min: 80
    offer_types:
      - immediate
    financing:
      - cash
  action:
    auto_buy: true
    max_bid_increment: 200000
```

### 리스 전용 모니터링

```yaml
- name: ATR 리스 딜 찾기
  enabled: true
  priority: 5
  match:
    family_id: "1600200"
    max_age: 25
    condition_min: 40
    financing:
      - lease
  action:
    auto_buy: false
```

## 규칙 평가 흐름

```
항공기 매물 도착
  ↓
규칙 목록 반복 (우선순위 정렬)
  ↓
규칙 활성화 확인
  ↓
유형 일치 확인
  ↓
가격 범위 확인
  ↓
기령 확인
  ↓
사이클 확인
  ↓
상태 확인
  ↓
제공 유형 확인
  ↓
금융 옵션 확인
  ↓
점수 계산
  ↓
자동 구매 실행 (auto_buy=true인 경우)
```

## 관련 문서

- [설정 가이드](/ko/guide/configuration) - 전체 설정 옵션
- [설정 참조](/ko/reference/config) - 모든 설정 필드 상세 참조
- [API 참조](/ko/reference/api) - REST API 엔드포인트
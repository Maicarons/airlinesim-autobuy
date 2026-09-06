# API 참조

AirlineSim Autobuy Web UI는 RESTful API를 제공합니다. 모든 API 엔드포인트는 `/api` 접두사 아래에 있습니다.

## 기본 정보

- **기본 URL**: `http://localhost:9090/api`
- **콘텐츠 타입**: `application/json`
- **인증**: 현재 Web UI는 별도의 인증을 제공하지 않습니다.

## 상태 및 제어

### 엔진 상태 조회

```
GET /api/status
```

엔진의 현재 상태와 통계를 반환합니다.

**응답 예시:**

```json
{
  "running": true,
  "start_time": "2024-01-01T00:00:00Z",
  "servers": [
    {
      "host": "free1",
      "scan_count": 150,
      "found_count": 12,
      "bought_count": 3,
      "failed_count": 1,
      "last_scan": "2024-01-01T12:00:00Z",
      "last_error": ""
    }
  ],
  "scan_count": 150,
  "found_count": 12,
  "bought_count": 3,
  "failed_count": 1,
  "last_error": ""
}
```

**응답 필드:**

| 필드 | 타입 | 설명 |
|------|------|------|
| `running` | boolean | 엔진 실행 여부 |
| `start_time` | string (ISO 8601) | 엔진 시작 시간 |
| `servers` | array | 서버별 상태 배열 |
| `servers[].host` | string | 서버 호스트명 |
| `servers[].scan_count` | integer | 스캔 횟수 |
| `servers[].found_count` | integer | 발견된 항공기 수 |
| `servers[].bought_count` | integer | 구매 성공 수 |
| `servers[].failed_count` | integer | 구매 실패 수 |
| `servers[].last_scan` | string (ISO 8601) | 마지막 스캔 시간 |
| `servers[].last_error` | string | 마지막 오류 메시지 |
| `scan_count` | integer | 전체 스캔 횟수 |
| `found_count` | integer | 전체 발견 수 |
| `bought_count` | integer | 전체 구매 수 |
| `failed_count` | integer | 전체 실패 수 |
| `last_error` | string | 전체 마지막 오류 |

---

### 엔진 시작

```
POST /api/control/start
```

모니터링 엔진을 시작합니다.

**응답:**

```json
{
  "status": "started"
}
```

**오류 응답 (이미 실행 중):**

```json
{
  "error": "engine is already running"
}
```

상태 코드: `409 Conflict`

---

### 엔진 중지

```
POST /api/control/stop
```

모니터링 엔진을 중지합니다.

**응답:**

```json
{
  "status": "stopped"
}
```

---

### 규칙 리로드

```
POST /api/reload
```

설정 파일에서 규칙을 다시 로드합니다. 엔진을 중지하지 않고 규칙을 업데이트합니다.

**응답:**

```json
{
  "status": "reloaded"
}
```

## 설정

### 설정 조회

```
GET /api/config
```

전체 설정을 반환합니다.

**응답 예시:**

```json
{
  "servers": [
    {
      "host": "free1",
      "base_url": "https://free1.airlinesim.aero"
    }
  ],
  "auths": [
    {
      "username": "user@example.com",
      "password": "****",
      "session_file": "session.json"
    }
  ],
  "monitor": {
    "interval": 30,
    "jitter": 10,
    "request_timeout": 30,
    "min_balance": 1000000
  },
  "notifier": {
    "console": true,
    "discord_webhook": ""
  },
  "webui": {
    "enabled": true,
    "host": "0.0.0.0",
    "port": 9090
  },
  "rules": []
}
```

---

### 설정 업데이트

```
PUT /api/config
```

전역 설정을 업데이트합니다. 규칙은 이 엔드포인트로 업데이트되지 않으며, 별도의 규칙 API를 사용합니다.

**요청 본문 예시:**

```json
{
  "servers": [
    {
      "host": "free1",
      "base_url": "https://free1.airlinesim.aero"
    }
  ],
  "auths": [
    {
      "username": "user@example.com",
      "password": "new-password",
      "session_file": "session.json"
    }
  ],
  "monitor": {
    "interval": 45,
    "jitter": 15,
    "request_timeout": 30,
    "min_balance": 2000000
  },
  "notifier": {
    "console": true,
    "discord_webhook": ""
  },
  "webui": {
    "enabled": true,
    "host": "127.0.0.1",
    "port": 9090
  }
}
```

**참고:** 요청에 포함된 필드만 업데이트됩니다. 보내지 않은 필드는 현재 값을 유지합니다.

**응답:** 업데이트된 전체 설정 객체

## 규칙 관리

### 규칙 목록 조회

```
GET /api/rules
```

모든 구매 규칙을 배열로 반환합니다.

**응답 예시:**

```json
[
  {
    "name": "A320-200 Heavy Lease Monitor",
    "enabled": true,
    "priority": 10,
    "server_id": 0,
    "auth_id": 0,
    "match": {
      "family_id": "1200300",
      "type_id": "",
      "types": [],
      "price_range": {
        "min": 0,
        "max": 5000000
      },
      "max_age": 20,
      "max_cycles": 50000,
      "condition_min": 50,
      "offer_types": ["auction", "immediate"],
      "financing": ["cash", "credit", "lease"],
      "sort_by": "price_asc"
    },
    "action": {
      "auto_buy": false,
      "snatch": false,
      "max_bid_increment": 100000
    }
  }
]
```

---

### 새 규칙 생성

```
POST /api/rules
```

새 구매 규칙을 생성합니다.

**요청 본문:**

```json
{
  "name": "새 규칙",
  "enabled": true,
  "priority": 10,
  "server_id": 0,
  "auth_id": 0,
  "match": {
    "family_id": "",
    "type_id": "",
    "types": ["Airbus A320-200 heavy"],
    "price_range": {
      "min": 500000,
      "max": 5000000
    },
    "max_age": 15,
    "max_cycles": 30000,
    "condition_min": 70,
    "offer_types": ["auction", "immediate"],
    "financing": ["cash", "credit", "lease"],
    "sort_by": "price_asc"
  },
  "action": {
    "auto_buy": true,
    "snatch": false,
    "max_bid_increment": 100000
  }
}
```

**응답:** 생성된 규칙 객체 (상태 코드: `201 Created`)

---

### 특정 규칙 조회

```
GET /api/rules/{id}
```

`{id}`는 규칙의 배열 인덱스입니다.

**응답:** 규칙 객체

**오류 응답:**

```json
{
  "error": "rule not found"
}
```

상태 코드: `404 Not Found`

---

### 규칙 수정

```
PUT /api/rules/{id}
```

특정 규칙을 수정합니다.

**요청 본문:** 전체 규칙 객체

**응답:** 업데이트된 규칙 객체

---

### 규칙 삭제

```
DELETE /api/rules/{id}
```

특정 규칙을 삭제합니다.

**응답:**

```json
{
  "status": "deleted"
}
```

---

### 규칙 활성화 전환

```
PATCH /api/rules/{id}/toggle
```

규칙의 `enabled` 상태를 반전시킵니다.

**응답:** 전환된 규칙 객체

---

### 규칙 순서 변경

```
PUT /api/rules/reorder
```

규칙의 순서를 변경합니다.

**요청 본문:**

```json
[2, 0, 1]
```

배열의 각 요소는 새 순서의 규칙 인덱스를 나타냅니다. 위 예시는: `rules[2]`, `rules[0]`, `rules[1]` 순서로 재정렬합니다.

**응답:** 재정렬된 규칙 배열

## 데이터

### 항공기 데이터 조회

```
GET /api/aircraft-data
```

UI 드롭다운에 사용되는 항공기 패밀리 및 유형 데이터를 반환합니다.

**응답 예시:**

```json
{
  "families": [
    {
      "id": "1200300",
      "name": "A320 / A321"
    },
    {
      "id": "2200400",
      "name": "737-600/700/800/900"
    }
  ],
  "types": [
    {
      "id": "16",
      "name": "Airbus A320-200 heavy"
    },
    {
      "id": "97",
      "name": "Boeing 737-800 BGW"
    }
  ]
}
```

## 상태 코드

| 상태 코드 | 설명 |
|-----------|------|
| `200 OK` | 요청 성공 |
| `201 Created` | 리소스 생성 성공 |
| `400 Bad Request` | 잘못된 요청 (JSON 파싱 실패, 유효성 검사 실패) |
| `404 Not Found` | 리소스를 찾을 수 없음 |
| `409 Conflict` | 충돌 (이미 실행 중인 엔진) |
| `500 Internal Server Error` | 서버 내부 오류 |

## 오류 응답 형식

모든 오류 응답은 다음 형식을 따릅니다:

```json
{
  "error": "오류 메시지"
}
```

## 관련 문서

- [웹 UI 가이드](/ko/guide/webui) - 웹 인터페이스 사용법
- [설정 참조](/ko/reference/config) - 설정 필드 상세
- [아키텍처](/ko/reference/architecture) - 시스템 구조
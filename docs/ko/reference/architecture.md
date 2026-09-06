# 아키텍처

AirlineSim Autobuy는 모듈식 파이프라인 아키텍처를 기반으로 설계되었습니다. 각 컴포넌트는 독립적인 책임을 가지며, 명확한 인터페이스를 통해 상호 작용합니다.

## 파이프라인 아키텍처

```
┌─────────────────────────────────────────────────────────────────┐
│                        엔진 (Engine)                            │
│  Orchestrates the entire monitoring pipeline                    │
│                                                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────┐  │
│  │  Auth    │  │Collector │  │  Parser  │  │ Rules Engine  │  │
│  │ (인증)   │→ │(수집)    │→ │(파싱)    │→ │(규칙 평가)    │  │
│  └──────────┘  └──────────┘  └──────────┘  └───────┬───────┘  │
│                                                      │          │
│                                                      ▼          │
│  ┌──────────┐  ┌──────────┐  ┌───────────────────────────────┐  │
│  │Notifier  │← │ Executor │← │    MatchResult (일치 결과)    │  │
│  │(알림)    │  │ (실행)   │  │                               │  │
│  └──────────┘  └──────────┘  └───────────────────────────────┘  │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Config Store (설정 저장소) ←→ Web UI (웹 관리 인터페이스) │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## 컴포넌트 개요

### 1. Config Store (`internal/config/`)

설정 저장소는 애플리케이션 설정을 관리하는 중앙 컴포넌트입니다.

**주요 책임:**
- YAML 설정 파일 읽기 및 쓰기
- 스레드 안전한 설정 접근 (`sync.RWMutex`)
- 설정 변경 핫 리로드 (파일 폴링)
- 레거시 설정 형식 자동 마이그레이션
- 기본 설정 생성

**핵심 타입:**

| 타입 | 설명 |
|------|------|
| `Config` | 최상위 설정 구조체 |
| `ServerConfig` | 게임 서버 연결 정보 |
| `AuthConfig` | 인증 계정 정보 |
| `MonitorConfig` | 모니터링 동작 설정 |
| `NotifierConfig` | 알림 채널 설정 |
| `WebUIConfig` | 웹 UI 설정 |
| `RuleConfig` | 구매 규칙 |
| `MatchConfig` | 규칙 일치 조건 |
| `ActionConfig` | 규칙 실행 작업 |

**파일:**
- `G:/GitHub/airlinesim-autobuy/internal/config/config.go` - 설정 저장소
- `G:/GitHub/airlinesim-autobuy/internal/config/types.go` - 설정 타입 정의
- `G:/GitHub/airlinesim-autobuy/internal/config/watcher.go` - 파일 변경 감지

### 2. Auth (`internal/auth/`)

인증 모듈은 AirlineSim 게임 서버에 대한 인증 및 세션 관리를 담당합니다.

**주요 책임:**
- `sar.simulogics.games` API를 통한 로그인
- 세션 쿠키(`as-sid`) 관리
- 세션 파일 저장 및 복원
- 세션 상태 헬스체크
- 자동 재인증

**보안 기능:**
- URL 유효성 검사 (허용된 도메인만 접근)
- 사설 IP/로컬호스트 차단
- `secureTransport` 래퍼를 통한 요청 검증

**파일:**
- `G:/GitHub/airlinesim-autobuy/internal/auth/auth.go`

### 3. Client (`internal/client/`)

HTTP 클라이언트는 요청 속도 제한, 재시도 로직, URL 검증을 제공합니다.

**주요 기능:**
- **Rate Limiting**: 토큰 버킷 패턴을 사용한 요청 간격 제어
- **Jitter**: 각 요청에 무작위 지연 추가
- **Retry**: 서버 오류(5xx, 429) 시 최대 3회 재시도
- **URL Validation**: 보안 검증 (사설 IP, 로컬호스트, 예약 도메인 차단)

**파일:**
- `G:/GitHub/airlinesim-autobuy/internal/client/client.go`

### 4. Collector (`internal/collector/`)

수집기는 게임 서버의 항공기 시장 페이지를 가져옵니다.

**주요 책임:**
- Wicket 시장 페이지 URL 발견
- 필터 파라미터를 적용한 페이지 요청
- 페이지 HTML 바이트 반환

**필터 파라미터:**
- `FamilyID`: 항공기 패밀리 ID
- `TypeID`: 항공기 유형 ID
- `SortBy`: 정렬 기준

**파일:**
- `G:/GitHub/airlinesim-autobuy/internal/collector/collector.go`

### 5. Parser (`internal/parser/`)

파서는 Wicket HTML 페이지를 구조화된 항공기 데이터로 변환합니다.

**주요 책임:**
- HTML 테이블 행 파싱
- 다양한 숫자 형식 처리 (US/유럽)
- 항공기 속성 추출 (유형, 가격, 기령, 상태, 사이클 등)
- 고유 ID 생성

**파싱되는 항공기 속성:**

| 속성 | 타입 | 설명 |
|------|------|------|
| `Type` | string | 항공기 유형명 |
| `Family` | string | 항공기 패밀리명 |
| `Age` | int | 기령 (년) |
| `Cycles` | int | 비행 사이클 |
| `Condition` | float64 | 상태 (0-100) |
| `Price` | float64 | 기준 가격 (AS$) |
| `LeaseRate` | float64 | 주간 리스 요금 |
| `OfferType` | string | "auction" 또는 "immediate" |
| `Financing` | []string | 사용 가능한 금융 옵션 |

**파일:**
- `G:/GitHub/airlinesim-autobuy/internal/parser/parser.go`

### 6. Rules Engine (`internal/rules/`)

규칙 엔진은 항공기 매물을 사용자 정의 규칙과 비교하여 평가합니다.

**주요 책임:**
- 규칙 기반 항공기 매칭
- 점수 계산 (가격, 상태, 기령, 우선순위 기반)
- 최적 매칭 결정

**매칭 조건:**
- 항공기 유형 (`types`)
- 가격 범위 (`price_range`)
- 최대 기령 (`max_age`)
- 최대 사이클 (`max_cycles`)
- 최소 상태 (`condition_min`)
- 제공 유형 (`offer_types`)
- 금융 옵션 (`financing`)

**점수 시스템:**
- 가격 점수 (50점): 최대 가격 대비 낮을수록 높음
- 상태 점수 (20점): 최소 상태 대비 높을수록 높음
- 기령 점수 (20점): 최대 기령 대비 낮을수록 높음
- 즉시 구매 보너스 (10점)
- 우선순위 보너스 (가변)

**파일:**
- `G:/GitHub/airlinesim-autobuy/internal/rules/engine.go`

### 7. Executor (`internal/executor/`)

실행기는 실제 항공기 구매 또는 입찰을 처리합니다.

**주요 책임:**
- 즉시 구매 요청
- 경매 입찰
- 가격 확인

**파일:**
- `G:/GitHub/airlinesim-autobuy/internal/executor/executor.go`

### 8. Notifier (`internal/notifier/`)

알림 모듈은 이벤트를 다양한 채널로 전송합니다.

**이벤트 유형:**

| 이벤트 | 설명 |
|--------|------|
| `AircraftFound` | 규칙 일치 항공기 발견 |
| `PurchaseMade` | 구매 성공 |
| `PurchaseFailed` | 구매 실패 |
| `BidPlaced` | 입찰 완료 |
| `Error` | 오류 발생 |
| `Info` | 정보 메시지 |

**알림 채널:**
- 콘솔 (구조화된 로그)
- Discord 웹훅 (준비 중)

**파일:**
- `G:/GitHub/airlinesim-autobuy/internal/notifier/notifier.go`

### 9. Engine (`internal/engine/`)

엔진은 모든 컴포넌트를 조정하는 오케스트레이터입니다.

**주요 책임:**
- 서버별 모니터링 루프 관리
- 세션 상태 확인 및 재인증
- 수집 → 파싱 → 규칙 평가 → 실행 파이프라인 조정
- 중복 항공기 감지 (`seen` 맵)
- 전역 상태 추적

**파일:**
- `G:/GitHub/airlinesim-autobuy/internal/engine/engine.go`

### 10. Web UI (`internal/webui/`)

웹 UI는 REST API와 Vue 3 SPA를 제공하는 내장 HTTP 서버입니다.

**주요 책임:**
- REST API 엔드포인트 제공
- Vue 3 프론트엔드 정적 파일 서빙
- CORS 미들웨어 (개발 모드 지원)
- SPA 폴백 (Vue Router 지원)

**사용된 라이브러리:**
- `chi/v5` - HTTP 라우터
- Vue 3 + Vite - 프론트엔드
- `embed` - Go 내장 파일 시스템

**파일:**
- `G:/GitHub/airlinesim-autobuy/internal/webui/server.go` - 서버 및 API 핸들러
- `G:/GitHub/airlinesim-autobuy/internal/webui/frontend/` - Vue 3 프론트엔드

## 데이터 흐름

### 메인 파이프라인

```
[Config Store] → 설정 로드
       ↓
[Engine] → 각 서버별 고루틴 시작
       ↓
[Auth] → 서버 로그인 (세션 수립)
       ↓
[Collector] → 시장 페이지 URL 발견
       ↓
[Monitor Loop] (interval + jitter 주기)
       ↓
[Collector.Fetch] → HTTP GET 시장 페이지
       ↓
[Parser.Parse] → HTML → AircraftOffer[]
       ↓
[Rules Engine.Evaluate] → 규칙 매칭
       ↓
[Executor.Execute] → 구매/입찰 실행
       ↓
[Notifier] → 알림 전송
       ↓
[Config Store] → 변경 감지 (핫 리로드)
```

### 동시성 모델

```
┌─────────────────────────────────────┐
│         메인 고루틴 (Main)          │
│  - 설정 로드                        │
│  - Web UI 서버 시작                 │
│  - 시그널 처리 (SIGINT/SIGTERM)     │
└──────────────┬──────────────────────┘
               │
     ┌─────────┴─────────┐
     │                   │
     ▼                   ▼
┌─────────────┐   ┌─────────────┐
│ 서버 1 루프  │   │ 서버 2 루프  │  ... 각 서버별 고루틴
│ (고루틴)     │   │ (고루틴)     │
│             │   │             │
│ ticker      │   │ ticker      │
│ ↓           │   │ ↓           │
│ fetch       │   │ fetch       │
│ parse       │   │ parse       │
│ evaluate    │   │ evaluate    │
│ execute     │   │ execute     │
└─────────────┘   └─────────────┘
```

**동시성 특징:**
- 각 서버는 독립적인 고루틴에서 실행됩니다.
- `sync.RWMutex`로 설정 저장소와 상태 보호
- `context.WithCancel`로 graceful shutdown 지원
- `seen` 맵은 `sync.Mutex`로 보호 (중복 탐지)
- Rate limiter는 채널 기반 토큰 버킷 사용

## 의존성 그래프

```
cmd/autobuy/main.go
  ├── internal/config
  │     └── gopkg.in/yaml.v3
  ├── internal/engine
  │     ├── internal/auth
  │     ├── internal/client
  │     ├── internal/collector
  │     ├── internal/parser
  │     │     └── github.com/PuerkitoBio/goquery
  │     ├── internal/rules
  │     ├── internal/executor
  │     └── internal/notifier
  └── internal/webui
        └── github.com/go-chi/chi/v5
```

## 관련 문서

- [API 참조](/ko/reference/api) - REST API 엔드포인트
- [설정 참조](/ko/reference/config) - 설정 필드 상세
- [기여하기](/ko/reference/contributing) - 개발 환경 설정
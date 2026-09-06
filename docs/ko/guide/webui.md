# 웹 UI

AirlineSim Autobuy는 Vue 3 + Vite로 구축된 웹 관리 인터페이스를 Go 바이너리에 내장하고 있습니다. 브라우저를 통해 엔진을 제어하고, 규칙을 관리하며, 실시간 로그를 확인할 수 있습니다.

## 접속

기본 설정에서 Web UI는 `http://localhost:9090`에서 접속할 수 있습니다.

```bash
# 기본 주소
http://localhost:9090

# 원격 접속 (host가 0.0.0.0인 경우)
http://your-server-ip:9090
```

## 대시보드 (Dashboard)

대시보드는 엔진의 전반적인 상태를 한눈에 보여줍니다.

### 엔진 상태 배너

- **Running**: 엔진이 실행 중이며 시장을 모니터링하고 있습니다.
- **Stopped**: 엔진이 중지된 상태입니다.

### 엔진 제어 버튼

| 버튼 | 동작 |
|------|------|
| **Start Engine** | 모니터링 엔진 시작 |
| **Stop Engine** | 모니터링 엔진 중지 |

### 통계 카드

| 메트릭 | 설명 |
|--------|------|
| **Market Scans** | 수행된 시장 스캔 횟수 |
| **Aircraft Found** | 규칙에 일치하는 항공기 발견 수 |
| **Purchased** | 성공적으로 구매한 항공기 수 |
| **Failed** | 실패한 구매 시도 수 |

### 기타 정보

- **Last Engine Error**: 마지막 엔진 오류 메시지
- **Activity**: 마지막 스캔 시간

대시보드는 3초마다 자동으로 상태를 갱신합니다.

## 규칙 페이지 (Rules)

규칙 페이지에서는 구매 규칙을 생성, 수정, 삭제, 활성화/비활성화할 수 있습니다.

### 규칙 목록

각 규칙은 카드 형태로 표시되며 다음 정보를 포함합니다:

- **규칙 이름** 및 메타 정보 (우선순위, 항공기 유형)
- **활성/비활성** 배지
- **규칙 상세 정보**: 가격 범위, 최대 기령, 최소 상태, 최대 사이클, 자동 구매 상태, 스내치 상태

### 규칙 작업

| 작업 | 설명 |
|------|------|
| **토글 버튼** (✓) | 규칙 활성화/비활성화 전환 |
| **편집 버튼** (연필 아이콘) | 규칙 수정 모달 열기 |
| **삭제 버튼** (휴지통 아이콘) | 규칙 삭제 (확인 후 실행) |
| **New Rule** 버튼 | 새 규칙 생성 모달 열기 |

### 규칙 생성/편집 폼

**기본 정보:**
- **Rule Name**: 규칙 이름 (필수)
- **Server**: 대상 서버 선택
- **Account**: 사용할 계정 선택

**항공기 조건:**
- **Aircraft Family**: 항공기 패밀리 선택
- **Aircraft Type**: 특정 항공기 유형 선택
- **Sort By**: 시장 정렬 기준

**가격 및 상태 조건:**
- **Min / Max Price**: 가격 범위 (AS$)
- **Max Age**: 최대 기령 (년)
- **Max Cycles**: 최대 사이클
- **Min Condition**: 최소 상태 (%)

**옵션:**
- **Offer Types**: 제공 유형 선택 (Auction, Immediate)
- **Financing Options**: 금융 옵션 선택 (Cash, Credit, Lease)

**작업 설정:**
- **Priority**: 규칙 우선순위
- **Max Bid Increment**: 최대 입찰 증가액 (AS$)
- **Enabled**: 규칙 활성화
- **Auto Buy**: 자동 구매 활성화
- **Snatch**: 스내치 모드 활성화

## 설정 페이지 (Settings)

설정 페이지에서는 애플리케이션의 전역 설정을 관리할 수 있습니다.

### 서버 연결

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
```

- **+ Add Server** 버튼으로 새 서버 추가
- 각 서버의 `host`와 `base_url` 설정
- 서버 삭제 버튼 (2개 이상일 때 활성화)

### 인증 계정

```yaml
auths:
  - username: account@example.com
    password: ""
    session_file: session.json
```

- **+ Add Account** 버튼으로 새 계정 추가
- 각 계정의 사용자명, 비밀번호 설정
- 세션 파일은 자동으로 관리됩니다

### 모니터링 설정

| 필드 | 설명 |
|------|------|
| **Poll Interval** | 서버 스캔 간격 (초, 최소 15) |
| **Jitter** | 무작위 지연 시간 (초) |
| **Request Timeout** | HTTP 요청 타임아웃 (초) |
| **Minimum Account Balance** | 구매 후 최소 유지 잔액 (AS$) |

### 알림 설정

| 필드 | 설명 |
|------|------|
| **Console Output** | 콘솔 알림 활성화/비활성화 |
| **Discord Webhook** | Discord 웹훅 URL (선택 사항) |

### Web UI 설정

| 필드 | 설명 |
|------|------|
| **Host** | 웹 서버 바인딩 주소 |
| **Port** | 웹 서버 포트 |

### 저장

설정을 변경한 후 **Save Settings** 버튼을 클릭하면 설정 파일에 저장되고 엔진이 자동으로 리로드됩니다.

## 로그 페이지 (Logs)

로그 페이지는 엔진의 실시간 활동 로그를 제공합니다.

### 기능

- **실시간 로그 스트리밍**: Server-Sent Events (SSE)를 통한 실시간 로그 수신
- **로그 필터링**: 키워드로 로그 필터링
- **자동 스크롤**: 새 로그 자동 스크롤 (토글 가능)
- **로그 지우기**: Clear 버튼으로 로그 초기화
- **최대 1000개 항목**: 로그는 최대 1000개까지 유지되며 초과 시 오래된 항목부터 제거됩니다

### 로그 색상

| 로그 레벨 | 색상 |
|----------|------|
| ERROR | 빨간색 |
| WARN | 노란색 |
| INFO | 기본 색상 |

## API 엔드포인트

Web UI는 다음 REST API 엔드포인트를 제공합니다:

### 상태 및 제어

| 메서드 | 엔드포인트 | 설명 |
|--------|----------|------|
| GET | `/api/status` | 엔진 상태 조회 |
| POST | `/api/control/start` | 엔진 시작 |
| POST | `/api/control/stop` | 엔진 중지 |
| POST | `/api/reload` | 규칙 리로드 |

### 설정

| 메서드 | 엔드포인트 | 설명 |
|--------|----------|------|
| GET | `/api/config` | 설정 조회 |
| PUT | `/api/config` | 설정 업데이트 |

### 규칙 관리

| 메서드 | 엔드포인트 | 설명 |
|--------|----------|------|
| GET | `/api/rules` | 규칙 목록 조회 |
| POST | `/api/rules` | 새 규칙 생성 |
| GET | `/api/rules/{id}` | 특정 규칙 조회 |
| PUT | `/api/rules/{id}` | 특정 규칙 수정 |
| DELETE | `/api/rules/{id}` | 특정 규칙 삭제 |
| PATCH | `/api/rules/{id}/toggle` | 규칙 활성화 전환 |
| PUT | `/api/rules/reorder` | 규칙 순서 변경 |

### 데이터

| 메서드 | 엔드포인트 | 설명 |
|--------|----------|------|
| GET | `/api/aircraft-data` | 항공기 데이터 (패밀리, 유형) 조회 |

## 프론트엔드 개발

프론트엔드는 `internal/webui/frontend/` 디렉토리에 위치하며, Vue 3 + TypeScript로 작성되었습니다.

### 개발 서버 실행

```bash
cd internal/webui/frontend
npm run dev
```

개발 서버는 `http://localhost:5173`에서 실행되며, API 요청은 Go 백엔드(`localhost:9090`)로 프록시됩니다. CORS는 개발 모드에서 자동으로 처리됩니다.

### 프로덕션 빌드

```bash
make frontend-build
```

빌드된 파일은 `internal/webui/frontend/dist/`에 생성되며, Go 바이너리에 내장됩니다.

## 관련 문서

- [API 참조](/ko/reference/api) - 모든 REST API 엔드포인트 상세
- [설정 가이드](/ko/guide/configuration) - Web UI 설정 옵션
- [빠른 시작](/ko/guide/quickstart) - 설치 및 실행
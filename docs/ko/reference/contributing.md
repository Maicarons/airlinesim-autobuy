# 기여하기

AirlineSim Autobuy에 기여해 주셔서 감사합니다! 이 문서는 개발 환경 설정, 코드 스타일, 테스트, 풀 리퀘스트 프로세스를 안내합니다.

## 개발 환경 설정

### 1. 필수 도구 설치

| 도구 | 버전 | 설치 방법 |
|------|------|----------|
| Go | 1.25.5+ | [go.dev/dl](https://go.dev/dl/) |
| Node.js | 18+ | [nodejs.org](https://nodejs.org/) |
| npm | 9+ | Node.js와 함께 설치 |
| Git | 최신 | [git-scm.com](https://git-scm.com/) |

### 2. 선택 도구 설치

```bash
# Air (핫 리로드)
go install github.com/air-verse/air@latest

# golangci-lint (린터)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 3. 저장소 클론

```bash
git clone https://github.com/Maicarons/airlinesim-autobuy.git
cd airlinesim-autobuy
```

### 4. 의존성 설치

```bash
# Go 의존성
go mod download

# 프론트엔드 의존성
make frontend-install
```

### 5. 개발 서버 실행

```bash
# 백엔드 개발 서버 (핫 리로드)
make dev

# 별도 터미널에서 프론트엔드 개발 서버
cd internal/webui/frontend && npm run dev
```

## 프로젝트 구조

```
airlinesim-autobuy/
├── cmd/
│   └── autobuy/
│       └── main.go              # 애플리케이션 진입점
├── configs/
│   └── config.yaml              # 설정 파일
├── internal/
│   ├── auth/                    # AirlineSim 인증
│   │   └── auth.go
│   ├── client/                  # HTTP 클라이언트 (속도 제한, 재시도)
│   │   └── client.go
│   ├── collector/               # 시장 페이지 수집
│   │   └── collector.go
│   ├── config/                  # 설정 관리
│   │   ├── config.go
│   │   ├── types.go
│   │   └── watcher.go
│   ├── engine/                  # 모니터링 엔진 (오케스트레이터)
│   │   └── engine.go
│   ├── executor/                # 구매 실행
│   │   └── executor.go
│   ├── marketdata/              # 항공기 참조 데이터
│   │   └── data.go
│   ├── notifier/                # 알림 전송
│   │   └── notifier.go
│   ├── parser/                  # HTML 파싱
│   │   └── parser.go
│   ├── rules/                   # 규칙 엔진
│   │   └── engine.go
│   └── webui/                   # 웹 관리 인터페이스
│       ├── server.go            # HTTP 서버 및 API 핸들러
│       └── frontend/            # Vue 3 프론트엔드
│           ├── src/
│           │   ├── api/         # API 클라이언트
│           │   ├── views/       # Vue 컴포넌트
│           │   └── ...
│           └── ...
├── docs/                        # VitePress 문서
├── Makefile                     # 빌드 및 개발 명령어
├── go.mod
└── go.sum
```

## 코드 스타일

### Go

- **포맷팅**: `go fmt` 또는 `gofumpt` 사용
- **린팅**: `golangci-lint run ./...` 통과 필수
- **네이밍**: Go 관례 준수 (camelCase, PascalCase)
- **주석**: 모든 exported 식별자에 주석 작성
- **에러 처리**: 에러는 적절히 감싸서 반환 (`fmt.Errorf("...: %w", err)`)
- **로깅**: `slog` 패키지 사용

### Go 코드 예시

```go
// Package example는 예시 패키지입니다.
package example

import (
    "fmt"
    "log/slog"
)

// Config는 예시 설정 구조체입니다.
type Config struct {
    Name  string `json:"name" yaml:"name"`
    Count int    `json:"count" yaml:"count"`
}

// New은 새로운 Example을 생성합니다.
func New(cfg Config) *Example {
    return &Example{
        name:  cfg.Name,
        count: cfg.Count,
    }
}

// DoSomething은 작업을 수행합니다.
func (e *Example) DoSomething() error {
    if e.name == "" {
        return fmt.Errorf("name is required")
    }
    slog.Info("doing something", "name", e.name, "count", e.count)
    return nil
}
```

### TypeScript / Vue

- **타입**: 가능한 경우 TypeScript 타입 명시
- **컴포넌트**: `<script setup lang="ts">` 구문 사용
- **스타일**: `<style scoped>` 사용
- **API 호출**: `src/api/index.ts`의 함수 사용

### Vue 코드 예시

```typescript
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getData } from '../api'

interface Data {
  id: number
  name: string
}

const data = ref<Data | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    data.value = await getData()
  } catch (err) {
    console.error('Failed to load data:', err)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div v-if="loading">Loading...</div>
  <div v-else-if="data">
    <h1>{{ data.name }}</h1>
  </div>
</template>

<style scoped>
/* 컴포넌트 전용 스타일 */
</style>
```

## 테스트

### Go 테스트 실행

```bash
# 모든 테스트 실행
make test

# 또는 직접
go test ./... -v

# 특정 패키지 테스트
go test ./internal/rules/... -v

# 커버리지 확인
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 테스트 작성 가이드라인

- 테스트 파일은 `_test.go` 접미사 사용
- 테이블 기반 테스트 선호
- 테스트 함수명은 `TestXxx` 형식
- 외부 API에 의존하는 테스트는 모킹 고려

```go
func TestEngine_Evaluate(t *testing.T) {
    tests := []struct {
        name     string
        aircraft *parser.AircraftOffer
        rule     config.RuleConfig
        want     bool
    }{
        {
            name: "matches type and price",
            aircraft: &parser.AircraftOffer{
                Type:  "Airbus A320-200 heavy",
                Price: 2000000,
            },
            rule: config.RuleConfig{
                Match: config.MatchConfig{
                    Types:      []string{"Airbus A320-200 heavy"},
                    PriceRange: config.PriceRange{Max: 5000000},
                },
                Action: config.ActionConfig{AutoBuy: true},
            },
            want: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            engine := New([]config.RuleConfig{tt.rule})
            result := engine.BestMatch(tt.aircraft)
            if (result != nil) != tt.want {
                t.Errorf("BestMatch() = %v, want %v", result != nil, tt.want)
            }
        })
    }
}
```

### 프론트엔드 테스트

```bash
cd internal/webui/frontend
npm run test
```

## 풀 리퀘스트

### PR 프로세스

1. **이슈 생성**: 변경 사항에 대한 이슈를 먼저 생성하세요.
2. **포크 및 브랜치**: `main` 브랜치에서 기능 브랜치를 생성하세요.
3. **변경 사항 커밋**: 명확한 커밋 메시지를 작성하세요.
4. **테스트**: 모든 테스트가 통과하는지 확인하세요.
5. **린트**: `golangci-lint`가 통과하는지 확인하세요.
6. **PR 생성**: 변경 사항을 설명하는 PR을 생성하세요.

### 커밋 메시지 규칙

[Conventional Commits](https://www.conventionalcommits.org/) 형식을 따릅니다:

```
<type>(<scope>): <description>

[optional body]
```

**타입:**
- `feat`: 새 기능
- `fix`: 버그 수정
- `docs`: 문서 변경
- `style`: 코드 포맷팅 (로직 변경 없음)
- `refactor`: 리팩토링
- `test`: 테스트 추가/수정
- `chore`: 빌드, 의존성, 기타

**예시:**
```
feat(rules): add auction-only offer type filter
fix(engine): handle session expiry during long runs
docs(config): add multi-auth configuration example
```

### PR 체크리스트

- [ ] 모든 기존 테스트가 통과
- [ ] 새 기능에 테스트가 추가됨
- [ ] 코드가 `golangci-lint`를 통과
- [ ] 문서가 업데이트됨 (필요한 경우)
- [ ] 커밋 메시지가 명확함
- [ ] 변경 사항이 너무 크지 않음 (여러 PR로 분할 가능)

## 문서

문서는 VitePress로 작성됩니다. 문서 변경 사항이 있는 경우:

```bash
# 문서 개발 서버 실행
make docs-dev

# 문서 빌드
make docs-build
```

문서는 `docs/` 디렉토리에 위치하며, 한국어 문서는 `docs/ko/`에 있습니다.

## 릴리스 프로세스

1. `CHANGELOG.md` 업데이트
2. 버전 태그 생성 (`git tag v1.0.0`)
3. GitHub Release 생성
4. 모든 플랫폼용 바이너리 빌드 (`make build-all`)

## 행동 강령

이 프로젝트는 기여자가 존중하고 협력적인 환경을 조성할 것을 기대합니다. 모든 참여자는 다음을 준수해야 합니다:

- 포용적인 언어 사용
- 다양한 관점 존중
- 건설적인 피드백 수용
- 공동체 최선의 행동

## 문의

- **GitHub Issues**: [https://github.com/Maicarons/airlinesim-autobuy/issues](https://github.com/Maicarons/airlinesim-autobuy/issues)
- **GitHub Discussions**: [https://github.com/Maicarons/airlinesim-autobuy/discussions](https://github.com/Maicarons/airlinesim-autobuy/discussions)

## 관련 문서

- [아키텍처](/ko/reference/architecture) - 시스템 구조
- [API 참조](/ko/reference/api) - REST API 엔드포인트
- [설정 참조](/ko/reference/config) - 설정 필드 상세
# FAQ

## 일반 질문

### AirlineSim Autobuy란 무엇인가요?

AirlineSim Autobuy는 AirlineSim 게임의 중고 항공기 시장을 실시간으로 모니터링하고, 사용자가 정의한 규칙에 따라 자동으로 항공기를 구매하는 도구입니다. Go로 작성되었으며, Vue 3 기반의 웹 관리 인터페이스를 내장하고 있습니다.

### 이 도구는 AirlineSim 이용약관에 위배되나요?

이 도구는 게임 서버에 공개적으로 접근 가능한 웹 페이지를 읽기 전용으로 스캔하고, 사용자가 설정한 조건에 따라 자동화된 구매를 수행합니다. 그러나 각 게임의 이용약관은 자동화 도구 사용에 대해 다를 수 있으므로, 사용 전에 AirlineSim의 이용약관을 확인하는 것을 권장합니다.

### 무료인가요?

네, 이 프로젝트는 MIT 라이선스로 제공되는 오픈 소스 소프트웨어입니다.

### 어떤 언어를 지원하나요?

애플리케이션 자체는 영어 UI를 제공하지만, 문서는 영어, 중국어(간체), 한국어를 지원합니다.

## 설정

### 설정 파일을 찾을 수 없습니다

애플리케이션을 처음 실행하면 `configs/config.yaml` 파일이 자동으로 생성됩니다. 자동 생성되지 않는 경우:

```bash
# 기본 설정 파일 생성
mkdir -p configs
cat > configs/config.yaml << 'EOF'
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
rules: []
EOF
```

### 계정 인증이 계속 실패합니다

인증 실패의 일반적인 원인:

1. **잘못된 자격 증명**: 사용자명(이메일)과 비밀번호를 확인하세요.
2. **계정 잠김**: 너무 많은 로그인 시도로 계정이 일시적으로 잠겼을 수 있습니다.
3. **서버 다운**: 게임 서버가 유지보수 중이거나 다운되었을 수 있습니다.
4. **네트워크 문제**: 방화벽이 `sar.simulogics.games` API에 대한 접근을 차단하고 있을 수 있습니다.

로그 확인:
```bash
# 로그에서 인증 관련 메시지 확인
journalctl -u airlinesim-autobuy -n 50 | grep -i login
```

### 여러 서버를 어떻게 설정하나요?

```yaml
servers:
  - host: free1
    base_url: https://free1.airlinesim.aero
  - host: free2
    base_url: https://free2.airlinesim.aero
  - host: free3
    base_url: https://free3.airlinesim.aero
```

### 여러 계정을 어떻게 사용하나요?

```yaml
auths:
  - username: account1@example.com
    password: password1
    session_file: session1.json
  - username: account2@example.com
    password: password2
    session_file: session2.json
```

각 규칙에서 `auth_id` 필드로 사용할 계정을 지정합니다.

## 기술 문제

### Web UI에 접속할 수 없습니다

1. **서버 실행 확인**: 애플리케이션이 실행 중인지 확인하세요.
2. **포트 확인**: 기본 포트 9090이 사용 중인지 확인:
   ```bash
   netstat -tlnp | grep 9090
   ```
3. **방화벽 확인**: 방화벽이 포트를 차단하고 있지 않은지 확인하세요.
4. **바인딩 주소 확인**: `host`가 `0.0.0.0` 또는 `127.0.0.1`로 설정되어 있는지 확인하세요.

### 엔진이 시작되지 않습니다

가능한 원인:

1. **잘못된 설정 파일**: YAML 형식이 올바른지 확인하세요.
2. **인증 실패**: 계정 정보가 올바른지 확인하세요.
3. **서버 연결 불가**: 게임 서버에 접근할 수 있는지 확인하세요.
4. **로그 확인**: 자세한 오류 메시지는 로그를 확인하세요.

### 구매가 실행되지 않습니다

1. **규칙 활성화 확인**: 규칙의 `enabled` 필드가 `true`인지 확인하세요.
2. **Auto Buy 확인**: `action.auto_buy`가 `true`인지 확인하세요.
3. **일치 조건 확인**: 항공기가 규칙의 모든 조건을 충족하는지 확인하세요.
4. **잔액 확인**: 계정 잔액이 충분한지 확인하세요.
5. **로그 확인**: 구매 시도와 실패 원인이 로그에 기록됩니다.

### 로그가 너무 많습니다

로그 레벨은 `main.go`에서 설정됩니다. 현재는 `Info` 레벨 이상의 로그만 출력됩니다. 디버그 로그를 보려면 소스 코드에서 로그 레벨을 `slog.LevelDebug`로 변경하세요.

### 설정 파일 변경이 적용되지 않습니다

1. **파일 저장 확인**: 변경 사항을 저장했는지 확인하세요.
2. **핫 리로드 확인**: 변경 사항이 감지되면 로그에 "configuration hot-reloaded" 메시지가 출력됩니다.
3. **수동 리로드**: Web UI의 `/api/reload` 엔드포인트를 호출하여 수동으로 리로드할 수 있습니다.
4. **애플리케이션 재시작**: 최후의 방법으로 애플리케이션을 재시작하세요.

## 문제 해결

### 일반적인 오류 메시지

| 오류 메시지 | 원인 | 해결 방법 |
|------------|------|----------|
| `login failed` | 인증 실패 | 계정 정보 확인 |
| `failed to fetch market page` | 시장 페이지 로드 실패 | 서버 상태 확인, 네트워크 확인 |
| `session expired` | 세션 만료 | 자동 재인증 대기 (1-2분) |
| `failed to parse market page` | HTML 파싱 실패 | 게임 업데이트 확인, 버그 리포팅 |
| `purchase returned status 4xx` | 구매 요청 실패 | 계정 잔액, 자격 확인 |

### 디버깅 팁

1. **로그 레벨 변경**: 개발 중에는 로그 레벨을 `Debug`로 설정하세요:
   ```go
   // main.go
   slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
       Level: slog.LevelDebug,
   })))
   ```

2. **직접 API 테스트**: curl로 API를 직접 테스트하세요:
   ```bash
   curl http://localhost:9090/api/status
   curl http://localhost:9090/api/rules
   ```

3. **세션 파일 확인**: `session.json` 파일이 올바르게 생성되었는지 확인하세요.

4. **네트워크 디버깅**: 게임 서버에 접근할 수 있는지 확인:
   ```bash
   curl -v https://free1.airlinesim.aero
   ```

### 지원 받기

- **GitHub Issues**: [https://github.com/Maicarons/airlinesim-autobuy/issues](https://github.com/Maicarons/airlinesim-autobuy/issues)
- **문서**: 전체 문서는 VitePress 사이트에서 확인할 수 있습니다.

## 관련 문서

- [빠른 시작](/ko/guide/quickstart) - 설치 및 실행
- [설정 가이드](/ko/guide/configuration) - 설정 옵션
- [구매 규칙](/ko/guide/rules) - 규칙 설정
- [배포 가이드](/ko/guide/deployment) - 프로덕션 배포
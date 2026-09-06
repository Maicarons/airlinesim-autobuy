# 参与贡献

感谢你考虑为 AirlineSim Autobuy 贡献代码！本文档将帮助你了解开发流程、代码规范和测试要求。

## 开发环境

### 必需工具

- **Go 1.22+** — 下载：[https://golang.org/dl](https://golang.org/dl)
- **Node.js 18+** — 下载：[https://nodejs.org](https://nodejs.org)
- **Git** — 版本控制
- **Make** — 构建自动化（Windows 可使用 GNU Make 或使用 Go 命令替代）

### 推荐工具

```bash
# Air 热重载工具（开发时自动重启）
make tools

# golangci-lint 代码检查
# 已包含在 make tools 中
```

### 环境配置

```bash
# 克隆仓库
git clone https://github.com/Maicarons/airlinesim-autobuy.git
cd airlinesim-autobuy

# 安装前端依赖
make frontend-install

# 验证构建
make all
```

## 项目结构

```
airlinesim-autobuy/
├── cmd/autobuy/              # 应用程序入口
│   └── main.go
├── internal/
│   ├── auth/                 # 认证与会话管理
│   ├── client/               # 限速 HTTP 客户端
│   ├── collector/            # 市场数据采集
│   ├── config/               # 配置管理与热重载
│   ├── engine/               # 管道编排
│   ├── executor/             # 购买执行
│   ├── marketdata/           # 飞机市场参考数据
│   ├── notifier/             # 通知系统
│   ├── parser/               # HTML 页面解析
│   ├── rules/                # 规则引擎
│   └── webui/                # Web 管理界面
│       ├── server.go
│       └── frontend/         # Vue 3 + Vite 前端
├── configs/                  # 配置文件
├── docs/                     # VitePress 文档
└── Makefile
```

## 开发工作流

### 1. 分支策略

- `main` 分支保持稳定，始终可构建
- 新功能在 `feature/xxx` 分支开发
- Bug 修复在 `fix/xxx` 分支进行
- 提交 Pull Request 到 `main` 分支

### 2. 开发流程

```bash
# 创建功能分支
git checkout -b feature/my-feature

# 开发后端
make run

# 开发前端（另一个终端）
cd internal/webui/frontend && npm run dev

# 运行测试
make test

# 运行代码检查
make lint
```

### 3. 提交规范

我们使用常规的 Git 提交信息格式：

```
<type>(<scope>): <subject>

<body>
```

**类型（type）：**

| 类型 | 说明 |
|------|------|
| `feat` | 新功能 |
| `fix` | Bug 修复 |
| `docs` | 文档更新 |
| `style` | 代码格式调整 |
| `refactor` | 代码重构 |
| `test` | 测试相关 |
| `chore` | 构建/工具链变更 |

**示例：**

```
feat(rules): 支持按飞机系列 ID 过滤

添加 family_id 匹配字段，允许用户通过 Wicket 系列 ID 筛选飞机。
```

## 代码风格

### Go 后端

- 遵循 [Go 官方编码规范](https://go.dev/doc/effective_go)
- 使用 `gofmt` 或 `go fmt` 格式化代码
- 使用 `log/slog` 进行结构化日志记录
- 错误处理：始终检查错误，不要使用 `_` 忽略错误
- 包名使用小写单数形式
- 注释使用英文，公共函数和类型必须有注释

### 命名约定

| 类型 | 约定 | 示例 |
|------|------|------|
| 包名 | 小写，单数 | `config`, `engine` |
| 接口 | 方法名 + `er` 后缀 | `Notifier`, `Parser` |
| 结构体 | 驼峰命名 | `ServerConfig`, `MatchConfig` |
| 变量 | 驼峰命名 | `cfgStore`, `rulesEng` |
| 常量 | 驼峰命名 | `EventAircraftFound` |

### 导入顺序

```
标准库
空行
第三方包
空行
内部包
```

示例：

```go
import (
    "fmt"
    "log/slog"
    "os"

    "github.com/go-chi/chi/v5"
    "gopkg.in/yaml.v3"

    "github.com/Maicarons/airlinesim-autobuy/internal/config"
)
```

### Vue 前端

- 遵循 [Vue 3 风格指南](https://vuejs.org/style-guide/)
- 使用 TypeScript
- 组件使用组合式 API（Composition API）
- 使用 `<script setup>` 语法

## 测试

### 运行测试

```bash
# 运行所有测试
make test

# 运行特定包测试
go test ./internal/rules/... -v

# 运行测试并生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 测试要求

- 新功能应包含单元测试
- Bug 修复应包含防止回归的测试
- 包级别的测试放在 `*_test.go` 文件中
- 使用表驱动测试（Table-driven tests）模式

### 测试示例

```go
func TestMatchRule(t *testing.T) {
    tests := []struct {
        name     string
        aircraft *parser.AircraftOffer
        rule     *config.RuleConfig
        want     bool
    }{
        {
            name: "匹配成功",
            aircraft: &parser.AircraftOffer{
                Type:  "Airbus A320-200 heavy",
                Price: 1000000,
                Age:   5,
            },
            rule: &config.RuleConfig{
                Match: config.MatchConfig{
                    Types: []string{"A320-200"},
                    PriceRange: config.PriceRange{
                        Max: 5000000,
                    },
                },
            },
            want: true,
        },
    }

    engine := New(nil)
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := engine.matchRule(tt.aircraft, tt.rule)
            if (result != nil) != tt.want {
                t.Errorf("matchRule() = %v, want %v", result, tt.want)
            }
        })
    }
}
```

## 拉取请求流程

### 提交 PR 前检查清单

- [ ] 代码通过 `make lint` 检查
- [ ] 所有测试通过 `make test`
- [ ] 新增代码包含适当的测试
- [ ] 公共函数和类型有文档注释
- [ ] 提交信息遵循规范格式
- [ ] 如有必要，更新了文档

### PR 描述模板

```markdown
## 变更内容

简要描述本次 PR 的变更。

## 相关 Issue

Closes #123

## 测试说明

描述如何测试这些变更。

## 截图（如适用）

如果是 UI 变更，请提供前后对比截图。
```

### 审查流程

1. 提交 PR 后，维护者会进行代码审查
2. 可能会要求修改或补充代码
3. 所有对话 resolved 后，PR 会被合并
4. 合并前确保分支是最新的（rebase 到 main）

## 文档

- 文档使用 VitePress 构建，位于 `docs/` 目录
- 文档内容支持多语言（英语、简体中文、韩语）
- 新增功能时请同步更新相关文档
- 文档编写时需要在 `docs/.vitepress/config.ts` 中添加对应侧边栏条目

### 本地预览文档

```bash
cd docs
npm install
npm run dev
```

## 新增功能指南

### 添加新的匹配条件

1. 在 `internal/config/types.go` 的 `MatchConfig` 结构体中添加新字段
2. 在 `internal/rules/engine.go` 的 `matchRule` 方法中添加匹配逻辑
3. 在 `internal/rules/engine.go` 的 `calculateScore` 方法中调整评分（可选）
4. 更新 WebUI 前端的规则表单（`internal/webui/frontend/src/views/Rules.vue`）
5. 更新 API 类型定义（`internal/webui/frontend/src/api/types.ts`）
6. 更新文档

### 添加新的通知渠道

1. 在 `internal/notifier/notifier.go` 中添加新的通知方法
2. 在 `internal/config/types.go` 的 `NotifierConfig` 中添加配置字段
3. 在 `internal/notifier/notifier.go` 的 `Notify` 方法中调用新渠道
4. 更新 WebUI 设置页面和文档

## 获取帮助

- 提交 Issue：[GitHub Issues](https://github.com/Maicarons/airlinesim-autobuy/issues)
- 发起讨论：[GitHub Discussions](https://github.com/Maicarons/airlinesim-autobuy/discussions)

感谢你的贡献！
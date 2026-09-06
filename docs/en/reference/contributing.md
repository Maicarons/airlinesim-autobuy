# Contributing

Thank you for your interest in contributing to AirlineSim Autobuy. This guide covers the development workflow, code style, testing, and pull request process.

## Development Setup

### Prerequisites

- Go 1.22+
- Node.js 18+
- Git

### Clone the Repository

```bash
git clone https://github.com/Maicarons/airlinesim-autobuy.git
cd airlinesim-autobuy
```

### Install Development Tools

```bash
make tools
```

This installs:

- [Air](https://github.com/air-verse/air) -- Hot reload for Go development
- [golangci-lint](https://golangci-lint.run/) -- Go linter

### Install Frontend Dependencies

```bash
make frontend-install
```

### Run in Development Mode

Start the backend with hot reload:

```bash
make dev
```

In a separate terminal, start the Vue frontend dev server:

```bash
cd internal/webui/frontend
npm run dev
```

The frontend dev server runs on `http://localhost:5173` and proxies API requests to the backend at `http://localhost:9090`.

## Project Structure

```
airlinesim-autobuy/
├── cmd/autobuy/           # Entry point
├── internal/
│   ├── auth/              # Authentication and session management
│   ├── client/            # Rate-limited HTTP client
│   ├── collector/         # Market page fetcher
│   ├── config/            # Configuration management
│   ├── engine/            # Pipeline orchestration
│   ├── executor/          # Purchase execution
│   ├── marketdata/        # Aircraft reference data
│   ├── notifier/          # Notification system
│   ├── parser/            # HTML market page parser
│   ├── rules/             # Rules engine
│   └── webui/             # Web UI server (Vue SPA + REST API)
│       └── frontend/      # Vue 3 + Vite frontend source
├── configs/               # YAML configuration files
├── docs/                  # VitePress documentation
├── scripts/               # Utility scripts
├── Makefile               # Build targets
├── go.mod                 # Go module definition
└── go.sum                 # Go module checksums
```

## Code Style

### Go

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) guidelines
- Run `golangci-lint run` before committing
- Use `gofmt` to format your code

```bash
# Format and lint
gofmt -s -w .
golangci-lint run
```

### Naming Conventions

| Element | Convention | Example |
|---|---|---|
| Packages | Lowercase, single word | `config`, `engine` |
| Types | PascalCase | `ServerConfig`, `MatchResult` |
| Functions | PascalCase (exported), camelCase (unexported) | `Evaluate()`, `matchRule()` |
| Variables | camelCase | `cfgStore`, `serverRunner` |
| Constants | PascalCase | `EventAircraftFound` |
| Files | Same as package name | `rules/engine.go`, `config/types.go` |

### Comments

- All exported types, functions, and constants must have doc comments
- Use complete sentences for doc comments: `// Package config manages...`
- Inline comments explain *why*, not *what*

### Error Handling

- Return errors rather than panicking (except for unrecoverable setup failures)
- Use `fmt.Errorf` with `%w` for error wrapping
- Log errors at the appropriate level using `slog`

```go
if err != nil {
    return fmt.Errorf("failed to process request: %w", err)
}
```

### Logging

- Use `slog` for all logging
- Use structured logging with key-value pairs
- Log levels: `Debug` for development details, `Info` for normal operations, `Warn` for issues, `Error` for failures

```go
slog.Info("engine started", "servers", count)
slog.Debug("processing aircraft", "type", ac.Type, "price", ac.Price)
slog.Warn("request failed", "error", err, "attempt", attempt)
slog.Error("failed to connect", "host", host, "error", err)
```

### Vue 3 / TypeScript

- Follow the [Vue 3 Style Guide](https://vuejs.org/style-guide/)
- Use TypeScript for type safety
- Use Composition API with `<script setup>` syntax
- Use Pinia for state management
- Use Vue Router for client-side routing

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test ./... -cover

# Run tests for a specific package
go test ./internal/rules/... -v

# Run tests with race detection
go test ./... -race
```

### Writing Tests

- Write unit tests for all new functionality
- Use table-driven tests for functions with multiple input/output combinations
- Mock external HTTP calls in tests

```go
func TestMatchRule(t *testing.T) {
    tests := []struct {
        name     string
        aircraft *AircraftOffer
        rule     *RuleConfig
        want     *MatchResult
    }{
        {
            name: "matches by type",
            aircraft: &AircraftOffer{
                Type:  "Airbus A320-200 heavy",
                Price: 2500000,
            },
            rule: &RuleConfig{
                Match: MatchConfig{
                    Types: []string{"Airbus A320-200 heavy"},
                    PriceRange: PriceRange{Max: 5000000},
                },
                Action: ActionConfig{AutoBuy: true},
            },
            want: &MatchResult{
                ShouldBuy: true,
            },
        },
    }

    e := New(nil)
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := e.matchRule(tt.aircraft, tt.rule)
            if got.ShouldBuy != tt.want.ShouldBuy {
                t.Errorf("matchRule().ShouldBuy = %v, want %v", got.ShouldBuy, tt.want.ShouldBuy)
            }
        })
    }
}
```

### Test Coverage

- Aim for at least 70% code coverage for new code
- Critical packages (`rules`, `config`, `parser`) should have 80%+ coverage

## Pull Request Process

### Before Submitting

1. Ensure your code builds: `make build`
2. Run all tests: `make test`
3. Run the linter: `make lint`
4. Update documentation if adding or changing features
5. Add or update tests for new functionality

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]
```

Types:

| Type | Description |
|---|---|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation changes |
| `style` | Code style changes (formatting, etc.) |
| `refactor` | Code refactoring |
| `test` | Adding or updating tests |
| `chore` | Build process or tooling changes |

Examples:

```
feat(rules): add financing filter option
fix(parser): handle European number format in prices
docs(api): add aircraft-data endpoint documentation
test(rules): add test cases for scoring calculation
```

### Pull Request Checklist

- [ ] Code builds without errors
- [ ] All tests pass
- [ ] Linter passes
- [ ] Documentation updated (if applicable)
- [ ] Tests added (if applicable)
- [ ] Commit messages follow conventional commits
- [ ] Branch is up to date with main

### Review Process

1. Maintainers will review your PR within a few days
2. Address any feedback or requested changes
3. Once approved, a maintainer will merge your PR

## Feature Requests

Open a [GitHub issue](https://github.com/Maicarons/airlinesim-autobuy/issues/new) with the `enhancement` label describing:

- The problem or use case
- The proposed solution
- Any alternatives considered

## Bug Reports

Open a [GitHub issue](https://github.com/Maicarons/airlinesim-autobuy/issues/new) with the `bug` label including:

- A clear description of the bug
- Steps to reproduce
- Expected behavior
- Actual behavior
- Logs or error messages
- Environment (OS, Go version, build date)

## Documentation

Documentation is built with [VitePress](https://vitepress.dev/).

### Running the Documentation Server

```bash
make docs-dev
```

This starts a VitePress dev server at `http://localhost:5173` with hot reload.

### Building Documentation

```bash
make docs-build
```

The output is in `docs/.vitepress/dist/`.

### Adding a New Page

1. Create the markdown file in the appropriate directory under `docs/en/`
2. Update `docs/.vitepress/config.ts` to add the page to the sidebar navigation
3. Update the corresponding locale configs for `zh-CN` and `ko` if applicable

## License

By contributing, you agree that your contributions will be licensed under the project's MIT License.

## See Also

- [Architecture](/reference/architecture) for understanding the codebase
- [API Reference](/reference/api) for REST API endpoints
- [Configuration Reference](/reference/config) for configuration fields
.PHONY: build clean run dev test docs-build docs-dev

APP_NAME = autobuy
GO_BUILD = go build -ldflags="-s -w" -o $(APP_NAME)

# Build the Go binary
build:
	$(GO_BUILD) ./cmd/$(APP_NAME)

# Build for current platform
build-all:
	GOOS=windows GOARCH=amd64 $(GO_BUILD)-windows-amd64.exe ./cmd/$(APP_NAME)
	GOOS=linux GOARCH=amd64 $(GO_BUILD)-linux-amd64 ./cmd/$(APP_NAME)
	GOOS=darwin GOARCH=amd64 $(GO_BUILD)-darwin-amd64 ./cmd/$(APP_NAME)

# Run the application
run:
	go run ./cmd/$(APP_NAME)

# Run in development mode (with hot reload if air is installed)
dev:
	air -- --config configs/config.yaml

# Run tests
test:
	go test ./... -v

# Run linter
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -f $(APP_NAME)
	rm -f $(APP_NAME)-windows-*
	rm -f $(APP_NAME)-linux-*
	rm -f $(APP_NAME)-darwin-*
	rm -rf internal/webui/frontend/dist

# Frontend dependencies
frontend-install:
	cd internal/webui/frontend && npm install

# Build frontend
frontend-build:
	cd internal/webui/frontend && npm run build

# Build everything (frontend + backend)
all: frontend-install frontend-build build

# VitePress documentation
docs-dev:
	cd docs && npm run dev

docs-build:
	cd docs && npm run build

# Install development tools
tools:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the Go binary"
	@echo "  run            - Run the application"
	@echo "  dev            - Run with hot reload"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  all            - Full build (frontend + backend)"
	@echo "  frontend-install - Install frontend dependencies"
	@echo "  frontend-build - Build Vue frontend"
	@echo "  docs-dev       - Start VitePress dev server"
	@echo "  docs-build     - Build VitePress documentation"
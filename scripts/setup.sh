#!/bin/bash
# First-time setup script for AirlineSim Autobuy
# Usage: ./scripts/setup.sh

set -euo pipefail

echo "╔══════════════════════════════════════════════╗"
echo "║   AirlineSim Autobuy - Setup                ║"
echo "╚══════════════════════════════════════════════╝"
echo ""

# Check Go
if ! command -v go &>/dev/null; then
    echo "✗ Go is not installed."
    echo "  Please install Go 1.22+ from https://go.dev"
    exit 1
fi
echo "✓ Go $(go version | grep -oP 'go\d+\.\d+')"

# Check Node
if command -v node &>/dev/null; then
    echo "✓ Node.js $(node --version)"
else
    echo "⚠ Node.js not found - frontend build will be skipped"
fi

echo ""

# Install frontend
if command -v node &>/dev/null; then
    echo "→ Installing frontend dependencies..."
    cd internal/webui/frontend
    npm install --silent
    echo "→ Building frontend..."
    npm run build
    cd ../..
    echo "✓ Frontend built"
fi

# Build backend
echo "→ Building Go binary..."
go build -ldflags="-s -w" -o autobuy ./cmd/autobuy/
echo "✓ Binary built: ./autobuy"

# Config
if [ ! -f configs/config.yaml ]; then
    echo "→ Creating config from template..."
    cp configs/config.example.yaml configs/config.yaml
    echo "⚠ Please edit configs/config.yaml with your AirlineSim credentials"
else
    echo "✓ Config file exists: configs/config.yaml"
fi

echo ""
echo "╔══════════════════════════════════════════════╗"
echo "║   Setup complete!                           ║"
echo "║                                            ║"
echo "║   Next steps:                              ║"
echo "║   1. Edit configs/config.yaml              ║"
echo "║   2. Run: ./autobuy                        ║"
echo "║   3. Open http://localhost:9090             ║"
echo "╚══════════════════════════════════════════════╝"
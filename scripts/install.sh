#!/bin/bash
# Install script for AirlineSim Autobuy
# Usage: ./scripts/install.sh

set -euo pipefail

echo "============================================"
echo " AirlineSim Autobuy - Install"
echo "============================================"
echo ""

# Check prerequisites
echo "Checking prerequisites..."

check_cmd() {
    if ! command -v "$1" &>/dev/null; then
        echo "  ✗ $1 not found"
        return 1
    fi
    echo "  ✓ $1 found"
    return 0
}

check_cmd go
HAS_NODE=$(check_cmd node && check_cmd npm && echo "yes" || echo "no")

echo ""

# Build frontend
if [ "$HAS_NODE" = "yes" ]; then
    echo "Installing frontend dependencies..."
    cd internal/webui/frontend
    npm install
    echo "Building frontend..."
    npm run build
    cd ../..
    echo "  ✓ Frontend built"
else
    echo "  ⚠ Node.js not found - skipping frontend build"
    echo "  The frontend won't be embedded. Use pre-built binary or install Node.js."
fi

echo ""

# Build backend
echo "Building Go backend..."
go build -ldflags="-s -w" -o autobuy ./cmd/autobuy/
echo "  ✓ Backend built: ./autobuy"

echo ""

# Setup config
if [ ! -f configs/config.yaml ]; then
    echo "Creating default config..."
    cp configs/config.example.yaml configs/config.yaml
    echo "  ✓ Created configs/config.yaml"
    echo "  ⚠ Edit configs/config.yaml with your credentials before running"
else
    echo "  ✓ configs/config.yaml already exists"
fi

echo ""
echo "============================================"
echo " Installation complete!"
echo ""
echo " Next steps:"
echo "   1. Edit configs/config.yaml with your credentials"
echo "   2. Run: ./autobuy"
echo "   3. Open http://localhost:9090"
echo "============================================"
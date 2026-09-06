#!/bin/bash
# Development script with hot-reload for AirlineSim Autobuy
# Requires: air (https://github.com/air-verse/air)
# Usage: ./scripts/dev.sh

set -euo pipefail

if ! command -v air &>/dev/null; then
    echo "Installing air (hot-reload tool)..."
    go install github.com/air-verse/air@latest
fi

echo "Starting development server with hot-reload..."
echo "Web UI will be available at http://localhost:9090"
echo ""

air -- --config configs/config.yaml
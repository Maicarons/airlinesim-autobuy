#!/bin/bash
# Run script for AirlineSim Autobuy
# Usage: ./scripts/run.sh [config_path]

set -euo pipefail

CONFIG="${1:-configs/config.yaml}"

if [ ! -f "$CONFIG" ]; then
    echo "Error: Config file not found: $CONFIG"
    echo "Usage: $0 [config_path]"
    echo ""
    echo "If you haven't set up the config yet, run:"
    echo "  cp configs/config.example.yaml configs/config.yaml"
    echo "  # Then edit configs/config.yaml with your credentials"
    exit 1
fi

echo "Starting AirlineSim Autobuy with config: $CONFIG"
echo "Web UI will be available at http://localhost:9090"
echo ""

if [ ! -f ./autobuy ]; then
    echo "Binary not found. Building first..."
    make build
    echo ""
fi

exec ./autobuy
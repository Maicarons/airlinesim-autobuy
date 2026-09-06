#!/bin/bash
# Docker build helper for AirlineSim Autobuy
# Usage: ./scripts/docker-build.sh [tag]

set -euo pipefail

TAG="${1:-latest}"
IMAGE="airlinesim-autobuy:${TAG}"

echo "Building Docker image: ${IMAGE}"
echo ""

# Build frontend first
echo "Building frontend..."
cd internal/webui/frontend
npm install
npm run build
cd ../..

echo ""
echo "Building Docker image..."
docker build -t "${IMAGE}" .

echo ""
echo "Docker image built: ${IMAGE}"
echo ""
echo "Run with:"
echo "  docker run -d -p 9090:9090 -v \$(pwd)/configs:/app/configs ${IMAGE}"
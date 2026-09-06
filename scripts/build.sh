#!/bin/bash
# Build script for AirlineSim Autobuy
# Usage: ./scripts/build.sh [linux|windows|darwin|all]

set -euo pipefail

APP_NAME="autobuy"
BUILD_DIR="build"
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

mkdir -p "${BUILD_DIR}"

build_platform() {
    local GOOS="$1"
    local GOARCH="$2"
    local EXT=""
    [ "${GOOS}" = "windows" ] && EXT=".exe"
    local OUTPUT="${BUILD_DIR}/${APP_NAME}-${GOOS}-${GOARCH}${EXT}"

    echo "Building for ${GOOS}/${GOARCH}..."
    GOOS="${GOOS}" GOARCH="${GOARCH}" CGO_ENABLED=0 \
        go build -ldflags="${LDFLAGS}" -o "${OUTPUT}" ./cmd/autobuy/
    echo "  → ${OUTPUT}"
}

case "${1:-}" in
    linux)
        build_platform linux amd64
        build_platform linux arm64
        ;;
    windows)
        build_platform windows amd64
        build_platform windows arm64
        ;;
    darwin|macos)
        build_platform darwin amd64
        build_platform darwin arm64
        ;;
    all|"")
        build_platform linux amd64
        build_platform linux arm64
        build_platform windows amd64
        build_platform darwin amd64
        build_platform darwin arm64
        ;;
    *)
        echo "Usage: $0 [linux|windows|darwin|all]"
        exit 1
        ;;
esac

echo ""
echo "Build complete! Binaries in ${BUILD_DIR}/"
ls -lh "${BUILD_DIR}/"
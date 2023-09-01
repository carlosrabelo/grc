#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BINARY_NAME="grc"
BUILD_DIR="$ROOT_DIR/bin"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS="-s -w -X main.version=$VERSION -X main.buildTime=$BUILD_TIME"

mkdir -p "$BUILD_DIR"
echo "Building $BINARY_NAME..."
CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o "$BUILD_DIR/$BINARY_NAME" ./grc/cmd/grc
echo "Done: $BUILD_DIR/$BINARY_NAME"

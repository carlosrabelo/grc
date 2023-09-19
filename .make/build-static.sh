#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BINARY_NAME="grc"
BUILD_DIR="$ROOT_DIR/bin"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS="-s -w -X github.com/carlosrabelo/grc/grc/internal/app.Version=$VERSION -X github.com/carlosrabelo/grc/grc/internal/app.BuildTime=$BUILD_TIME"

mkdir -p "$BUILD_DIR"
echo "Building static $BINARY_NAME..."
CGO_ENABLED=0 go build \
  -tags netgo \
  -ldflags="$LDFLAGS" \
  -o "$BUILD_DIR/$BINARY_NAME" \
  ./grc/cmd/grc
echo "Done: $BUILD_DIR/$BINARY_NAME"

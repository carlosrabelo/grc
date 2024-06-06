#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

cd "$ROOT_DIR"
echo "=== go fmt ==="
go fmt ./...
echo ""
echo "=== go vet ==="
go vet ./...
echo ""
echo "=== go test ==="
go test ./... -count=1

#!/usr/bin/env bash
set -euo pipefail

echo "=== go fmt ==="
go fmt ./...
echo ""
echo "=== go vet ==="
go vet ./...
echo ""
echo "=== golangci-lint ==="
golangci-lint run ./... 2>/dev/null || echo "(golangci-lint not installed)"
echo ""
echo "=== go test ==="
go test ./... -count=1

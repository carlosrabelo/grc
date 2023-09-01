#!/usr/bin/env bash
set -euo pipefail

echo "=== go fmt ==="
go fmt ./...
echo ""
echo "=== go vet ==="
go vet ./...
echo ""
echo "=== go test ==="
go test ./... -count=1

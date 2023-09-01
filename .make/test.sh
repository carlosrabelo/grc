#!/usr/bin/env bash
set -euo pipefail

echo "Running tests..."
go test ./... -v -count=1

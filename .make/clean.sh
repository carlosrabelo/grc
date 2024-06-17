#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

cd "$ROOT_DIR"
go clean
# Remove built binaries but keep bin/.gitkeep so the empty tree stays tracked.
find "$ROOT_DIR/bin" -mindepth 1 ! -name '.gitkeep' -delete 2>/dev/null || true
rm -rf "$ROOT_DIR/dist"
echo "Cleaned build artifacts"

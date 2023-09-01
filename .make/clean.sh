#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
echo "Cleaning..."
rm -rf "$ROOT_DIR/bin"
echo "Done"

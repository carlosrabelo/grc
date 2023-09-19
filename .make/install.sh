#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BINARY_NAME="grc"
SOURCE="$ROOT_DIR/bin/$BINARY_NAME"

if [ ! -f "$SOURCE" ]; then
  echo "Binary not found. Running build first..."
  "$ROOT_DIR/.make/build.sh"
fi

if [ "$(id -u)" -eq 0 ]; then
  DEST="/usr/local/bin/$BINARY_NAME"
else
  DEST="${HOME}/.local/bin/$BINARY_NAME"
fi

mkdir -p "$(dirname "$DEST")"
cp "$SOURCE" "$DEST"
chmod 755 "$DEST"
echo "Installed: $DEST"

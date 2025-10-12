#!/bin/bash
set -euo pipefail

# This script make install of grc binary with friendly steps.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

BIN=${1:-grc}
BUILD_DIR=${2:-"${PROJECT_ROOT}/bin"}
BIN_PATH="${BUILD_DIR}/${BIN}"

if [ ! -f "${BIN_PATH}" ]; then
    echo "Binary ${BIN_PATH} not found. Please build before install."
    exit 1
fi

echo "Installing ${BIN}..."
if [ "$(id -u)" = "0" ]; then
    echo "Installing to /usr/local/bin (system-wide)"
    install -m 755 "${BIN_PATH}" /usr/local/bin/
else
    echo "Installing to ${HOME}/.local/bin (user-local)"
    mkdir -p "${HOME}/.local/bin"
    install -m 755 "${BIN_PATH}" "${HOME}/.local/bin/"
    echo "Make sure ${HOME}/.local/bin is in your PATH"
fi

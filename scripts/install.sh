#!/bin/bash
set -euo pipefail

# This script make install of grc binary with friendly steps.
# 
# Usage:
#   ./scripts/install.sh [binary_name] [build_dir]
#   GRC_INSTALL_DIR=/custom/path ./scripts/install.sh
#   GRC_INSTALL_DIR=/custom/path ./scripts/install.sh my-binary

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

BIN=${1:-grc}
BUILD_DIR=${2:-"${PROJECT_ROOT}/bin"}
BIN_PATH="${BUILD_DIR}/${BIN}"

# Custom installation directory via environment variable
INSTALL_DIR="${GRC_INSTALL_DIR:-}"

if [ ! -f "${BIN_PATH}" ]; then
    echo "Binary ${BIN_PATH} not found. Please build before install."
    exit 1
fi

echo "Installing ${BIN}..."

if [ -n "${INSTALL_DIR}" ]; then
    # Custom installation directory
    echo "Installing to ${INSTALL_DIR} (custom directory)"
    mkdir -p "${INSTALL_DIR}"
    install -m 755 "${BIN_PATH}" "${INSTALL_DIR}/"
    echo "Make sure ${INSTALL_DIR} is in your PATH"
elif [ "$(id -u)" = "0" ]; then
    # System-wide installation
    echo "Installing to /usr/local/bin (system-wide)"
    install -m 755 "${BIN_PATH}" /usr/local/bin/
else
    # User-local installation
    echo "Installing to ${HOME}/.local/bin (user-local)"
    mkdir -p "${HOME}/.local/bin"
    install -m 755 "${BIN_PATH}" "${HOME}/.local/bin/"
    echo "Make sure ${HOME}/.local/bin is in your PATH"
fi

echo "Installation completed successfully!"

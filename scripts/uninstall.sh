#!/bin/bash
set -euo pipefail

# This script make uninstall of grc binary in calm manner.
#
# Usage:
#   ./scripts/uninstall.sh [binary_name]
#   GRC_INSTALL_DIR=/custom/path ./scripts/uninstall.sh
#   GRC_INSTALL_DIR=/custom/path ./scripts/uninstall.sh my-binary

BIN=${1:-grc}

# Custom installation directory via environment variable
INSTALL_DIR="${GRC_INSTALL_DIR:-}"

echo "Uninstalling ${BIN}..."

if [ -n "${INSTALL_DIR}" ]; then
    # Custom installation directory
    if [ -f "${INSTALL_DIR}/${BIN}" ]; then
        rm -f "${INSTALL_DIR}/${BIN}"
        echo "Removed from ${INSTALL_DIR}"
    else
        echo "Binary ${BIN} not found in ${INSTALL_DIR}"
    fi
elif [ "$(id -u)" = "0" ]; then
    # System-wide installation
    if [ -f "/usr/local/bin/${BIN}" ]; then
        rm -f "/usr/local/bin/${BIN}"
        echo "Removed from /usr/local/bin"
    else
        echo "Binary ${BIN} not found in /usr/local/bin"
    fi
else
    # User-local installation
    if [ -f "${HOME}/.local/bin/${BIN}" ]; then
        rm -f "${HOME}/.local/bin/${BIN}"
        echo "Removed from ${HOME}/.local/bin"
    else
        echo "Binary ${BIN} not found in ${HOME}/.local/bin"
    fi
fi

echo "Uninstallation completed!"

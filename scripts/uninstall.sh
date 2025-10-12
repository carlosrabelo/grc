#!/bin/bash
set -euo pipefail

# This script make uninstall of grc binary in calm manner.

BIN=${1:-grc}

echo "Uninstalling ${BIN}..."
if [ "$(id -u)" = "0" ]; then
    rm -f "/usr/local/bin/${BIN}"
    echo "Removed from /usr/local/bin"
else
    rm -f "${HOME}/.local/bin/${BIN}"
    echo "Removed from ${HOME}/.local/bin"
fi

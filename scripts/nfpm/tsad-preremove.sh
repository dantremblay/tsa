#!/bin/sh
set -e

# Stop and disable the daemon before the package is removed.
if command -v systemctl >/dev/null 2>&1; then
	systemctl stop tsad.service >/dev/null 2>&1 || true
	systemctl disable tsad.service >/dev/null 2>&1 || true
fi

exit 0

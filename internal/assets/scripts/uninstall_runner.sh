#!/usr/bin/env bash
set -euo pipefail

if [ -f ./svc.sh ]; then
  ./svc.sh stop || true
fi

if [ -f ./config.sh ]; then
  ./config.sh remove --unattended || true
fi

echo "Runner uninstalled"


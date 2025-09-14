#!/usr/bin/env bash
set -euo pipefail

echo "[fastauto] Deploying repo at $(pwd)"

if command -v git >/dev/null 2>&1; then
  echo "[fastauto] git pull --ff-only"
  git pull --ff-only
fi

if [ -f package.json ]; then
  if command -v npm >/dev/null 2>&1; then
    echo "[fastauto] npm ci && npm run build"
    npm ci && npm run build || true
  fi
fi

if [ -f go.mod ]; then
  if command -v go >/dev/null 2>&1; then
    echo "[fastauto] go build ./..."
    go build ./... || true
  fi
fi

echo "[fastauto] Deploy script finished"


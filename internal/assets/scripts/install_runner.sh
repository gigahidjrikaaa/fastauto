#!/usr/bin/env bash
set -euo pipefail

# A minimal helper to download and configure GitHub Actions runner
# Usage: GH_REPO_URL=https://github.com/owner/repo GH_TOKEN=XXXX ./install_runner.sh

ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH=linux-x64 ;;
  aarch64|arm64) ARCH=linux-arm64 ;;
  *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac

VER="2.317.0" # pinned; can be updated
URL="https://github.com/actions/runner/releases/download/v${VER}/actions-runner-${ARCH}-${VER}.tar.gz"

mkdir -p .
cd .

echo "Downloading $URL"
curl -fsSL -o runner.tgz "$URL"
tar xzf runner.tgz
rm -f runner.tgz

if [ -z "${GH_REPO_URL:-}" ] || [ -z "${GH_TOKEN:-}" ]; then
  echo "Set GH_REPO_URL and GH_TOKEN environment variables" >&2
  exit 1
fi

./config.sh --url "$GH_REPO_URL" --token "$GH_TOKEN" --unattended --replace || true
echo "Runner configured. You can now start it with ./run.sh"


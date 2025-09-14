#!/usr/bin/env bash
set -euo pipefail

REPO_OWNER="gigahidjrikaaa"
REPO_NAME="fastauto"
PROJECT="${REPO_NAME}"

err() { echo "[${PROJECT}] $*" >&2; }
info() { echo "[${PROJECT}] $*"; }

need() { command -v "$1" >/dev/null 2>&1 || { err "missing dependency: $1"; exit 1; }; }
need curl
need tar

# Determine OS/ARCH
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH_RAW="$(uname -m)"
case "$OS" in
  linux) OS=linux ;;
  *) err "unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH_RAW" in
  x86_64|amd64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) err "unsupported ARCH: $ARCH_RAW"; exit 1 ;;
esac

# Install dir
BIN_DIR="${FASTAUTO_BIN_DIR:-}"
if [[ -z "${BIN_DIR}" ]]; then
  if [[ "${EUID:-$(id -u)}" -eq 0 ]]; then
    BIN_DIR="/usr/local/bin"
  else
    BIN_DIR="$HOME/.local/bin"
  fi
fi
mkdir -p "$BIN_DIR"

TMP="$(mktemp -d)"; trap 'rm -rf "$TMP"' EXIT

UA="${PROJECT}-installer"
owner_repo="${REPO_OWNER}/${REPO_NAME}"
api_base="https://api.github.com/repos/${owner_repo}/releases"

# Desired tag (optional): set FASTAUTO_VERSION=vX.Y.Z
TAG="${FASTAUTO_VERSION:-}"
if [[ -n "$TAG" ]]; then
  TAG="v${TAG#v}"
else
  # Resolve latest tag via redirect without parsing JSON
  latest_url=$(curl -fsSLI -o /dev/null -w '%{url_effective}' "https://github.com/${owner_repo}/releases/latest")
  TAG=$(sed -E 's#.*/tag/(v[^/]+)$#\1#' <<<"$latest_url")
fi

if [[ -z "$TAG" ]]; then
  err "failed to determine release tag; set FASTAUTO_VERSION=vX.Y.Z and retry"
  exit 1
fi

ASSET_URL="https://github.com/${owner_repo}/releases/download/${TAG}/${PROJECT}_${TAG#v}_${OS}_${ARCH}.tar.gz"
info "downloading $ASSET_URL"
ARCHIVE="$TMP/${PROJECT}.tar.gz"
curl -fL --retry 3 --retry-delay 1 -o "$ARCHIVE" "$ASSET_URL"

info "extracting archive"
tar -xzf "$ARCHIVE" -C "$TMP"

# Locate binary in extracted tree
BIN_SRC="$(find "$TMP" -maxdepth 2 -type f -name "${PROJECT}" | head -n1 || true)"
if [[ -z "$BIN_SRC" ]]; then
  err "failed to locate ${PROJECT} binary in archive"
  exit 1
fi

DEST="$BIN_DIR/${PROJECT}"
BACKUP="${DEST}.bak.$(date +%Y%m%dT%H%M%S)"
if [[ -f "$DEST" ]]; then
  mv -f "$DEST" "$BACKUP" || true
  info "backed up existing to $BACKUP"
fi

install -m 0755 "$BIN_SRC" "$DEST"
info "installed ${PROJECT} to $DEST"

# PATH hint
case ":$PATH:" in
  *:"$BIN_DIR":*) ;;
  *) info "Note: $BIN_DIR is not in PATH. Add it to your shell profile." ;;
esac

"$DEST" version || true

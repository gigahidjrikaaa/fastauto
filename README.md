# fastauto

[![CI](https://github.com/gigahidjrikaaa/fastauto/actions/workflows/ci.yml/badge.svg)](https://github.com/gigahidjrikaaa/fastauto/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/gigahidjrikaaa/fastauto?display_name=tag&sort=semver)](https://github.com/gigahidjrikaaa/fastauto/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gigahidjrikaaa/fastauto)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/gigahidjrikaaa/fastauto)](https://goreportcard.com/report/github.com/gigahidjrikaaa/fastauto)
[![lint](https://img.shields.io/badge/lint-golangci--lint-blue)](https://golangci-lint.run)
[![Downloads](https://img.shields.io/github/downloads/gigahidjrikaaa/fastauto/total.svg)](https://github.com/gigahidjrikaaa/fastauto/releases)
[![Platforms](https://img.shields.io/badge/platform-linux%20amd64%20%7C%20arm64-2ea44f)](#)

Fastauto is a tiny, safe tool that keeps your Git repo up-to-date and runs your deploy script automatically. One command sets it up; after that, pushes to your repo deploy your code.

You can run fastauto in two simple modes:
- Webhook Mode: Runs a small HTTP server that verifies GitHub webhooks (HMAC) and runs `deploy.sh` on pushes.
- Runner Mode: Installs a self‑hosted GitHub Actions runner and adds a minimal workflow that runs `deploy.sh` on pushes.

## Highlights

- Single static CLI binary (Linux amd64/arm64)
- Beginner-friendly setup with safe defaults
- Commands: `init`, `install`, `status`, `logs`, `deploy --now`, `secret rotate`, `uninstall`, `version`
- Config lives in your repo (`.fastauto.yml`) and globally (`$XDG_CONFIG_HOME/fastauto/config.yml`)
- Journald logging, idempotent writes, backups, and atomic secret rotation

## Supported Platforms

- Linux x86_64 (amd64) and ARM64 (aarch64)

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Webhook Mode](#webhook-mode)
  - [HTTPS/TLS](#https-tls)
- [Runner Mode](#runner-mode)
- [Tokens & Credentials](#tokens--credentials)
- [Commands Reference](#commands-reference)
- [Configuration Files](#configuration-files)
- [Deploy Script](#deploy-script)
- [Troubleshooting](#troubleshooting)
- [Security Notes](#security-notes)
- [Uninstall](#uninstall)
- [License](#license)

## Installation

### Quick install

Option A — Go toolchain (recommended if you have Go):

```bash
GOFLAGS=-trimpath CGO_ENABLED=0 go install github.com/gigahidjrikaaa/fastauto/cmd/fastauto@latest
export PATH="$(go env GOPATH)/bin:$PATH"  # if not already
fastauto version
```

Option B — Build, then self-install to your PATH:

```bash
git clone https://github.com/gigahidjrikaaa/fastauto.git
cd fastauto
go build -o bin/fastauto ./cmd/fastauto
./bin/fastauto self install
# Re-open your shell if ~/.local/bin was added
fastauto version
```

Option C — Curl installer (no Go needed):

```bash
curl -fsSL https://raw.githubusercontent.com/gigahidjrikaaa/fastauto/main/scripts/install.sh | bash
# pin to a version (e.g., v0.1.0):
FASTAUTO_VERSION=v0.1.0 curl -fsSL https://raw.githubusercontent.com/gigahidjrikaaa/fastauto/main/scripts/install.sh | bash
# install system-wide (requires root):
curl -fsSL https://raw.githubusercontent.com/gigahidjrikaaa/fastauto/main/scripts/install.sh | sudo FASTAUTO_BIN_DIR=/usr/local/bin bash
```

### From source (needs Go 1.22+)

```bash
git clone https://github.com/gigahidjrikaaa/fastauto.git
cd fastauto
go build -o bin/fastauto ./cmd/fastauto
# or
GOFLAGS=-trimpath CGO_ENABLED=0 go install ./cmd/fastauto
```

### From releases (recommended)

1. Download the tarball for your platform from the Releases page
2. Extract and place `fastauto` somewhere in your `PATH` (e.g., `/usr/local/bin`)

## Quick Start

1. Go to your app's Git repo on the server:

   ```bash
   cd /path/to/your/repo
   ```

2. Initialize fastauto and choose a mode when prompted:

   ```bash
   fastauto init
   ```

   This writes `.fastauto.yml`, adds `deploy.sh` (if missing), and can install a systemd unit for you.

3. Install and start the service:

   ```bash
   fastauto install
   ```

4. Watch logs:

   ```bash
   fastauto logs -f
   ```

## Webhook Mode

- During init, choose "webhook" or run:

  ```bash
  fastauto init --mode webhook --port 8080
  ```

- Create a GitHub webhook pointing to `http://your-server:8080/hook`.
  - Secret lives in `$XDG_CONFIG_HOME/fastauto/config.yml` (auto-generated if missing).
- On each push (matching your configured branches), fastauto does a `git pull --ff-only` and runs `./deploy.sh`.

### HTTPS / TLS

Add to `.fastauto.yml`:

```yaml
webhook:
  address: ":8443"
  tls_cert_file: "/path/to/cert.pem"
  tls_key_file: "/path/to/key.pem"
```

Or set environment variables:

```bash
export FASTAUTO_WEBHOOK_TLS_CERT_FILE=/path/to/cert.pem
export FASTAUTO_WEBHOOK_TLS_KEY_FILE=/path/to/key.pem
```

## Runner Mode

- During init, choose "runner" or run:

  ```bash
  fastauto init --mode runner
  ```

- Install the service and write the workflow:

  ```bash
  fastauto install
  ```

- Configure the runner with an ephemeral token:

  ```bash
  cd .fastauto/runner
  GH_REPO_URL=https://github.com/OWNER/REPO \
  GH_TOKEN=$(gh api -X POST repos/OWNER/REPO/actions/runners/registration-token -q .token) \
  ./install_runner.sh
  ```

- The `fastauto-runner.service` supervises the runner; pushes will trigger the workflow, which runs `./deploy.sh`.

## Tokens & Credentials

### Webhook Secret (Webhook Mode)

- Where: stored in `$XDG_CONFIG_HOME/fastauto/config.yml` under `webhook_secret`.
- How it's created: auto-generated by `fastauto install` if missing.
- How to use in GitHub:
  1. Repo → Settings → Webhooks → Add webhook
  2. Payload URL: `http://your-server:8080/hook` (or your TLS URL)
  3. Content type: `application/json`
  4. Secret: paste the value from the global config file
  5. Events: "Just the push event"
- Rotate anytime: `fastauto secret rotate` (then update the webhook Secret in GitHub).
- Show current value locally: `fastauto secret show`

### Git Credentials for Autopull

The server must be able to `git pull` your repo:

- Recommended: SSH Deploy Key

  ```bash
  ssh-keygen -t ed25519 -C "fastauto@$(hostname)"
  # Add the public key to GitHub → Repo → Settings → Deploy keys (read-only)
  git remote set-url origin git@github.com:OWNER/REPO.git
  ```

- Alternative: HTTPS with a Personal Access Token (PAT)

  ```bash
  # Create a fine-scoped token with read access to the repo
  git remote set-url origin https://<TOKEN>@github.com/OWNER/REPO.git
  ```

### Runner Registration Token (Runner Mode)

- Where in GitHub UI:
  - Repo → Settings → Actions → Runners → "New self-hosted runner" → copy the registration token
  - For org runners: Organization → Settings → Actions → Runners → "New runner"
- CLI/API alternative (requires permissions):

  ```bash
  # Repo-scoped token
  gh api -X POST repos/OWNER/REPO/actions/runners/registration-token -q .token
  # Org-scoped token
  gh api -X POST orgs/ORG/actions/runners/registration-token -q .token
  ```

- Use with installer:

  ```bash
  cd .fastauto/runner
  GH_REPO_URL=https://github.com/OWNER/REPO \
  GH_TOKEN=$(gh api -X POST repos/OWNER/REPO/actions/runners/registration-token -q .token) \
  ./install_runner.sh
  ```

## Manual Deploy

Run your deploy script on demand:

```bash
fastauto deploy --now
```

## Commands Reference

- `init`: Detects repo, branch, writes `.fastauto.yml` and `deploy.sh` (if missing), and can install units.
- `install`: Installs and enables the systemd service for your chosen mode.
- `status`: Shows systemd service status.
- `logs [-f]`: Shows journald logs for the running service.
- `deploy --now`: Pulls latest and runs `./deploy.sh` once.
- `secret rotate`: Rotates the webhook secret in the global config (keep GitHub in sync).
- `secret show`: Prints the current webhook secret and config path.
- `uninstall`: Stops, disables, and removes the systemd unit.
- `version`: Prints version info.
- `self install [--bin-dir DIR]`: Install this binary into your PATH (defaults to `~/.local/bin` or `/usr/local/bin` as root).
- `self uninstall [--bin-dir DIR]`: Remove the installed binary.
- `completion [bash|zsh|fish|powershell]`: Generate shell completion script.

## Configuration Files

- Repo: `./.fastauto.yml` (mode, branches, webhook address and TLS, etc.)
- Global: `$XDG_CONFIG_HOME/fastauto/config.yml` (webhook secret)

Example `.fastauto.yml`:

```yaml
mode: webhook
repo_path: /path/to/your/repo
branches: ["main"]
webhook:
  address: ":8080"
  # tls_cert_file: "/path/to/cert.pem"
  # tls_key_file: "/path/to/key.pem"
# runner:
#   labels: ["self-hosted", "linux", "x64"]
```

## Deploy Script

A simple Bash script that you control. The default template does:

- `git pull --ff-only`
- If Node project: `npm ci && npm run build`
- If Go project: `go build ./...`

Put your real deploy steps here (restart services, copy files, etc.).

## Troubleshooting

- Check service state: `fastauto status`
- Follow logs: `fastauto logs -f`
- See TROUBLESHOOTING.md for common fixes.

## Security Notes

- Every webhook is HMAC-verified using the secret in global config.
- Secrets are rotated atomically and existing files are backed up with timestamps.

## Uninstall

`fastauto uninstall` removes the systemd unit. Delete `.fastauto.yml` and scripts if you no longer need them.

## License

MIT — see `LICENSE`

fastauto
=========

<div align="left">

[![CI](https://github.com/gigahidjrikaaa/fastauto/actions/workflows/ci.yml/badge.svg)](https://github.com/gigahidjrikaaa/fastauto/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/gigahidjrikaaa/fastauto?display_name=tag&sort=semver)](https://github.com/gigahidjrikaaa/fastauto/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gigahidjrikaaa/fastauto)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/gigahidjrikaaa/fastauto)](https://goreportcard.com/report/github.com/gigahidjrikaaa/fastauto)
[![lint](https://img.shields.io/badge/lint-golangci--lint-blue)](https://golangci-lint.run)
[![Downloads](https://img.shields.io/github/downloads/gigahidjrikaaa/fastauto/total.svg)](https://github.com/gigahidjrikaaa/fastauto/releases)
[![Platforms](https://img.shields.io/badge/platform-linux%20amd64%20%7C%20arm64-2ea44f)](#)

</div>

Fastauto is a tiny, safe tool that keeps your Git repo up‑to‑date and runs your deploy script automatically. One command sets it up; after that, pushes to your repo deploy your code.

You can run fastauto in two simple modes:
- Webhook Mode: Runs a small HTTP server that verifies GitHub webhooks (HMAC) and runs `deploy.sh` on pushes.
- Runner Mode: Installs a self‑hosted GitHub Actions runner and adds a minimal workflow that runs `deploy.sh` on pushes.

Highlights
- Single static CLI binary (Linux amd64/arm64)
- Beginner‑friendly setup with safe defaults
- Commands: `init`, `install`, `status`, `logs`, `deploy --now`, `secret rotate`, `uninstall`, `version`
- Config lives in your repo (`.fastauto.yml`) and globally (`$XDG_CONFIG_HOME/fastauto/config.yml`)
- Journald logging, idempotent writes, backups, and atomic secret rotation

Supported Platforms
- Linux x86_64 (amd64) and ARM64 (aarch64)

Install
- From source (needs Go 1.22+):
  - `git clone <this repo>`
  - `cd fastauto`
  - `go build -o bin/fastauto ./cmd/fastauto`
  - Optionally: `go install ./cmd/fastauto`

- From releases (recommended):
  - Download the tarball for your platform from Releases
  - Extract and place `fastauto` somewhere in your `PATH` (e.g., `/usr/local/bin`)

Quick Start (Most Users)
1) Go to your app’s Git repo on the server: `cd /path/to/your/repo`
2) Initialize fastauto and choose a mode when prompted:
   - `fastauto init`
   - This writes `.fastauto.yml`, adds `deploy.sh` (if missing), and can install a systemd unit for you.
3) Install and start the service: `fastauto install`
4) Watch logs: `fastauto logs -f`

Using Webhook Mode
- During `fastauto init`, choose “webhook” (or run: `fastauto init --mode webhook --port 8080`).
- Create a GitHub webhook pointing to `http://your-server:8080/hook`.
  - Secret lives in `$XDG_CONFIG_HOME/fastauto/config.yml` (auto‑generated if missing).
- On each push (matching your configured branches), fastauto does a `git pull --ff-only` and runs `./deploy.sh`.

Optional HTTPS for Webhooks
- Add to `.fastauto.yml`:
  - `webhook.address: ":8443"`
  - `webhook.tls_cert_file: "/path/to/cert.pem"`
  - `webhook.tls_key_file: "/path/to/key.pem"`
- Or set env vars: `FASTAUTO_WEBHOOK_TLS_CERT_FILE` and `FASTAUTO_WEBHOOK_TLS_KEY_FILE`.

Using Runner Mode
- During `fastauto init`, choose “runner” (or run: `fastauto init --mode runner`).
- `fastauto install` will add a systemd unit for the runner and write `.github/workflows/fastauto.yml`.
- Configure the runner with an ephemeral token:
  - `cd .fastauto/runner`
  - `GH_REPO_URL=https://github.com/OWNER/REPO GH_TOKEN=YOUR_SHORT_LIVED_TOKEN ./install_runner.sh`
- The `fastauto-runner.service` supervises the runner; pushes will trigger the workflow, which runs `./deploy.sh`.

Manual Deploy
- Run your deploy script on demand: `fastauto deploy --now`

Commands Reference
- `init`: Detects repo, branch, writes `.fastauto.yml` and `deploy.sh` (if missing), and can install units.
- `install`: Installs and enables the systemd service for your chosen mode.
- `status`: Shows systemd service status.
- `logs [-f]`: Shows journald logs for the running service.
- `deploy --now`: Pulls latest and runs `./deploy.sh` once.
- `secret rotate`: Rotates the webhook secret in the global config (keep GitHub in sync).
- `uninstall`: Stops, disables, and removes the systemd unit.
- `version`: Prints version info.

Configuration Files
- Repo: `./.fastauto.yml` (mode, branches, webhook address and TLS, etc.)
- Global: `$XDG_CONFIG_HOME/fastauto/config.yml` (webhook secret)

Deploy Script (`deploy.sh`)
- A simple Bash script that you control. The default template does:
  - `git pull --ff-only`
  - If Node project: `npm ci && npm run build`
  - If Go project: `go build ./...`
- Put your real deploy steps here (restart services, copy files, etc.).

Troubleshooting
- Check service state: `fastauto status`
- Follow logs: `fastauto logs -f`
- See TROUBLESHOOTING.md for common fixes.

Security Notes
- Every webhook is HMAC‑verified using the secret in global config.
- Secrets are rotated atomically and existing files are backed up with timestamps.

Uninstall
- `fastauto uninstall` removes the systemd unit. Delete `.fastauto.yml` and scripts if you no longer need them.

License
- MIT — see LICENSE

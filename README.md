fastauto
=========

Fast, safe “autopull + auto-deploy” for any Git repo with one command. Runs in two modes:

- Webhook mode: HMAC-verified HTTP server triggers `deploy.sh` on push.
- Runner mode: Self-hosted GitHub Actions runner executes a minimal workflow.

Features
- Single static CLI binary (Linux amd64/arm64)
- Commands: `init`, `install`, `status`, `logs`, `deploy --now`, `secret rotate`, `uninstall`
- Config persists in repo `.fastauto.yml` and global `$XDG_CONFIG_HOME/fastauto/config.yml`
- Safe defaults: idempotent writes, backups, atomic secret updates, journald logging
- Assets: systemd unit templates, example `deploy.sh`, runner install scripts
- Goreleaser, CI workflow, golangci-lint

See QUICKSTART.md to get going in minutes.


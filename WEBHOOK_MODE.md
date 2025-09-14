Webhook Mode
============

What it does
- Starts an HTTP server with HMAC verification for GitHub webhooks.
- On push events matching configured branches, runs `deploy.sh` after a `git pull`.

Install steps
1. `fastauto init --mode webhook --port 8080`
2. `fastauto install`
3. Configure a GitHub webhook to `http://your-server:8080/hook` using the secret from `$XDG_CONFIG_HOME/fastauto/config.yml`.

Security
- Uses `X-Hub-Signature-256` (or `sha1`) to verify payloads. See SECURITY.md for details.


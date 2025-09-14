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

TLS (optional)
- Provide certificate and key paths in `.fastauto.yml` under `webhook:`

  ```yaml
  webhook:
    address: ":8443"
    tls_cert_file: "/etc/ssl/certs/your.crt"
    tls_key_file: "/etc/ssl/private/your.key"
  ```

- Or set env vars `FASTAUTO_WEBHOOK_TLS_CERT_FILE` and `FASTAUTO_WEBHOOK_TLS_KEY_FILE`.
- Systemd unit will use the repo config when starting the internal webhook server.

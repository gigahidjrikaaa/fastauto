Security Overview
=================

- HMAC verification: `X-Hub-Signature-256` and `X-Hub-Signature` are validated using a secret stored in `$XDG_CONFIG_HOME/fastauto/config.yml`.
- Secret rotation: `fastauto secret rotate` updates the secret atomically and keeps a timestamped backup.
- Least privileges: services run as the invoking user with systemd user units by default.
- Journald logging: service output goes to journald for traceability.

Operational Tips
- Protect the host firewall to expose the webhook port only as needed.
- For runner mode, prefer repository-level runner with limited scope.


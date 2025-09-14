Troubleshooting
===============

Common issues

- Service not starting
  - Check `fastauto status` and `fastauto logs -f`
  - Try `systemctl --user daemon-reload`

- Webhook 401 invalid signature
  - Ensure the secret in GitHub matches `$XDG_CONFIG_HOME/fastauto/config.yml`
  - Rotate the secret: `fastauto secret rotate` and update GitHub

- Runner not picking jobs
  - Confirm `GH_REPO_URL` and `GH_TOKEN` used during config
  - Check `.fastauto/runner/_diag` logs


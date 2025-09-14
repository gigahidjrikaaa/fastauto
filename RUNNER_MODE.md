Runner Mode
===========

What it does
- Installs a self-hosted GitHub Actions runner under `.fastauto/runner`.
- Adds `.github/workflows/fastauto.yml` to run `deploy.sh` on push.
- Installs a systemd user unit `fastauto-runner.service` to supervise the runner.

Install steps
1. `fastauto init --mode runner`
2. `fastauto install`
3. Configure runner: set `GH_REPO_URL` and `GH_TOKEN` then run `.fastauto/runner/install_runner.sh`.

Notes
- `GH_TOKEN` is a short-lived registration token obtained via the GitHub API or web UI.
- The unit restarts the runner if it exits.


Quickstart
==========

1) Install fastauto

- From source: `go install ./cmd/fastauto`
- Or use Goreleaser artifacts (see Releases)

2) Initialize in your repo

- `cd your-repo`
- `fastauto init` and choose mode (webhook or runner)

3) Install services

- `fastauto install`
- Check: `fastauto status` and `fastauto logs -f`

Webhook mode
- Expose the listen port (default :8080) and create a GitHub webhook pointing to `/hook` with the secret from `$XDG_CONFIG_HOME/fastauto/config.yml`.

Runner mode
- Provide `GH_REPO_URL` and `GH_TOKEN` then run the generated installer under `.fastauto/runner/install_runner.sh`. The systemd unit runs the `run.sh` script.

Manual deploy
- `fastauto deploy --now`


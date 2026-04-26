# M1 Closure Report

Date: 2026-04-25
Branch: m1-closure-v0100

## Summary
Implemented and validated the `kiln dev` M1 behavior on a branch forked from the latest version (`v0.10.0`): `--no-serve` watch-only mode and `--on-rebuild` webhook behavior are present and validated with reproducible smoke scripts under `demo/scripts`.

## Implementation Status
- `--no-serve` flag: implemented in `internal/cli` and validated.
- `--on-rebuild <url>` flag: implemented in `internal/cli` and validated.
- Webhook timeout: 5 seconds (`http.Client{Timeout: 5 * time.Second}`).
- Webhook errors are non-fatal to rebuild loop.
- CLI docs updated in `docs/Commands/dev.md`.

## Validation Results
Command: `devenv shell -- go test ./...`
- Result: fail (exit 1)
- Evidence: `demo/logs/go_test_m1.log`
- Failing tests:
  - `internal/builder TestGeneratePageOGImages`
  - `internal/templates TestHead_OGMetaTags`

Command: `devenv shell -- ./demo/scripts/smoke_no_serve.sh`
- Result: pass (exit 0)
- Evidence:
  - `demo/logs/smoke_no_serve.result`
  - `demo/logs/smoke_no_serve.log`

Command: `devenv shell -- ./demo/scripts/smoke_on_rebuild.sh`
- Result: pass (exit 0)
- Evidence:
  - `demo/logs/smoke_on_rebuild.result`
  - `demo/logs/on_rebuild_webhook.log`

## Webhook Evidence
`demo/logs/on_rebuild_webhook.log` contains:
- `POST /rebuild`
- `{"type":"rebuilt"}`

## Upstreaming
- PR URL: deferred
- Status: deferred (not opened in this run)

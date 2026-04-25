# Progress

Updated: 2026-04-25

---

## Status

| Task | Status |
|------|--------|
| Add `noServe` and `onRebuildURL` vars to `commands.go` | pending |
| Register `--no-serve` and `--on-rebuild` flags in `dev.go` `init()` | pending |
| Add `postRebuildWebhook` helper to `dev.go` | pending |
| Fire webhook in `OnRebuild` callback | pending |
| Conditionally skip `server.Serve` when `--no-serve` set | pending |
| `go build ./cmd/kiln` passes | pending |
| Smoke test `--no-serve` | pending |
| Smoke test `--on-rebuild` | pending |
| `go test ./...` passes | pending |
| Tag `v0.9.5-forge.1` | pending |
| Open upstream PR to `otaleghani/kiln` | pending |

## Current State

Not started. Ready to implement. All decisions resolved.

## Next Action

Begin Step 1: add vars to `internal/cli/commands.go`. See PLAN.md.

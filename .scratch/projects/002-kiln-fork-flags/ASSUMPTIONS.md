# Assumptions

Updated: 2026-04-25

---

## A-001: This repo is already the correct fork baseline

The current `kiln-fork` repo is a fork of `otaleghani/kiln`. It already has the commit history
and structure of kiln. The module path in `go.mod` is still `github.com/otaleghani/kiln` — this
does not need to change for the fork to work correctly in the forge-v2 stack.

## A-002: Only `internal/cli/dev.go` and `internal/cli/commands.go` need to change

The two new flags (`--no-serve`, `--on-rebuild`) belong entirely in the dev command. No other
packages need modification. The webhook POST uses only `net/http` from the standard library.

## A-003: Flag variables follow the existing package-level var pattern

The existing code puts all flag variables as package-level vars in `commands.go`. The two new
vars (`noServe bool`, `onRebuildURL string`) follow the same pattern.

## A-004: The webhook is fire-and-forget with a 5s timeout

If the POST fails (forge-overlay not running, network error), kiln logs an error and continues.
The build is not aborted. This matches the design in 001-forge-v2-refactor-overview/DECISIONS.md D-002.

## A-005: `--no-serve` only suppresses `server.Serve(...)` — the watcher still runs

When `--no-serve` is set, `runDev` still performs the initial build, populates the mtime store,
builds the dep graph, and starts the watcher goroutine. Only the final `server.Serve(ctx, ...)`
call is skipped. The process blocks on `<-ctx.Done()` instead.

## A-006: The existing test suite must pass unchanged

The diff adds new code paths behind new flags. No existing behavior changes. All existing tests
pass without modification. New tests cover the new code paths.

## A-007: devenv shell is required for `go test` and `go build`

Per REPO_RULES.md: use `devenv shell -- go test ./...` for all Go execution.

## A-008: Tag format is `v0.9.5-forge.1`

This follows the kiln version baseline (v0.9.5) plus a forge-specific suffix. Git tag is
created after tests pass and smoke test is confirmed.

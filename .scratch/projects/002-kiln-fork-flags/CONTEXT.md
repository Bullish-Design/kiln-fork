# Context

Updated: 2026-04-25

---

## What This Project Is

M1 of the forge-v2 refactor (see `.scratch/projects/001-forge-v2-refactor-overview/` for the
full system overview). This project is scoped entirely to the `kiln-fork` repo.

The only work here is adding two flags to `kiln dev`:

1. `--no-serve` — suppresses `server.Serve(...)`, lets kiln run as a headless build-and-watch process
2. `--on-rebuild <url>` — POSTs `{"type":"rebuilt"}` to the given URL after each successful rebuild

## Current State

Not started. All decisions resolved. Ready to implement.

## Key Files

| File | What it does |
|------|-------------|
| `internal/cli/commands.go` | Package-level flag vars and flag name constants. Add `noServe`, `onRebuildURL` here. |
| `internal/cli/dev.go` | `cmdDev`, `init()`, `runDev`. All changes happen here. |

## What "done" looks like

- Both flags work as described in PLAN.md
- `go test ./...` passes with no changes to existing tests
- Tagged `v0.9.5-forge.1`
- Upstream PR filed to `otaleghani/kiln`

## Resumption Instructions

1. Read CRITICAL_RULES.md — especially NO SUBAGENTS rule
2. Read this file
3. Check PROGRESS.md for current task status
4. Read PLAN.md for implementation steps
5. Pick up the first pending task and continue

## System Context

- The forge-v2 system has 5 components; this repo only implements M1
- After this is tagged, M2 (forge-overlay) begins in a separate repo
- obsidian-ops and obsidian-agent already exist and need no changes
- forge (Python orchestrator, M3) will start kiln-fork with these new flags

# Plan

Updated: 2026-04-25

---

> **CRITICAL RULE — NO SUBAGENTS**
> NEVER use the Task tool. Do ALL work directly. No delegation. No exceptions.

---

## Objective

Add two flags to `kiln dev` in this fork:

1. `--no-serve` — suppress `server.Serve(...)` so kiln runs as a pure build-and-watch process
2. `--on-rebuild <url>` — POST `{"type":"rebuilt"}` to a URL after each successful incremental build

Diff from upstream kiln should be ≤50 lines when complete.

---

## Files to Change

| File | Change |
|------|--------|
| `internal/cli/commands.go` | Add `noServe bool` and `onRebuildURL string` package-level vars |
| `internal/cli/dev.go` | Register flags in `init()`, implement `--no-serve` and webhook in `runDev` |

No other files change.

---

## Implementation Steps

### Step 1 — Add vars to `commands.go`

In the `var (...)` block (after `accentColor string`), add:

```go
noServe      bool   // Skip starting the HTTP server in dev mode
onRebuildURL string // URL to POST {"type":"rebuilt"} after each successful rebuild
```

### Step 2 — Register flags in `dev.go` `init()`

Add after the existing `port` flag registration:

```go
cmdDev.Flags().
    BoolVar(&noServe, "no-serve", false, "Skip starting the HTTP server (for use with an external overlay)")
cmdDev.Flags().
    StringVar(&onRebuildURL, "on-rebuild", "", "URL to POST {\"type\":\"rebuilt\"} after each successful rebuild")
```

### Step 3 — Add webhook helper to `dev.go`

Add a package-level helper (below the imports, before `cmdDev`):

```go
var rebuildClient = &http.Client{Timeout: 5 * time.Second}

func postRebuildWebhook(url string, log *slog.Logger) {
    body := strings.NewReader(`{"type":"rebuilt"}`)
    resp, err := rebuildClient.Post(url, "application/json", body)
    if err != nil {
        log.Error("on-rebuild webhook failed", "url", url, "err", err)
        return
    }
    resp.Body.Close()
}
```

Add required imports: `"net/http"`, `"strings"`, `"time"` (plus existing imports).

### Step 4 — Fire webhook in `OnRebuild` callback

At the end of the `OnRebuild` func body (after `graph.UpdateFiles(...)`), add:

```go
if onRebuildURL != "" {
    postRebuildWebhook(onRebuildURL, log)
}
return nil
```

### Step 5 — Conditionally serve in `runDev`

Replace the final block:

```go
// Serve on main goroutine
localBaseURL := "http://localhost:" + port
server.Serve(ctx, port, builder.OutputDir, localBaseURL, log)
```

With:

```go
if noServe {
    // Block until signal; serving is handled by an external overlay
    <-ctx.Done()
} else {
    localBaseURL := "http://localhost:" + port
    server.Serve(ctx, port, builder.OutputDir, localBaseURL, log)
}
```

---

## Verification Steps

### Smoke test 1 — `--no-serve`

```bash
devenv shell -- go build -o /tmp/kiln-fork ./cmd/kiln
/tmp/kiln-fork dev --no-serve --input ./vault --output ./public &
# Confirm: no HTTP server starts, process stays running, ctrl-C exits cleanly
```

### Smoke test 2 — `--on-rebuild`

```bash
# Terminal 1: simple echo server
devenv shell -- go run .scratch/projects/002-kiln-fork-flags/echo_server.go
# Terminal 2: kiln in watch mode
/tmp/kiln-fork dev --no-serve --on-rebuild http://localhost:9999/rebuild --input ./vault --output ./public
# Edit a file in ./vault
# Terminal 1 should show: POST /rebuild with body {"type":"rebuilt"}
```

### Test suite

```bash
devenv shell -- go test ./...
```

All existing tests must pass.

---

## Acceptance Criteria

- [ ] `go build ./cmd/kiln` succeeds
- [ ] `kiln dev --no-serve` runs watcher without HTTP server, exits on SIGINT
- [ ] `kiln dev --on-rebuild http://...` POSTs `{"type":"rebuilt"}` after each rebuild
- [ ] Webhook failure is logged, build continues (not aborted)
- [ ] `go test ./...` passes with no changes to existing tests
- [ ] Diff from upstream kiln ≤50 lines
- [ ] Tagged `v0.9.5-forge.1`

---

> **REMINDER — NO SUBAGENTS**
> NEVER use the Task tool. Do ALL work directly. No delegation. No exceptions.

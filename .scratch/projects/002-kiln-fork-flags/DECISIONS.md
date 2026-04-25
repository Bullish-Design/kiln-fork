# Decisions

Updated: 2026-04-25

All architecture decisions for kiln-fork are inherited from project 001. The decisions relevant
to this repo are summarized here for quick reference.

---

## D-001: Flag placement — package-level vars in `commands.go` [RESOLVED]

**Decision:** Add `noServe bool` and `onRebuildURL string` as package-level vars in
`internal/cli/commands.go`, matching the existing pattern for all other flag variables.

**Rationale:** Consistency with existing code. All other flag vars (`port`, `flatUrls`, etc.)
live there. Flag registration stays in `dev.go`'s `init()`.

---

## D-002: Webhook implementation — `net/http` with 5s timeout [RESOLVED]

**Decision:** Use a package-level `http.Client` with `Timeout: 5 * time.Second`. POST
`{"type":"rebuilt"}` with `Content-Type: application/json`. Log error, do not abort build.

**Rationale:** No new dependencies needed. Standard library is sufficient. Timeout prevents
a hung forge-overlay from stalling the watcher.

---

## D-003: `--no-serve` behavior — block on ctx.Done() instead of serving [RESOLVED]

**Decision:** When `--no-serve` is set, skip the `server.Serve(ctx, ...)` call. The process
must still block until SIGINT/SIGTERM. Use `<-ctx.Done()` after starting the watcher goroutine.

**Rationale:** Without this, the process would exit immediately after starting the watcher.
`ctx` is already set up with `signal.NotifyContext` — blocking on it is the right idiom.

---

## D-004: Webhook fires from `OnRebuild` callback, not the watcher directly [RESOLVED]

**Decision:** The `onRebuildURL` POST happens at the end of the `OnRebuild` func in `runDev`,
after `builder.IncrementalBuild(...)` and graph refresh complete successfully.

**Rationale:** Fires once per completed rebuild. The forge-overlay SSE fires exactly when the
build is done — no timing ambiguity. Matches D-002 from 001/DECISIONS.md.

---

## D-005: Upstream PR filed after tag [RESOLVED from 001]

**Decision:** After tagging `v0.9.5-forge.1`, open a PR to `otaleghani/kiln` with the two flags.
If accepted, kiln-fork eventually becomes a direct kiln dependency with no fork needed.

**Rationale:** Minimal diff, easy to upstream. Reduces long-term fork maintenance burden.

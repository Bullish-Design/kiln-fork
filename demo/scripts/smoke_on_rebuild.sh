#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/../.."

mkdir -p demo/logs demo/public demo/out demo/tmp
LOG_FILE="demo/logs/smoke_on_rebuild.log"
RESULT_FILE="demo/logs/smoke_on_rebuild.result"
WEBHOOK_LOG="demo/logs/on_rebuild_webhook.log"
ECHO_OUT="demo/logs/echo_server.stdout"
ECHO_ERR="demo/logs/echo_server.stderr"
TARGET_FILE="demo/vault/snippets/reusable-cta.md"
BACKUP_FILE="demo/tmp/reusable-cta.md.bak"
WEBHOOK_URL="http://127.0.0.1:9999/rebuild"

: > "$LOG_FILE"
: > "$RESULT_FILE"
: > "$WEBHOOK_LOG"
: > "$ECHO_OUT"
: > "$ECHO_ERR"

if [ ! -x ./demo/out/kiln ]; then
  go build -o ./demo/out/kiln ./cmd/kiln
fi

cp "$TARGET_FILE" "$BACKUP_FILE"
restore_file() {
  if [ -f "$BACKUP_FILE" ]; then
    mv "$BACKUP_FILE" "$TARGET_FILE"
  fi
}

node -e '
const fs = require("fs");
const http = require("http");
const logPath = "./demo/logs/on_rebuild_webhook.log";
const server = http.createServer((req, res) => {
  let body = "";
  req.on("data", chunk => { body += chunk; });
  req.on("end", () => {
    fs.appendFileSync(logPath, `${req.method} ${req.url}\n${body}\n`);
    res.statusCode = 204;
    res.end();
  });
});
server.listen(9999, "127.0.0.1");
' >"$ECHO_OUT" 2>"$ECHO_ERR" &
ECHO_PID=$!

./demo/out/kiln dev --input ./demo/vault --output ./demo/public --no-serve --on-rebuild "$WEBHOOK_URL" >"$LOG_FILE" 2>&1 &
DEV_PID=$!

stop_process() {
  local pid="$1"
  if ! kill -0 "$pid" 2>/dev/null; then
    return 0
  fi
  kill -INT "$pid" 2>/dev/null || true
  for _ in $(seq 1 20); do
    if ! kill -0 "$pid" 2>/dev/null; then
      return 0
    fi
    sleep 0.2
  done
  kill -TERM "$pid" 2>/dev/null || true
  for _ in $(seq 1 10); do
    if ! kill -0 "$pid" 2>/dev/null; then
      return 0
    fi
    sleep 0.2
  done
  kill -KILL "$pid" 2>/dev/null || true
}

cleanup() {
  stop_process "$DEV_PID"
  wait "$DEV_PID" 2>/dev/null || true
  stop_process "$ECHO_PID"
  wait "$ECHO_PID" 2>/dev/null || true
  restore_file
}
trap cleanup EXIT

# Wait until initial build has completed and watcher is active.
for _ in $(seq 1 60); do
  if rg -q "Build complete" "$LOG_FILE"; then
    break
  fi
  sleep 0.5
done

printf "\n" >> "$TARGET_FILE"

for _ in $(seq 1 40); do
  if rg -q 'POST /rebuild' "$WEBHOOK_LOG" && rg -F -q '{"type":"rebuilt"}' "$WEBHOOK_LOG"; then
    cleanup
    trap - EXIT
    echo "PASS: webhook received expected payload" | tee "$RESULT_FILE"
    exit 0
  fi
  sleep 1
done

echo "FAIL: webhook payload not observed" | tee "$RESULT_FILE"
exit 1

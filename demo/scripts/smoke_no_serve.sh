#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/../.."

mkdir -p demo/logs demo/public demo/out demo/tmp
LOG_FILE="demo/logs/smoke_no_serve.log"
RESULT_FILE="demo/logs/smoke_no_serve.result"
HTTP_ERR="demo/logs/smoke_no_serve_http.err"
HTTP_OUT="demo/logs/smoke_no_serve_http.out"
PORT=18777

: > "$LOG_FILE"
: > "$RESULT_FILE"
: > "$HTTP_ERR"
: > "$HTTP_OUT"

if [ ! -x ./demo/out/kiln ]; then
  go build -o ./demo/out/kiln ./cmd/kiln
fi

./demo/out/kiln dev --input ./demo/vault --output ./demo/public --port "$PORT" --no-serve >"$LOG_FILE" 2>&1 &
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
}
trap cleanup EXIT

sleep 3
if curl -sS "http://127.0.0.1:${PORT}/" >"$HTTP_OUT" 2>"$HTTP_ERR"; then
  echo "FAIL: server unexpectedly responded on port ${PORT}" | tee "$RESULT_FILE"
  exit 1
fi

cleanup
trap - EXIT

echo "PASS: no server on port and process handled shutdown" | tee "$RESULT_FILE"

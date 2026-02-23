#!/usr/bin/env bash
# test_shortener.sh — Build, run, test, and clean up the link shortener stack.
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

PASSED=0
FAILED=0

cleanup() {
    echo "==> Cleaning up..."
    docker compose down --rmi local -v 2>/dev/null || true
}
trap cleanup EXIT

echo "==> Building and starting Docker Compose stack..."
docker compose up -d --build

echo "==> Waiting for services to start..."
sleep 5

# --- Test 1: POST /shorten returns 201 and a key ---
echo "==> Test 1: POST /shorten with valid URL..."
SHORTEN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/shorten \
    -H "Content-Type: application/json" \
    -d '{"url": "https://go.dev"}')
SHORTEN_BODY=$(echo "$SHORTEN_RESPONSE" | head -1)
SHORTEN_STATUS=$(echo "$SHORTEN_RESPONSE" | tail -1)
SHORTEN_KEY=$(echo "$SHORTEN_BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['key'])")

if [ "$SHORTEN_STATUS" = "201" ] && [ -n "$SHORTEN_KEY" ]; then
    echo "PASS: got 201 with key '$SHORTEN_KEY'"
    PASSED=$((PASSED + 1))
else
    echo "FAIL: expected 201 with key, got status=$SHORTEN_STATUS body=$SHORTEN_BODY"
    FAILED=$((FAILED + 1))
fi

# --- Test 2: GET /r/{key} returns 302 redirect ---
echo "==> Test 2: GET /r/$SHORTEN_KEY returns 302 redirect..."
REDIRECT_RESPONSE=$(curl -s -o /dev/null -w "%{http_code} %{redirect_url}" http://localhost:8080/r/$SHORTEN_KEY)
REDIRECT_STATUS=$(echo "$REDIRECT_RESPONSE" | awk '{print $1}')
REDIRECT_URL=$(echo "$REDIRECT_RESPONSE" | awk '{print $2}')

if [ "$REDIRECT_STATUS" = "302" ] && [ "$REDIRECT_URL" = "https://go.dev/" ]; then
    echo "PASS: got 302 redirecting to https://go.dev/"
    PASSED=$((PASSED + 1))
else
    echo "FAIL: expected 302 → https://go.dev/, got status=$REDIRECT_STATUS url=$REDIRECT_URL"
    FAILED=$((FAILED + 1))
fi

# --- Test 3: GET /r/nonexistent returns 404 ---
echo "==> Test 3: GET /r/nonexistent returns 404..."
NOT_FOUND_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/r/nonexistent)

if [ "$NOT_FOUND_STATUS" = "404" ]; then
    echo "PASS: got 404 for nonexistent key"
    PASSED=$((PASSED + 1))
else
    echo "FAIL: expected 404, got $NOT_FOUND_STATUS"
    FAILED=$((FAILED + 1))
fi

# --- Test 4: POST /shorten with empty URL returns 400 ---
echo "==> Test 4: POST /shorten with empty URL..."
EMPTY_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/shorten \
    -H "Content-Type: application/json" \
    -d '{"url": ""}')

if [ "$EMPTY_STATUS" = "400" ]; then
    echo "PASS: got 400 for empty URL"
    PASSED=$((PASSED + 1))
else
    echo "FAIL: expected 400, got $EMPTY_STATUS"
    FAILED=$((FAILED + 1))
fi

# --- Test 5: Persistence — data survives app restart ---
echo "==> Test 5: Persistence — restarting app container..."
docker compose restart app
sleep 5

PERSIST_RESPONSE=$(curl -s -o /dev/null -w "%{http_code} %{redirect_url}" http://localhost:8080/r/$SHORTEN_KEY)
PERSIST_STATUS=$(echo "$PERSIST_RESPONSE" | awk '{print $1}')
PERSIST_URL=$(echo "$PERSIST_RESPONSE" | awk '{print $2}')

if [ "$PERSIST_STATUS" = "302" ] && [ "$PERSIST_URL" = "https://go.dev/" ]; then
    echo "PASS: data survived app restart (Redis persistence)"
    PASSED=$((PASSED + 1))
else
    echo "FAIL: expected 302 → https://go.dev/ after restart, got status=$PERSIST_STATUS url=$PERSIST_URL"
    FAILED=$((FAILED + 1))
fi

# --- Test 6: Load balancing — multiple hostnames from /health ---
echo "==> Test 6: Load balancing — hitting /health 10 times..."
HOSTNAMES=""
for i in $(seq 1 10); do
    H=$(curl -s http://localhost:8080/health | python3 -c "import sys,json; print(json.load(sys.stdin)['hostname'])")
    HOSTNAMES="$HOSTNAMES $H"
done
UNIQUE=$(echo "$HOSTNAMES" | tr ' ' '\n' | sort -u | grep -c .)

if [ "$UNIQUE" -gt 1 ]; then
    echo "PASS: saw $UNIQUE unique hostnames (load balancing works)"
    PASSED=$((PASSED + 1))
else
    echo "FAIL: expected >1 unique hostname, got $UNIQUE (hostnames:$HOSTNAMES)"
    FAILED=$((FAILED + 1))
fi

# --- Summary ---
echo ""
echo "==> Results: $PASSED passed, $FAILED failed out of 6 tests"
if [ "$FAILED" -gt 0 ]; then
    exit 1
fi
echo "==> All tests passed!"

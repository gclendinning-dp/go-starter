#!/usr/bin/env bash
# test_compose.sh — Build, run, test, and clean up the Docker Compose stack.
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

cleanup() {
    echo "==> Cleaning up..."
    docker compose down --rmi local 2>/dev/null || true
}
trap cleanup EXIT

echo "==> Building and starting Docker Compose stack..."
docker compose up -d --build

echo "==> Waiting for services to start..."
sleep 5

# --- Test 1: Multiple hostnames prove load balancing ---
echo "==> Test 1: Checking load balancing across replicas..."
HOSTNAMES=""
for i in $(seq 1 10); do
    HOSTNAME=$(curl -s http://localhost:8080/message | python3 -c "import sys,json; print(json.load(sys.stdin)['hostname'])")
    HOSTNAMES="$HOSTNAMES $HOSTNAME"
done

UNIQUE=$(echo "$HOSTNAMES" | tr ' ' '\n' | sort -u | wc -l | tr -d ' ')
if [ "$UNIQUE" -gt 1 ]; then
    echo "PASS: saw $UNIQUE unique hostnames — load balancing works"
else
    echo "FAIL: only saw 1 hostname — expected multiple (got:$HOSTNAMES)"
    exit 1
fi

# --- Test 2: POST /greet through load balancer ---
echo "==> Test 2: POST /greet through load balancer..."
GREET_RESPONSE=$(curl -s -X POST http://localhost:8080/greet \
    -H "Content-Type: application/json" \
    -d '{"name": "Alice"}')

# Check that the greeting contains "Alice"
if echo "$GREET_RESPONSE" | python3 -c "import sys,json; g=json.load(sys.stdin)['greeting']; assert 'Alice' in g" 2>/dev/null; then
    echo "PASS: greeting contains 'Alice'"
else
    echo "FAIL: expected greeting with 'Alice', got: $GREET_RESPONSE"
    exit 1
fi

echo "==> All tests passed!"

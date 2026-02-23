#!/usr/bin/env bash
# test_container.sh — Build, run, test, and clean up the Docker container.
set -e

IMAGE_NAME="go-server-test"
CONTAINER_NAME="go-server-test-run"
PORT=8080

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "==> Building Docker image..."
docker build -t "$IMAGE_NAME" "$SCRIPT_DIR"

echo "==> Starting container on port $PORT..."
docker run -d --name "$CONTAINER_NAME" -p "$PORT:$PORT" "$IMAGE_NAME"

# Give the server a moment to start.
sleep 2

echo "==> Testing /message endpoint..."
RESPONSE=$(curl -s "http://localhost:$PORT/message")
EXPECTED='{"message": "Hello from the backend!"}'

if [ "$RESPONSE" = "$EXPECTED" ]; then
    echo "PASS: got expected response"
else
    echo "FAIL: expected '$EXPECTED', got '$RESPONSE'"
    docker stop "$CONTAINER_NAME" >/dev/null 2>&1
    docker rm "$CONTAINER_NAME" >/dev/null 2>&1
    docker rmi "$IMAGE_NAME" >/dev/null 2>&1
    exit 1
fi

echo "==> Cleaning up..."
docker stop "$CONTAINER_NAME" >/dev/null 2>&1
docker rm "$CONTAINER_NAME" >/dev/null 2>&1
docker rmi "$IMAGE_NAME" >/dev/null 2>&1

echo "==> All done! Container built, ran, responded correctly, and cleaned up."

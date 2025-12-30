#!/bin/bash
set -e

echo "=== casspeed Test Suite ==="

BINDIR="binaries"
SERVER="$BINDIR/casspeed"
CLIENT="$BINDIR/casspeed-cli"
PORT=64580

if [ ! -f "$SERVER" ]; then
    echo "Error: Server binary not found. Run 'make build' first."
    exit 1
fi

if [ ! -f "$CLIENT" ]; then
    echo "Error: Client binary not found. Run 'make build' first."
    exit 1
fi

echo "✓ Binaries found"

echo "Starting server on port $PORT..."
$SERVER --port $PORT &
SERVER_PID=$!
sleep 3

if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "✗ Server failed to start"
    exit 1
fi

echo "✓ Server started (PID: $SERVER_PID)"

echo "Testing health endpoint..."
HEALTH=$(curl -s http://localhost:$PORT/health)
if echo "$HEALTH" | grep -q '"status":"ok"'; then
    echo "✓ Health check passed"
else
    echo "✗ Health check failed"
    kill $SERVER_PID
    exit 1
fi

echo "Testing API root..."
API=$(curl -s http://localhost:$PORT/api/v1/)
if echo "$API" | grep -q '"version":"v1"'; then
    echo "✓ API root check passed"
else
    echo "✗ API root check failed"
    kill $SERVER_PID
    exit 1
fi

echo "Testing web UI..."
UI=$(curl -s http://localhost:$PORT/)
if echo "$UI" | grep -q "casspeed"; then
    echo "✓ Web UI check passed"
else
    echo "✗ Web UI check failed"
    kill $SERVER_PID
    exit 1
fi

echo "Stopping server..."
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null

echo ""
echo "=== All Tests Passed ✓ ==="

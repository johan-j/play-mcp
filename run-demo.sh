#!/bin/bash

echo "Starting MCP Server with plugins..."

# Build the project
echo "Building the server..."
make build

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Build successful!"

# Start the server in the background
echo "Starting server on localhost:8080..."
./bin/mcp-server &
SERVER_PID=$!

# Wait for server to start
sleep 2

# Test the server
echo "Testing server endpoints..."
go run test-client.go

# Test individual plugin tools (example JSON-RPC calls)
echo ""
echo "=== Example Tool Calls ==="
echo ""

echo "Financial Plugin - Get Apple stock data:"
echo "You can test this with a WebSocket client by sending:"
echo '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "get_stock_data",
    "arguments": {
      "symbol": "AAPL"
    }
  }
}'

echo ""
echo "Housing Plugin - Search properties in San Francisco:"
echo "You can test this with a WebSocket client by sending:"
echo '{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "search_properties",
    "arguments": {
      "city": "San Francisco",
      "state": "CA",
      "maxPrice": 2000000
    }
  }
}'

echo ""
echo "Server is running with PID: $SERVER_PID"
echo "WebSocket endpoint: ws://localhost:8080/mcp"
echo "Press Ctrl+C to stop the server"

# Keep the script running and handle Ctrl+C
trap "echo 'Stopping server...'; kill $SERVER_PID; exit 0" INT

# Wait for the server process
wait $SERVER_PID

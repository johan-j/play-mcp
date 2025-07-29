#!/bin/bash

echo "Testing MCP Server WebSocket connection..."

# Start the server in the background
echo "Starting MCP server..."
./bin/mcp-server &
SERVER_PID=$!

# Wait for server to start
sleep 3

echo "Testing server endpoints..."

# Test health endpoint
echo "1. Health check:"
curl -s http://localhost:8080/health | jq

echo ""
echo "2. Available tools:"
curl -s http://localhost:8080/tools | jq '.tools | length'

echo ""
echo "3. Available resources:"
curl -s http://localhost:8080/resources | jq '.resources | length'

echo ""
echo "WebSocket endpoint available at: ws://localhost:8080/mcp"
echo ""
echo "To test WebSocket manually:"
echo "1. Install wscat: npm install -g wscat"
echo "2. Connect: wscat -c ws://localhost:8080/mcp"
echo "3. Send messages like:"
echo '   {"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}'
echo '   {"jsonrpc":"2.0","id":2,"method":"tools/list"}'
echo '   {"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_stock_data","arguments":{"symbol":"AAPL"}}}'
echo ""
echo "For GitHub Copilot integration, see: GITHUB_COPILOT_SETUP.md"

# Stop the server
echo ""
echo "Stopping test server..."
kill $SERVER_PID

echo "Test complete!"

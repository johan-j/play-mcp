#!/bin/bash

echo "🧪 Testing MCP WebSocket Connection..."

# Check if server is running
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "❌ Server not running. Starting server..."
    ./bin/mcp-server &
    SERVER_PID=$!
    sleep 3
    echo "✅ Server started with PID: $SERVER_PID"
else
    echo "✅ Server is already running"
    SERVER_PID=""
fi

# Test WebSocket with Node.js if available
if command -v node > /dev/null 2>&1; then
    echo "🔌 Testing WebSocket connection with Node.js..."
    
    cat > test-websocket.js << 'EOF'
const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:8080/mcp');

ws.on('open', function() {
    console.log('✅ WebSocket connected successfully!');
    
    // Send initialize message
    const initMessage = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "initialize",
        "params": {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {
                "name": "test-client",
                "version": "1.0.0"
            }
        }
    };
    
    ws.send(JSON.stringify(initMessage));
});

ws.on('message', function(data) {
    const response = JSON.parse(data);
    console.log('📨 Received:', JSON.stringify(response, null, 2));
    
    if (response.id === 1) {
        // After initialize, list tools
        const toolsMessage = {
            "jsonrpc": "2.0",
            "id": 2,
            "method": "tools/list"
        };
        ws.send(JSON.stringify(toolsMessage));
    } else if (response.id === 2) {
        console.log('🛠️  Available tools:', response.result.tools.length);
        ws.close();
    }
});

ws.on('error', function(error) {
    console.log('❌ WebSocket error:', error.message);
    process.exit(1);
});

ws.on('close', function() {
    console.log('🔌 WebSocket connection closed');
    process.exit(0);
});

// Timeout after 10 seconds
setTimeout(() => {
    console.log('⏰ Test timeout');
    ws.close();
    process.exit(1);
}, 10000);
EOF

    # Install ws package if needed
    if [ ! -d "node_modules" ]; then
        echo "📦 Installing WebSocket package..."
        npm init -y > /dev/null 2>&1
        npm install ws > /dev/null 2>&1
    fi
    
    node test-websocket.js
    
    # Clean up
    rm -f test-websocket.js
    
else
    echo "⚠️  Node.js not available, skipping WebSocket test"
    echo "💡 Install Node.js to test WebSocket connection"
    echo "💡 Or use: npm install -g wscat && wscat -c ws://localhost:8080/mcp"
fi

# Clean up server if we started it
if [ ! -z "$SERVER_PID" ]; then
    echo "🛑 Stopping test server..."
    kill $SERVER_PID
fi

echo ""
echo "🎯 Connection Summary:"
echo "✅ WebSocket Endpoint: ws://localhost:8080/mcp"
echo "✅ HTTP Health Check: http://localhost:8080/health"
echo "✅ VS Code Settings: .vscode/settings.json"
echo ""
echo "📖 Next Steps:"
echo "1. Open VS Code: code ."
echo "2. Install MCP extension"
echo "3. Start server: ./bin/mcp-server"
echo "4. Ask Copilot: 'Get Apple stock data'"

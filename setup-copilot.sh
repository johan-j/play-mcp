#!/bin/bash

# GitHub Copilot MCP Server Setup Script

echo "🚀 Setting up MCP Server for GitHub Copilot..."

# Get the absolute path to the project
PROJECT_DIR=$(pwd)
BINARY_PATH="$PROJECT_DIR/bin/mcp-server"

echo "Project directory: $PROJECT_DIR"
echo "Binary path: $BINARY_PATH"

# Build the server
echo "📦 Building MCP server..."
make build

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"

# Check if binary exists
if [ ! -f "$BINARY_PATH" ]; then
    echo "❌ Binary not found at $BINARY_PATH"
    exit 1
fi

# Test the server
echo "🧪 Testing server..."
./bin/mcp-server &
SERVER_PID=$!

sleep 3

# Test endpoints
echo "📊 Testing endpoints..."
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
if [ $? -eq 0 ]; then
    echo "✅ Health endpoint working"
    echo "$HEALTH_RESPONSE" | jq
else
    echo "❌ Health endpoint failed"
fi

TOOLS_COUNT=$(curl -s http://localhost:8080/tools | jq '.tools | length' 2>/dev/null)
if [ ! -z "$TOOLS_COUNT" ]; then
    echo "✅ Tools endpoint working - $TOOLS_COUNT tools available"
else
    echo "❌ Tools endpoint failed"
fi

# Stop test server
kill $SERVER_PID
echo "🛑 Test server stopped"

echo ""
echo "🎯 GitHub Copilot Integration Options:"
echo ""
echo "1. 🌐 WebSocket Endpoint:"
echo "   ws://localhost:8080/mcp"
echo ""
echo "2. 🔗 HTTP REST API:"
echo "   Health:    http://localhost:8080/health"
echo "   Tools:     http://localhost:8080/tools"
echo "   Resources: http://localhost:8080/resources"
echo ""
echo "3. 📚 Available Tools (9 total):"
echo "   Financial Tools:"
echo "   - get_stock_data: Get stock information by symbol"
echo "   - search_companies: Search companies by name"
echo "   - get_market_summary: Get market overview"
echo "   - get_historical_data: Get historical stock data"
echo ""
echo "   Housing Tools:"
echo "   - search_properties: Search real estate properties"
echo "   - get_property_details: Get property information"
echo "   - get_market_stats: Get market statistics"
echo "   - get_price_history: Get price history"
echo "   - estimate_property_value: Get value estimates"
echo ""
echo "🔧 Next Steps:"
echo ""
echo "1. Start the server:"
echo "   ./bin/mcp-server"
echo ""
echo "2. For VS Code integration:"
echo "   - Install an MCP extension from VS Code marketplace"
echo "   - Configure to use: ws://localhost:8080/mcp"
echo ""
echo "3. For custom integrations:"
echo "   - See: GITHUB_COPILOT_SETUP.md"
echo "   - Test with: ./test-websocket.sh"
echo ""
echo "🎉 Your MCP server is ready for GitHub Copilot integration!"
echo ""
echo "💡 Pro tip: Run with -debug flag for detailed logging:"
echo "   ./bin/mcp-server -debug"

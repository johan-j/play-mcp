# How to Connect GitHub Copilot to Your MCP Server

## 🚀 Quick Connection Guide

Your MCP server endpoint: **`ws://localhost:8080/mcp`**

## Method 1: VS Code with MCP Extension (Recommended)

### Step 1: Install MCP Extension
1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X / Cmd+Shift+X)
3. Search for "MCP" or "Model Context Protocol"
4. Install a compatible MCP extension

### Step 2: Configure VS Code Settings
Add this to your workspace settings (`.vscode/settings.json`):

```json
{
  "mcp.servers": {
    "financial-housing-data": {
      "name": "Financial & Housing Data Server",
      "url": "ws://localhost:8080/mcp",
      "description": "Provides financial market data and real estate information",
      "capabilities": ["tools", "resources"],
      "enabled": true
    }
  },
  "mcp.enableLogging": true,
  "mcp.autoConnect": true
}
```

### Step 3: Start Your Server
```bash
# In your terminal
./bin/mcp-server

# Or with debug logging
./bin/mcp-server -debug
```

### Step 4: Verify Connection
- Check VS Code status bar for MCP connection indicator
- Look for "Financial & Housing Data Server" in connected servers
- Test with a simple query like "Get Apple stock data"

## Method 2: GitHub Copilot Chat Extension

If using the GitHub Copilot Chat extension in VS Code:

### Step 1: Configure Copilot Settings
Add to your VS Code settings:

```json
{
  "github.copilot.advanced.agent.tools": {
    "financial-housing-mcp": {
      "type": "mcp",
      "url": "ws://localhost:8080/mcp",
      "enabled": true
    }
  }
}
```

### Step 2: Use in Copilot Chat
In the Copilot Chat panel, you can now ask:
- "Get current Apple stock price"
- "Search for properties in San Francisco under $2M"
- "Show me market summary"

## Method 3: Direct WebSocket Connection (For Custom Clients)

If you're building a custom Copilot client:

```javascript
const socket = new WebSocket('ws://localhost:8080/mcp');

// Initialize connection
socket.onopen = () => {
  socket.send(JSON.stringify({
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {
        "name": "github-copilot-client",
        "version": "1.0.0"
      }
    }
  }));
};

// Handle responses
socket.onmessage = (event) => {
  const response = JSON.parse(event.data);
  console.log('MCP Response:', response);
};
```

## Method 4: Test Connection First

Before setting up Copilot, test the connection manually:

### Using wscat:
```bash
# Install wscat if you haven't
npm install -g wscat

# Connect to your server
wscat -c ws://localhost:8080/mcp

# Send initialize message
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}

# List available tools
{"jsonrpc":"2.0","id":2,"method":"tools/list"}

# Test a tool call
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_stock_data","arguments":{"symbol":"AAPL"}}}
```

### Using curl (HTTP endpoints):
```bash
# Test server health
curl http://localhost:8080/health

# List tools
curl http://localhost:8080/tools | jq

# List resources
curl http://localhost:8080/resources | jq
```

## 🔧 Troubleshooting

### Connection Issues:
1. **Server not running:** Make sure `./bin/mcp-server` is running
2. **Port busy:** Try different port: `./bin/mcp-server -port 9090`
3. **WebSocket blocked:** Check firewall settings

### VS Code Issues:
1. **Extension not found:** Search for alternative MCP extensions
2. **Settings not applied:** Reload VS Code window (Ctrl+Shift+P → "Reload Window")
3. **Connection failed:** Check VS Code Developer Console (Help → Toggle Developer Tools)

### Debug Mode:
```bash
# Run server with detailed logging
./bin/mcp-server -debug
```

## 🎯 Example Copilot Interactions

Once connected, you can ask Copilot:

**Financial Queries:**
- "What's Apple's current stock price?"
- "Find tech companies for investment"
- "Show me today's market summary"
- "Get historical data for Tesla stock"

**Real Estate Queries:**
- "Search properties in San Francisco under $2M"
- "Get property details for 123 Main St"
- "What are market stats for Los Angeles?"
- "Estimate value of a property in Seattle"

**Code Generation:**
- "Write a Python function to analyze stock data"
- "Create a real estate investment calculator"
- "Build a portfolio tracker using current market data"

## 📝 Configuration Files

Your repo already includes:
- `vscode-mcp-config.json` - Ready-to-use VS Code settings
- `GITHUB_COPILOT_SETUP.md` - Detailed setup guide
- `test-websocket.sh` - Connection testing script

## 🎉 Success Indicators

You'll know it's working when:
- ✅ VS Code shows MCP server connected
- ✅ Copilot can answer financial/housing questions with real data
- ✅ Tool calls appear in server logs (debug mode)
- ✅ HTTP endpoints return data: `curl http://localhost:8080/health`

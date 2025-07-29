# GitHub Copilot MCP Server Setup

This guide shows how to connect your MCP server to GitHub Copilot for enhanced code assistance with financial and housing data.

## Quick Start

1. **Build and Start the Server:**
   ```bash
   make build
   ./bin/mcp-server
   ```

2. **Verify Server is Running:**
   ```bash
   curl http://localhost:8080/health
   ```

## GitHub Copilot Integration

### Method 1: VS Code Extension (Recommended)

1. **Install MCP Extension in VS Code:**
   - Search for "MCP" or "Model Context Protocol" extensions
   - Install a compatible MCP extension

2. **Configure the Extension:**
   - Open VS Code settings
   - Find MCP extension settings
   - Add server endpoint: `ws://localhost:8080/mcp`

### Method 2: Direct WebSocket Connection

If using a custom Copilot client or integration:

```javascript
// Connect to WebSocket endpoint
const socket = new WebSocket('ws://localhost:8080/mcp');

// Initialize MCP connection
socket.send(JSON.stringify({
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "github-copilot",
      "version": "1.0.0"
    }
  }
}));
```

### Method 3: HTTP REST API

For simple integrations, use the REST endpoints:

```bash
# Get available tools
curl http://localhost:8080/tools

# Get server health
curl http://localhost:8080/health

# Get available resources
curl http://localhost:8080/resources
```

## Available Tools for Copilot

### Financial Data Tools (4 available)
- **get_stock_data** - Get current stock information by symbol
- **search_companies** - Search for companies by name
- **get_market_summary** - Get market overview and indices
- **get_historical_data** - Get historical stock price data

### Housing Data Tools (5 available)
- **search_properties** - Search real estate properties with filters
- **get_property_details** - Get detailed property information
- **get_market_stats** - Get market statistics for areas
- **get_price_history** - Get property price history
- **estimate_property_value** - Get property value estimates

## Example Usage in Code

When GitHub Copilot is connected, it can help with:

### Financial Analysis Code
```python
# Copilot can now access real stock data
import requests

def analyze_stock(symbol):
    # Copilot will know current stock prices and can suggest analysis
    # based on real market data from your MCP server
    pass
```

### Real Estate Applications
```python
def find_investment_properties(city, max_price):
    # Copilot can access property listings and market data
    # to suggest relevant investment opportunities
    pass
```

## Server Configuration

### Environment Variables
```bash
export MCP_HOST=localhost
export MCP_PORT=8080
export MCP_DEBUG=false
```

### Command Line Options
```bash
./bin/mcp-server -help
# Options:
#   -host string    Server host (default "localhost")
#   -port int       Server port (default 8080)
#   -debug          Enable debug logging
```

## Testing the Integration

### WebSocket Test with wscat
```bash
# Install wscat
npm install -g wscat

# Connect to server
wscat -c ws://localhost:8080/mcp

# Send initialize message
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}

# List available tools
{"jsonrpc":"2.0","id":2,"method":"tools/list"}

# Test financial tool
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_stock_data","arguments":{"symbol":"AAPL"}}}
```

### HTTP API Test
```bash
# Test all endpoints
curl http://localhost:8080/health | jq
curl http://localhost:8080/tools | jq
curl http://localhost:8080/resources | jq
```

## Troubleshooting

### Common Issues

1. **Server not starting:**
   - Check if port 8080 is available: `lsof -i :8080`
   - Try different port: `./bin/mcp-server -port 9090`

2. **WebSocket connection failed:**
   - Check firewall settings
   - Verify server is running: `curl http://localhost:8080/health`

3. **Tools not working:**
   - Enable debug mode: `./bin/mcp-server -debug`
   - Check server logs for errors

### Debug Mode
```bash
./bin/mcp-server -debug
```
This will show detailed logs of:
- Plugin registration
- Client connections
- Tool calls
- Error messages

## Security Considerations

### For Production Use
- Add authentication to endpoints
- Use HTTPS/WSS in production
- Implement rate limiting
- Validate all inputs

### Development Mode
- Server allows all origins by default
- No authentication required
- Debug logging available

## Integration Examples

### VS Code Workspace Settings
```json
{
  "mcp.servers": {
    "financial-housing": {
      "url": "ws://localhost:8080/mcp",
      "name": "Financial & Housing Data"
    }
  }
}
```

### Custom Copilot Extension
```typescript
import { MCPClient } from 'mcp-client';

const client = new MCPClient('ws://localhost:8080/mcp');
await client.initialize({
  protocolVersion: "2024-11-05",
  capabilities: {},
  clientInfo: { name: "copilot-extension", version: "1.0.0" }
});

const tools = await client.listTools();
// Use tools in your Copilot extension
```

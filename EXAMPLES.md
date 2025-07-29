# MCP Server Plugin Examples

## Starting the Server

```bash
# Quick start
make run

# Or build and run manually
make build
./bin/mcp-server

# With debug logging
./bin/mcp-server -debug

# On different port
./bin/mcp-server -port 9090
```

## WebSocket MCP Client Example

Connect to `ws://localhost:8080/mcp` and follow the MCP protocol:

### 1. Initialize Connection

```json
{
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
}
```

### 2. List Available Tools

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}
```

### 3. Financial Plugin Examples

#### Get Stock Data
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "get_stock_data",
    "arguments": {
      "symbol": "AAPL"
    }
  }
}
```

#### Search Companies
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "search_companies",
    "arguments": {
      "query": "Apple"
    }
  }
}
```

#### Get Market Summary
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "get_market_summary",
    "arguments": {}
  }
}
```

#### Get Historical Data
```json
{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "tools/call",
  "params": {
    "name": "get_historical_data",
    "arguments": {
      "symbol": "GOOGL",
      "period": "1mo"
    }
  }
}
```

### 4. Housing Plugin Examples

#### Search Properties
```json
{
  "jsonrpc": "2.0",
  "id": 7,
  "method": "tools/call",
  "params": {
    "name": "search_properties",
    "arguments": {
      "city": "San Francisco",
      "state": "CA",
      "minPrice": 500000,
      "maxPrice": 2000000,
      "bedrooms": 3,
      "propertyType": "house"
    }
  }
}
```

#### Get Property Details
```json
{
  "jsonrpc": "2.0",
  "id": 8,
  "method": "tools/call",
  "params": {
    "name": "get_property_details",
    "arguments": {
      "propertyId": "prop_001"
    }
  }
}
```

#### Get Market Statistics
```json
{
  "jsonrpc": "2.0",
  "id": 9,
  "method": "tools/call",
  "params": {
    "name": "get_market_stats",
    "arguments": {
      "city": "Los Angeles",
      "state": "CA"
    }
  }
}
```

#### Get Price History
```json
{
  "jsonrpc": "2.0",
  "id": 10,
  "method": "tools/call",
  "params": {
    "name": "get_price_history",
    "arguments": {
      "address": "123 Main St",
      "city": "San Francisco",
      "state": "CA"
    }
  }
}
```

#### Estimate Property Value
```json
{
  "jsonrpc": "2.0",
  "id": 11,
  "method": "tools/call",
  "params": {
    "name": "estimate_property_value",
    "arguments": {
      "address": "456 Oak Ave",
      "city": "Los Angeles",
      "state": "CA"
    }
  }
}
```

## HTTP REST API Examples

### Health Check
```bash
curl http://localhost:8080/health
```

### List Tools
```bash
curl http://localhost:8080/tools | jq
```

### List Resources
```bash
curl http://localhost:8080/resources | jq
```

## Testing with wscat

Install wscat: `npm install -g wscat`

```bash
# Connect to WebSocket
wscat -c ws://localhost:8080/mcp

# Then send JSON messages from the examples above
```

## Integration with Claude Desktop

Add to your Claude Desktop MCP configuration:

```json
{
  "mcpServers": {
    "play-mcp": {
      "command": "/path/to/play-mcp/bin/mcp-server",
      "args": ["-host", "localhost", "-port", "8080"]
    }
  }
}
```

## Available Tools Summary

### Financial Plugin (5 tools)
- `get_stock_data` - Current stock information
- `search_companies` - Company search
- `get_market_summary` - Market indices and movers
- `get_historical_data` - Historical price data

### Housing Plugin (5 tools)  
- `search_properties` - Property search with filters
- `get_property_details` - Detailed property information
- `get_market_stats` - Area market statistics
- `get_price_history` - Property price history
- `estimate_property_value` - Property value estimation

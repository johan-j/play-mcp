# Play MCP - Multi-Plugin MCP Server

A Go-based Model Context Protocol (MCP) server with plugin architecture supporting financial data and housing market data, designed for GitHub Copilot integration.

## Features

- **Plugin Architecture**: Extensible design for adding new data sources
- **Financial Data Plugin**: Stock prices, market data, company information
- **Housing Data Plugin**: Property listings, market statistics, price estimates
- **Chi Router**: Fast HTTP router with middleware support
- **WebSocket Support**: Real-time MCP protocol communication
- **REST API**: HTTP endpoints for easy integration
- **GitHub Copilot Ready**: Optimized for Copilot tool integration
- **Docker Support**: Containerized deployment

## Plugins

### Financial Data Plugin
- Get stock data by symbol
- Search companies by name
- Market summary and indices
- Historical price data
- Mock data with realistic examples

### Housing Data Plugin  
- Search properties by location and criteria
- Get detailed property information
- Market statistics for areas
- Price history and estimates
- Property value estimation (Zestimate-like)

## Installation

### Prerequisites
- Go 1.21 or higher
- Git

### Build from Source

```bash
git clone https://github.com/johan-j/play-mcp.git
cd play-mcp
make build
```

### Using Docker

```bash
make docker-build
make docker-run
```

## Usage

### Connect to GitHub Copilot (Recommended)

1. **Quick Setup:**
   ```bash
   ./setup-copilot.sh
   ```

2. **Manual Setup:**
   - Start server: `./bin/mcp-server`
   - WebSocket endpoint: `ws://localhost:8080/mcp`
   - Configure your Copilot client to connect to the endpoint
   - See: `GITHUB_COPILOT_SETUP.md` for detailed instructions

### Start the Server Manually

```bash
# HTTP/WebSocket mode (localhost:8080)
make run

# With debug logging
make run-debug

# On different port
./bin/mcp-server -port 9090

# With custom host
./bin/mcp-server -host 0.0.0.0 -port 8080
```

### Available Endpoints

#### WebSocket (MCP Protocol)
- `ws://localhost:8080/mcp` - MCP WebSocket endpoint for GitHub Copilot

#### REST API
- `GET /health` - Health check
- `GET /tools` - List available tools
- `GET /resources` - List available resources

## MCP Tools

### Financial Tools
- `get_stock_data` - Get current stock information
- `search_companies` - Search for companies by name
- `get_market_summary` - Get market indices and top movers
- `get_historical_data` - Get historical price data

### Housing Tools
- `search_properties` - Search properties by criteria
- `get_property_details` - Get detailed property information
- `get_market_stats` - Get market statistics for area
- `get_price_history` - Get property price history
- `estimate_property_value` - Get property value estimate

## Example Usage

### Using the Financial Plugin

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "get_stock_data",
    "arguments": {
      "symbol": "AAPL"
    }
  }
}
```

### Using the Housing Plugin

```json
{
  "jsonrpc": "2.0", 
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "search_properties",
    "arguments": {
      "city": "San Francisco",
      "state": "CA",
      "minPrice": 500000,
      "maxPrice": 2000000,
      "bedrooms": 3
    }
  }
}
```

## Development

### Project Structure

```
play-mcp/
├── cmd/mcp-server/          # Main application entry point
├── internal/
│   ├── plugins/             # Plugin system
│   │   ├── financial/       # Financial data plugin  
│   │   └── housing/         # Housing data plugin
│   └── server/              # MCP server implementation
├── pkg/mcp/                 # MCP protocol types
├── config/                  # Configuration files
├── Dockerfile               # Docker configuration
└── Makefile                 # Build automation
```

### Adding New Plugins

1. Create a new directory under `internal/plugins/`
2. Implement the `mcp.Plugin` interface:
   - `Name() string`
   - `Description() string` 
   - `GetTools() []mcp.Tool`
   - `HandleToolCall(ctx, request) (*mcp.ToolCallResponse, error)`
   - `GetResources() []mcp.Resource`
3. Register the plugin in `cmd/mcp-server/main.go`

### Development Commands

```bash
# Install dependencies
make deps

# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Clean build artifacts
make clean
```

## Configuration

Edit `config/config.yaml` to configure:
- Server settings (host, port)
- Plugin configuration
- API keys for real data sources
- Logging settings

## Real Data Integration

The plugins currently use mock data. To integrate with real APIs:

### Financial Data
- Alpha Vantage API
- Yahoo Finance API
- IEX Cloud API
- Polygon.io API

### Housing Data
- Zillow API
- Redfin API
- Realtor.com API
- RentSpree API

Add API keys to `config/config.yaml` and update plugin implementations.

## API Documentation

### MCP Protocol Flow

1. **Initialize** - Client establishes connection
2. **List Tools** - Get available tools
3. **Call Tools** - Execute specific tools
4. **Get Resources** - Access plugin resources

### Error Codes

- `-32601` - Method not found
- `-32602` - Invalid parameters  
- `-32603` - Internal error
- `-32002` - Client not initialized

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Run `make test` and `make lint`
5. Submit a pull request

## License

MIT License - see LICENSE file for details
# Play MCP - Financial Data Plugin

A Go-based Model Context Protocol (MCP) plugin for searching financial data.

## Features

- MCP server implementation in Go
- Financial data search capabilities
- RESTful API endpoints for financial queries
- Mock financial data provider (can be extended with real APIs)

## Installation

```bash
git clone https://github.com/johan-j/play-mcp.git
cd play-mcp
go mod tidy
go build -o bin/mcp-server ./cmd/mcp-server
```

## Usage

Run the MCP server:

```bash
./bin/mcp-server
```

The server will start on port 8080 and provide the following financial data search capabilities:

- Stock price lookups
- Company information search
- Market data queries
- Historical price data (mock implementation)

## API Endpoints

- `GET /api/stock/{symbol}` - Get stock information
- `GET /api/search/company?q={query}` - Search for companies
- `GET /api/market/summary` - Get market summary
- `GET /api/history/{symbol}?period={period}` - Get historical data

## Development

This project follows the Model Context Protocol specification and can be extended to integrate with real financial data providers like Alpha Vantage, Yahoo Finance, or Bloomberg APIs.

## License

MIT License
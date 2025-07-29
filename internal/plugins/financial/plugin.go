package financial

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/johan-j/play-mcp/pkg/mcp"
)

// Plugin implements the MCP plugin interface for financial data
type Plugin struct{}

// StockData represents stock information
type StockData struct {
	Symbol        string  `json:"symbol"`
	CompanyName   string  `json:"companyName"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
	Volume        int64   `json:"volume"`
	MarketCap     int64   `json:"marketCap"`
	PE            float64 `json:"pe"`
	Timestamp     string  `json:"timestamp"`
}

// MarketSummary represents market summary data
type MarketSummary struct {
	Indices    []IndexData `json:"indices"`
	TopGainers []StockData `json:"topGainers"`
	TopLosers  []StockData `json:"topLosers"`
	Timestamp  string      `json:"timestamp"`
}

// IndexData represents stock index information
type IndexData struct {
	Name          string  `json:"name"`
	Value         float64 `json:"value"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
}

// NewPlugin creates a new financial plugin instance
func NewPlugin() *Plugin {
	return &Plugin{}
}

// Name returns the plugin name
func (p *Plugin) Name() string {
	return "financial"
}

// Description returns the plugin description
func (p *Plugin) Description() string {
	return "Provides financial market data including stock prices, company information, and market summaries"
}

// GetTools returns the available tools for this plugin
func (p *Plugin) GetTools() []mcp.Tool {
	return []mcp.Tool{
		{
			Name:        "get_stock_data",
			Description: "Get current stock price and company information for a given symbol",
			InputSchema: mcp.ToolSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "Stock symbol (e.g., AAPL, GOOGL, MSFT)",
					},
				},
				Required: []string{"symbol"},
			},
		},
		{
			Name:        "search_companies",
			Description: "Search for companies by name or partial name",
			InputSchema: mcp.ToolSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Company name or partial name to search for",
					},
				},
				Required: []string{"query"},
			},
		},
		{
			Name:        "get_market_summary",
			Description: "Get current market summary including major indices and top movers",
			InputSchema: mcp.ToolSchema{
				Type:       "object",
				Properties: map[string]interface{}{},
			},
		},
		{
			Name:        "get_historical_data",
			Description: "Get historical price data for a stock symbol",
			InputSchema: mcp.ToolSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "Stock symbol (e.g., AAPL, GOOGL, MSFT)",
					},
					"period": map[string]interface{}{
						"type":        "string",
						"description": "Time period (1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max)",
						"default":     "1mo",
					},
				},
				Required: []string{"symbol"},
			},
		},
	}
}

// HandleToolCall handles tool calls for this plugin
func (p *Plugin) HandleToolCall(ctx context.Context, request mcp.ToolCallRequest) (*mcp.ToolCallResponse, error) {
	switch request.Name {
	case "get_stock_data":
		return p.handleGetStockData(ctx, request.Arguments)
	case "search_companies":
		return p.handleSearchCompanies(ctx, request.Arguments)
	case "get_market_summary":
		return p.handleGetMarketSummary(ctx, request.Arguments)
	case "get_historical_data":
		return p.handleGetHistoricalData(ctx, request.Arguments)
	default:
		return nil, fmt.Errorf("unknown tool: %s", request.Name)
	}
}

// GetResources returns available resources for this plugin
func (p *Plugin) GetResources() []mcp.Resource {
	return []mcp.Resource{
		{
			URI:         "financial://stocks",
			Name:        "Stock Data",
			Description: "Real-time and historical stock market data",
			MimeType:    "application/json",
		},
		{
			URI:         "financial://market",
			Name:        "Market Data",
			Description: "Market indices and summary information",
			MimeType:    "application/json",
		},
	}
}

func (p *Plugin) handleGetStockData(ctx context.Context, args map[string]interface{}) (*mcp.ToolCallResponse, error) {
	symbol, ok := args["symbol"].(string)
	if !ok {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: "symbol parameter is required and must be a string"}},
		}, nil
	}

	// Mock data - in real implementation, you would call a financial API
	stockData := p.getMockStockData(strings.ToUpper(symbol))

	data, err := json.MarshalIndent(stockData, "", "  ")
	if err != nil {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: fmt.Sprintf("Error marshaling stock data: %v", err)}},
		}, nil
	}

	return &mcp.ToolCallResponse{
		Content: []mcp.Content{
			{Type: "text", Text: fmt.Sprintf("Stock data for %s:", symbol)},
			{Type: "text", Text: string(data)},
		},
	}, nil
}

func (p *Plugin) handleSearchCompanies(ctx context.Context, args map[string]interface{}) (*mcp.ToolCallResponse, error) {
	query, ok := args["query"].(string)
	if !ok {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: "query parameter is required and must be a string"}},
		}, nil
	}

	// Mock search results
	companies := p.searchMockCompanies(query)

	data, err := json.MarshalIndent(companies, "", "  ")
	if err != nil {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: fmt.Sprintf("Error marshaling company data: %v", err)}},
		}, nil
	}

	return &mcp.ToolCallResponse{
		Content: []mcp.Content{
			{Type: "text", Text: fmt.Sprintf("Companies matching '%s':", query)},
			{Type: "text", Text: string(data)},
		},
	}, nil
}

func (p *Plugin) handleGetMarketSummary(ctx context.Context, args map[string]interface{}) (*mcp.ToolCallResponse, error) {
	summary := p.getMockMarketSummary()

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: fmt.Sprintf("Error marshaling market summary: %v", err)}},
		}, nil
	}

	return &mcp.ToolCallResponse{
		Content: []mcp.Content{
			{Type: "text", Text: "Current market summary:"},
			{Type: "text", Text: string(data)},
		},
	}, nil
}

func (p *Plugin) handleGetHistoricalData(ctx context.Context, args map[string]interface{}) (*mcp.ToolCallResponse, error) {
	symbol, ok := args["symbol"].(string)
	if !ok {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: "symbol parameter is required and must be a string"}},
		}, nil
	}

	period := "1mo"
	if p, ok := args["period"].(string); ok {
		period = p
	}

	// Mock historical data
	historicalData := p.getMockHistoricalData(strings.ToUpper(symbol), period)

	data, err := json.MarshalIndent(historicalData, "", "  ")
	if err != nil {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: fmt.Sprintf("Error marshaling historical data: %v", err)}},
		}, nil
	}

	return &mcp.ToolCallResponse{
		Content: []mcp.Content{
			{Type: "text", Text: fmt.Sprintf("Historical data for %s (%s):", symbol, period)},
			{Type: "text", Text: string(data)},
		},
	}, nil
}

// Mock data functions (replace with real API calls in production)

func (p *Plugin) getMockStockData(symbol string) StockData {
	mockData := map[string]StockData{
		"AAPL": {
			Symbol: "AAPL", CompanyName: "Apple Inc.", Price: 193.50, Change: 2.30, ChangePercent: 1.20,
			Volume: 45000000, MarketCap: 3000000000000, PE: 25.4, Timestamp: time.Now().Format(time.RFC3339),
		},
		"GOOGL": {
			Symbol: "GOOGL", CompanyName: "Alphabet Inc.", Price: 140.20, Change: -1.50, ChangePercent: -1.06,
			Volume: 25000000, MarketCap: 1800000000000, PE: 22.1, Timestamp: time.Now().Format(time.RFC3339),
		},
		"MSFT": {
			Symbol: "MSFT", CompanyName: "Microsoft Corporation", Price: 410.80, Change: 5.20, ChangePercent: 1.28,
			Volume: 35000000, MarketCap: 3100000000000, PE: 28.7, Timestamp: time.Now().Format(time.RFC3339),
		},
	}

	if data, exists := mockData[symbol]; exists {
		return data
	}

	// Return generic data for unknown symbols
	return StockData{
		Symbol: symbol, CompanyName: symbol + " Corporation", Price: 100.00, Change: 0.50, ChangePercent: 0.50,
		Volume: 1000000, MarketCap: 50000000000, PE: 20.0, Timestamp: time.Now().Format(time.RFC3339),
	}
}

func (p *Plugin) searchMockCompanies(query string) []StockData {
	allCompanies := []StockData{
		{Symbol: "AAPL", CompanyName: "Apple Inc.", Price: 193.50},
		{Symbol: "GOOGL", CompanyName: "Alphabet Inc.", Price: 140.20},
		{Symbol: "MSFT", CompanyName: "Microsoft Corporation", Price: 410.80},
		{Symbol: "AMZN", CompanyName: "Amazon.com Inc.", Price: 145.30},
		{Symbol: "TSLA", CompanyName: "Tesla Inc.", Price: 248.90},
		{Symbol: "META", CompanyName: "Meta Platforms Inc.", Price: 485.20},
		{Symbol: "NVDA", CompanyName: "NVIDIA Corporation", Price: 875.60},
	}

	var results []StockData
	queryLower := strings.ToLower(query)

	for _, company := range allCompanies {
		if strings.Contains(strings.ToLower(company.CompanyName), queryLower) ||
			strings.Contains(strings.ToLower(company.Symbol), queryLower) {
			results = append(results, company)
		}
	}

	return results
}

func (p *Plugin) getMockMarketSummary() MarketSummary {
	return MarketSummary{
		Indices: []IndexData{
			{Name: "S&P 500", Value: 5200.45, Change: 25.30, ChangePercent: 0.49},
			{Name: "NASDAQ", Value: 16800.20, Change: -45.80, ChangePercent: -0.27},
			{Name: "Dow Jones", Value: 38500.10, Change: 150.60, ChangePercent: 0.39},
		},
		TopGainers: []StockData{
			{Symbol: "XYZ", CompanyName: "XYZ Corp", Price: 25.50, Change: 5.20, ChangePercent: 25.60},
			{Symbol: "ABC", CompanyName: "ABC Inc", Price: 45.80, Change: 8.30, ChangePercent: 22.10},
		},
		TopLosers: []StockData{
			{Symbol: "DEF", CompanyName: "DEF Ltd", Price: 15.20, Change: -3.80, ChangePercent: -20.00},
			{Symbol: "GHI", CompanyName: "GHI Corp", Price: 32.10, Change: -6.20, ChangePercent: -16.20},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func (p *Plugin) getMockHistoricalData(symbol, period string) map[string]interface{} {
	// Generate mock historical data points
	dataPoints := make([]map[string]interface{}, 0)

	basePrice := 100.0
	if symbol == "AAPL" {
		basePrice = 190.0
	} else if symbol == "GOOGL" {
		basePrice = 135.0
	} else if symbol == "MSFT" {
		basePrice = 400.0
	}

	// Generate 30 days of mock data
	for i := 30; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		variation := (float64(i%10) - 5) * 2.0 // Simple variation
		price := basePrice + variation

		dataPoints = append(dataPoints, map[string]interface{}{
			"date":   date.Format("2006-01-02"),
			"open":   price - 1.0,
			"high":   price + 2.0,
			"low":    price - 2.5,
			"close":  price,
			"volume": 1000000 + (i * 50000),
		})
	}

	return map[string]interface{}{
		"symbol": symbol,
		"period": period,
		"data":   dataPoints,
	}
}

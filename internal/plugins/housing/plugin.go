package housing

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/johan-j/play-mcp/pkg/mcp"
	"github.com/johan-j/play-mcp/pkg/scraper"
)

// Plugin implements the MCP plugin interface for housing data
type Plugin struct{}

// PropertyData represents comprehensive property information for pricing analysis
type PropertyData struct {
	ID           string   `json:"id"`
	Address      string   `json:"address"`
	City         string   `json:"city"`
	State        string   `json:"state"`
	ZipCode      string   `json:"zipCode"`
	Price        int64    `json:"price"`
	Bedrooms     int      `json:"bedrooms"`
	Bathrooms    float64  `json:"bathrooms"`
	SquareFeet   int      `json:"squareFeet"`
	LotSize      float64  `json:"lotSize"`
	YearBuilt    int      `json:"yearBuilt"`
	PropertyType string   `json:"propertyType"`
	Status       string   `json:"status"`
	ListedDate   string   `json:"listedDate"`
	Description  string   `json:"description"`
	Features     []string `json:"features"`
	Images       []string `json:"images"`
	// Enhanced fields for comprehensive pricing analysis
	PropertyCondition  string   `json:"propertyCondition"`
	SchoolDistrict     string   `json:"schoolDistrict"`
	ElementarySchool   string   `json:"elementarySchool"`
	MiddleSchool       string   `json:"middleSchool"`
	HighSchool         string   `json:"highSchool"`
	SchoolRatings      []string `json:"schoolRatings"`
	Neighborhood       string   `json:"neighborhood"`
	PropertyTax        string   `json:"propertyTax"`
	HOAFees            string   `json:"hoaFees"`
	ParkingSpaces      int      `json:"parkingSpaces"`
	Garage             string   `json:"garage"`
	Heating            string   `json:"heating"`
	Cooling            string   `json:"cooling"`
	Flooring           []string `json:"flooring"`
	Appliances         []string `json:"appliances"`
	LastRenovated      string   `json:"lastRenovated"`
	PriceHistory       []string `json:"priceHistory"`
	NearbyComparables  []string `json:"nearbyComparables"`
	WalkScore          string   `json:"walkScore"`
	TransitScore       string   `json:"transitScore"`
	DistanceToBeach    string   `json:"distanceToBeach"`
	DistanceToDowntown string   `json:"distanceToDowntown"`
	FloodZone          string   `json:"floodZone"`
	HomeInsurance      string   `json:"homeInsurance"`
	SoldDate           string   `json:"soldDate"`
	DaysOnMarket       int      `json:"daysOnMarket"`
	PricePerSqFt       int      `json:"pricePerSqFt"`
	Agent              string   `json:"agent"`
	Brokerage          string   `json:"brokerage"`
}

// SearchFilters represents search criteria for properties
type SearchFilters struct {
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
}

// NewPlugin creates a new housing plugin instance
func NewPlugin() *Plugin {
	return &Plugin{}
}

// Name returns the plugin name
func (p *Plugin) Name() string {
	return "housing"
}

// Description returns the plugin description
func (p *Plugin) Description() string {
	return "Search for sold properties by neighborhood, city, and state"
}

// GetTools returns the available tools for this plugin
func (p *Plugin) GetTools() []mcp.Tool {
	return []mcp.Tool{
		{
			Name:        "search_sold_properties",
			Description: "Search for sold properties based on location",
			InputSchema: mcp.ToolSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"city": map[string]interface{}{
						"type":        "string",
						"description": "City name",
					},
					"state": map[string]interface{}{
						"type":        "string",
						"description": "State abbreviation (e.g., CA, NY, TX)",
					},
					"neighborhood": map[string]interface{}{
						"type":        "string",
						"description": "Neighborhood name (e.g., Manoa, Waikiki)",
					},
				},
				Required: []string{"city", "state", "neighborhood"},
			},
		},
		{
			Name:        "fetch_property_detail",
			Description: "Fetch detailed information about a specific property from homes.com or redfin.com",
			InputSchema: mcp.ToolSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "Full property URL from homes.com or redfin.com (e.g., https://www.homes.com/property/2819-poelua-st-honolulu-hi/n207sqkl8vl1p/ or https://www.redfin.com/HI/Honolulu/2819-Poelua-St-96822/home/88513618)",
					},
				},
				Required: []string{"url"},
			},
		},
	}
}

// HandleToolCall handles tool calls for this plugin
func (p *Plugin) HandleToolCall(_ context.Context, request mcp.ToolCallRequest) (*mcp.ToolCallResponse, error) {
	switch request.Name {
	case "search_sold_properties":
		return p.handleSearchSoldProperties(request.Arguments)
	case "fetch_property_detail":
		return p.handleFetchPropertyDetail(request.Arguments)
	default:
		return nil, fmt.Errorf("unknown tool: %s", request.Name)
	}
}

// GetResources returns available resources for this plugin
func (p *Plugin) GetResources() []mcp.Resource {
	return []mcp.Resource{
		{
			URI:         "housing://sold-properties",
			Name:        "Sold Properties",
			Description: "Recently sold real estate properties by location",
			MimeType:    "application/json",
		},
	}
}

func (p *Plugin) handleSearchSoldProperties(args map[string]interface{}) (*mcp.ToolCallResponse, error) {
	// Validate required parameters
	city, ok := args["city"].(string)
	if !ok || city == "" {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: "city parameter is required and must be a non-empty string"}},
		}, nil
	}

	state, ok := args["state"].(string)
	if !ok || state == "" {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: "state parameter is required and must be a non-empty string"}},
		}, nil
	}

	neighborhood, ok := args["neighborhood"].(string)
	if !ok || neighborhood == "" {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: "neighborhood parameter is required and must be a non-empty string"}},
		}, nil
	}

	filters := SearchFilters{
		City:         city,
		State:        state,
		Neighborhood: neighborhood,
	}

	// Check if this is a Manoa/Honolulu search and use scraper
	if (strings.ToLower(filters.City) == "honolulu" && strings.ToLower(filters.State) == "hi") ||
		strings.Contains(strings.ToLower(filters.City), "manoa") ||
		strings.ToLower(filters.Neighborhood) == "manoa" {
		properties, err := p.searchRealProperties(filters)
		if err == nil && len(properties) > 0 {
			data, err := json.MarshalIndent(properties, "", "  ")
			if err != nil {
				return &mcp.ToolCallResponse{
					IsError: true,
					Content: []mcp.Content{{Type: "text", Text: fmt.Sprintf("Error marshaling property data: %v", err)}},
				}, nil
			}

			return &mcp.ToolCallResponse{
				Content: []mcp.Content{
					{Type: "text", Text: fmt.Sprintf("Found %d sold properties matching your criteria:", len(properties))},
					{Type: "text", Text: string(data)},
				},
			}, nil
		}
	}

	// Fall back to mock properties for other locations
	properties := p.searchMockProperties(filters)

	data, err := json.MarshalIndent(properties, "", "  ")
	if err != nil {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: fmt.Sprintf("Error marshaling property data: %v", err)}},
		}, nil
	}

	return &mcp.ToolCallResponse{
		Content: []mcp.Content{
			{Type: "text", Text: fmt.Sprintf("Found %d sold properties matching your criteria:", len(properties))},
			{Type: "text", Text: string(data)},
		},
	}, nil
}

func (p *Plugin) handleFetchPropertyDetail(args map[string]interface{}) (*mcp.ToolCallResponse, error) {
	// Validate required parameters
	url, ok := args["url"].(string)
	if !ok || url == "" {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: "url parameter is required and must be a non-empty string"}},
		}, nil
	}

	// Validate URL format
	if !strings.Contains(url, "homes.com/property/") && !strings.Contains(url, "redfin.com/") {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: "URL must be a valid homes.com or redfin.com property URL"}},
		}, nil
	}

	homesScraper := scraper.NewHomesScraper()
	propertyDetail, err := homesScraper.ScrapePropertyDetail(url)
	if err != nil {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: fmt.Sprintf("Error fetching property details: %v", err)}},
		}, nil
	}

	// Convert scraper property to our comprehensive PropertyData format
	property := PropertyData{
		ID:           propertyDetail.ID,
		Address:      propertyDetail.Address,
		City:         propertyDetail.City,
		State:        propertyDetail.State,
		ZipCode:      propertyDetail.ZipCode,
		Price:        int64(propertyDetail.Price),
		Bedrooms:     propertyDetail.Bedrooms,
		Bathrooms:    propertyDetail.Bathrooms,
		SquareFeet:   propertyDetail.SquareFeet,
		LotSize:      propertyDetail.LotSize,
		PropertyType: propertyDetail.PropertyType,
		Status:       propertyDetail.Status,
		YearBuilt:    propertyDetail.YearBuilt,
		Description:  propertyDetail.Description,
		Features:     propertyDetail.Features,
		ListedDate:   propertyDetail.SoldDate,
		SoldDate:     propertyDetail.SoldDate,
		DaysOnMarket: propertyDetail.DaysOnMarket,
		PricePerSqFt: propertyDetail.PricePerSqFt,
		Agent:        propertyDetail.Agent,
		Brokerage:    propertyDetail.Brokerage,
		// Enhanced fields for pricing analysis
		PropertyCondition:  propertyDetail.PropertyCondition,
		SchoolDistrict:     propertyDetail.SchoolDistrict,
		ElementarySchool:   propertyDetail.ElementarySchool,
		MiddleSchool:       propertyDetail.MiddleSchool,
		HighSchool:         propertyDetail.HighSchool,
		SchoolRatings:      propertyDetail.SchoolRatings,
		Neighborhood:       propertyDetail.Neighborhood,
		PropertyTax:        propertyDetail.PropertyTax,
		HOAFees:            propertyDetail.HOAFees,
		ParkingSpaces:      propertyDetail.ParkingSpaces,
		Garage:             propertyDetail.Garage,
		Heating:            propertyDetail.Heating,
		Cooling:            propertyDetail.Cooling,
		Flooring:           propertyDetail.Flooring,
		Appliances:         propertyDetail.Appliances,
		LastRenovated:      propertyDetail.LastRenovated,
		PriceHistory:       propertyDetail.PriceHistory,
		NearbyComparables:  propertyDetail.NearbyComparables,
		WalkScore:          propertyDetail.WalkScore,
		TransitScore:       propertyDetail.TransitScore,
		DistanceToBeach:    propertyDetail.DistanceToBeach,
		DistanceToDowntown: propertyDetail.DistanceToDowntown,
		FloodZone:          propertyDetail.FloodZone,
		HomeInsurance:      propertyDetail.HomeInsurance,
	}

	data, err := json.MarshalIndent(property, "", "  ")
	if err != nil {
		return &mcp.ToolCallResponse{
			IsError: true,
			Content: []mcp.Content{{Type: "text", Text: fmt.Sprintf("Error marshaling property data: %v", err)}},
		}, nil
	}

	return &mcp.ToolCallResponse{
		Content: []mcp.Content{
			{Type: "text", Text: "Property details:"},
			{Type: "text", Text: string(data)},
		},
	}, nil
}

// Mock data functions (replace with real API calls in production)

func (p *Plugin) searchMockProperties(filters SearchFilters) []PropertyData {
	mockProperties := []PropertyData{
		{
			ID: "prop_001", Address: "123 Main St", City: "San Francisco", State: "CA", ZipCode: "94102",
			Price: 1250000, Bedrooms: 3, Bathrooms: 2.5, SquareFeet: 1800, LotSize: 0.12,
			YearBuilt: 1925, PropertyType: "house", Status: "sold", ListedDate: "2025-01-15",
			Description: "Beautiful Victorian home in prime location",
			Features:    []string{"fireplace", "hardwood floors", "updated kitchen"},
		},
		{
			ID: "prop_002", Address: "456 Oak Ave", City: "Los Angeles", State: "CA", ZipCode: "90210",
			Price: 850000, Bedrooms: 2, Bathrooms: 2.0, SquareFeet: 1200, LotSize: 0.08,
			YearBuilt: 1985, PropertyType: "condo", Status: "sold", ListedDate: "2025-01-20",
			Description: "Modern condo with city views",
			Features:    []string{"balcony", "parking", "pool", "gym"},
		},
		{
			ID: "prop_003", Address: "789 Pine St", City: "Seattle", State: "WA", ZipCode: "98101",
			Price: 675000, Bedrooms: 4, Bathrooms: 3.0, SquareFeet: 2200, LotSize: 0.15,
			YearBuilt: 2010, PropertyType: "house", Status: "sold", ListedDate: "2025-01-10",
			Description: "Contemporary home with mountain views",
			Features:    []string{"deck", "garage", "open floor plan"},
		},
	}

	// Simple filtering logic
	var results []PropertyData
	for _, prop := range mockProperties {
		matches := true

		if filters.City != "" && !strings.EqualFold(prop.City, filters.City) {
			matches = false
		}
		if filters.State != "" && !strings.EqualFold(prop.State, filters.State) {
			matches = false
		}
		if filters.Neighborhood != "" && !strings.EqualFold(prop.City, filters.Neighborhood) {
			matches = false
		}

		if matches {
			results = append(results, prop)
		}
	}

	return results
}

// searchRealProperties uses the scraper to get real property data
func (p *Plugin) searchRealProperties(filters SearchFilters) ([]PropertyData, error) {
	homesScraper := scraper.NewHomesScraper()

	var scrapedProperties []scraper.Property
	var err error

	// Use neighborhood-specific scraping if neighborhood is specified
	if filters.Neighborhood != "" {
		city := filters.City
		if city == "" {
			city = "Honolulu" // Default to Honolulu for Hawaii neighborhoods
		}
		state := filters.State
		if state == "" {
			state = "HI" // Default to Hawaii
		}
		scrapedProperties, err = homesScraper.ScrapeNeighborhood(city, state, filters.Neighborhood, "sold")
	} else {
		// Fall back to Manoa-specific scraping
		scrapedProperties, err = homesScraper.ScrapeManoa("sold")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scrape properties: %v", err)
	}

	var properties []PropertyData

	// Convert scraper properties to our PropertyData format
	for _, scraped := range scrapedProperties {
		property := PropertyData{
			ID:           scraped.ID,
			Address:      scraped.Address,
			City:         scraped.City,
			State:        scraped.State,
			ZipCode:      scraped.ZipCode,
			Price:        int64(scraped.Price),
			Bedrooms:     scraped.Bedrooms,
			Bathrooms:    scraped.Bathrooms,
			SquareFeet:   scraped.SquareFeet,
			PropertyType: scraped.PropertyType,
			Status:       scraped.Status,
			YearBuilt:    scraped.YearBuilt,
			Description:  scraped.Description,
			Features:     scraped.Features,
		}

		// Apply filters
		if p.matchesFilters(property, filters) {
			properties = append(properties, property)
		}
	}

	return properties, nil
}

// matchesFilters checks if a property matches the search filters
func (p *Plugin) matchesFilters(property PropertyData, filters SearchFilters) bool {
	// For now, since we removed the filtering fields, we'll just return true
	// In a real implementation, you might want to check city, state, or neighborhood
	return true
}

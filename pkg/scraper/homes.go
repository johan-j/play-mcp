package scraper

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// HomesScraper handles scraping data from homes.com
type HomesScraper struct {
	client *http.Client
}

// Property represents a property from homes.com or redfin.com
type Property struct {
	ID           string   `json:"id"`
	Address      string   `json:"address"`
	City         string   `json:"city"`
	State        string   `json:"state"`
	ZipCode      string   `json:"zipCode"`
	Price        int      `json:"price"`
	Bedrooms     int      `json:"bedrooms"`
	Bathrooms    float64  `json:"bathrooms"`
	SquareFeet   int      `json:"squareFeet"`
	LotSize      float64  `json:"lotSize"`
	PricePerSqFt int      `json:"pricePerSqFt"`
	YearBuilt    int      `json:"yearBuilt"`
	PropertyType string   `json:"propertyType"`
	Status       string   `json:"status"`
	SoldDate     string   `json:"soldDate"`
	DaysOnMarket int      `json:"daysOnMarket"`
	Description  string   `json:"description"`
	Features     []string `json:"features"`
	Agent        string   `json:"agent"`
	Brokerage    string   `json:"brokerage"`
	// Enhanced fields for pricing analysis
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
}

// MarketStats represents market statistics for an area
type MarketStats struct {
	Area                    string  `json:"area"`
	MedianSalePrice         int     `json:"medianSalePrice"`
	MedianSingleFamilyPrice int     `json:"medianSingleFamilyPrice"`
	MedianTownhousePrice    int     `json:"medianTownhousePrice"`
	AveragePricePerSqFt     int     `json:"averagePricePerSqFt"`
	HomesForSale            int     `json:"homesForSale"`
	SalesLast12Months       int     `json:"salesLast12Months"`
	AverageDaysOnMarket     int     `json:"averageDaysOnMarket"`
	MonthsOfSupply          float64 `json:"monthsOfSupply"`
	YearOverYearChange      float64 `json:"yearOverYearChange"`
	Timestamp               string  `json:"timestamp"`
}

// NewHomesScraper creates a new homes.com scraper
func NewHomesScraper() *HomesScraper {
	// Create a more robust HTTP transport that mimics browser behavior
	transport := &http.Transport{
		ForceAttemptHTTP2:     false, // Force HTTP/1.1 to avoid HTTP/2 stream errors
		DisableKeepAlives:     false,
		DisableCompression:    false,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		},
	}

	return &HomesScraper{
		client: &http.Client{
			Timeout:   45 * time.Second, // Increased timeout
			Transport: transport,
		},
	}
}

// ScrapeManoa scrapes properties from Manoa, Honolulu
func (h *HomesScraper) ScrapeManoa(status string) ([]Property, error) {
	return h.ScrapeNeighborhood("Honolulu", "HI", "manoa", status)
}

// ScrapeNeighborhood scrapes properties from a specific neighborhood
func (h *HomesScraper) ScrapeNeighborhood(city, state, neighborhood, status string) ([]Property, error) {
	var properties []Property
	propertyMap := make(map[string]Property) // Use map to deduplicate by ID

	// Build URL for specific neighborhood
	cityState := fmt.Sprintf("%s-%s", strings.ToLower(city), strings.ToLower(state))
	neighborhoodSlug := strings.ToLower(strings.ReplaceAll(neighborhood, " ", "-"))

	// Scrape multiple pages
	for page := 1; page <= 3; page++ {
		var url string
		if page == 1 {
			url = fmt.Sprintf("https://www.homes.com/%s/%s-neighborhood/%s/", cityState, neighborhoodSlug, status)
		} else {
			url = fmt.Sprintf("https://www.homes.com/%s/%s-neighborhood/%s/p%d/", cityState, neighborhoodSlug, status, page)
		}

		pageProperties, err := h.scrapePage(url)
		if err != nil {
			return nil, fmt.Errorf("failed to scrape page %d: %v", page, err)
		}

		log.Printf("Found %d properties on page %d", len(pageProperties), page)

		// Add properties to map to deduplicate by ID
		for _, prop := range pageProperties {
			if prop.ID != "" {
				propertyMap[prop.ID] = prop
			}
		}

		// If we got less than expected, we're probably at the end
		if len(pageProperties) < 20 {
			break
		}
	}

	// Convert map back to slice
	for _, prop := range propertyMap {
		properties = append(properties, prop)
	}

	log.Printf("Total properties found: %d (after deduplication)", len(properties))
	return properties, nil
}

// ScrapeNeighborhoodStats scrapes market statistics for a specific neighborhood
func (h *HomesScraper) ScrapeNeighborhoodStats(city, state, neighborhood string) (*MarketStats, error) {
	cityState := fmt.Sprintf("%s-%s", strings.ToLower(city), strings.ToLower(state))
	neighborhoodSlug := strings.ToLower(strings.ReplaceAll(neighborhood, " ", "-"))
	url := fmt.Sprintf("https://www.homes.com/%s/%s-neighborhood/sold/", cityState, neighborhoodSlug)

	log.Printf("Fetching market stats from homes.com URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers to mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	stats := &MarketStats{
		Area:      fmt.Sprintf("%s, %s, %s", neighborhood, city, state),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Extract market statistics from the page
	doc.Find(".housing-trends").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(j int, row *goquery.Selection) {
			label := strings.TrimSpace(row.Find("td:first-child").Text())
			value := strings.TrimSpace(row.Find("td:last-child").Text())

			switch {
			case strings.Contains(label, "Median Sale Price"):
				stats.MedianSalePrice = parsePrice(value)
			case strings.Contains(label, "Median Single Family Sale Price"):
				stats.MedianSingleFamilyPrice = parsePrice(value)
			case strings.Contains(label, "Median Townhouse Sale Price"):
				stats.MedianTownhousePrice = parsePrice(value)
			case strings.Contains(label, "Average Price Per Sq Ft"):
				stats.AveragePricePerSqFt = parsePrice(value)
			case strings.Contains(label, "Number of Homes for Sale"):
				stats.HomesForSale = parseInt(value)
			case strings.Contains(label, "Last 12 months Home Sales"):
				stats.SalesLast12Months = parseInt(value)
			case strings.Contains(label, "YoY Change"):
				stats.YearOverYearChange = parsePercent(value)
			}
		})
	})

	// If we didn't get data from table, try extracting from text
	if stats.MedianSalePrice == 0 {
		text := doc.Text()
		if match := regexp.MustCompile(`Median Sale Price \$([0-9,]+)`).FindStringSubmatch(text); len(match) > 1 {
			stats.MedianSalePrice = parsePrice(match[1])
		}
		if match := regexp.MustCompile(`Average Price Per Sq Ft \$([0-9,]+)`).FindStringSubmatch(text); len(match) > 1 {
			stats.AveragePricePerSqFt = parsePrice(match[1])
		}
		if match := regexp.MustCompile(`([0-9]+) days on the market`).FindStringSubmatch(text); len(match) > 1 {
			stats.AverageDaysOnMarket = parseInt(match[1])
		}
		if match := regexp.MustCompile(`down ([0-9]+)%`).FindStringSubmatch(text); len(match) > 1 {
			stats.YearOverYearChange = -parsePercent(match[1] + "%")
		}
		if match := regexp.MustCompile(`([0-9]+) homes for sale in`).FindStringSubmatch(text); len(match) > 1 {
			stats.HomesForSale = parseInt(match[1])
		}
		if match := regexp.MustCompile(`([0-9]+) home sales`).FindStringSubmatch(text); len(match) > 1 {
			stats.SalesLast12Months = parseInt(match[1])
		}
	}

	return stats, nil
}

// ScrapeMarketStats scrapes market statistics for an area
func (h *HomesScraper) ScrapeMarketStats(city, state string) (*MarketStats, error) {
	url := fmt.Sprintf("https://www.homes.com/%s-%s/manoa-neighborhood/sold/", strings.ToLower(city), strings.ToLower(state))

	log.Printf("Fetching market stats from homes.com URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers to mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	stats := &MarketStats{
		Area:      fmt.Sprintf("%s, %s", city, state),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Extract market statistics from the page
	doc.Find(".housing-trends").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(j int, row *goquery.Selection) {
			label := strings.TrimSpace(row.Find("td:first-child").Text())
			value := strings.TrimSpace(row.Find("td:last-child").Text())

			switch {
			case strings.Contains(label, "Median Sale Price"):
				stats.MedianSalePrice = parsePrice(value)
			case strings.Contains(label, "Median Single Family Sale Price"):
				stats.MedianSingleFamilyPrice = parsePrice(value)
			case strings.Contains(label, "Median Townhouse Sale Price"):
				stats.MedianTownhousePrice = parsePrice(value)
			case strings.Contains(label, "Average Price Per Sq Ft"):
				stats.AveragePricePerSqFt = parsePrice(value)
			case strings.Contains(label, "Number of Homes for Sale"):
				stats.HomesForSale = parseInt(value)
			case strings.Contains(label, "Last 12 months Home Sales"):
				stats.SalesLast12Months = parseInt(value)
			case strings.Contains(label, "YoY Change"):
				stats.YearOverYearChange = parsePercent(value)
			}
		})
	})

	// If we didn't get data from table, try extracting from text
	if stats.MedianSalePrice == 0 {
		text := doc.Text()
		if match := regexp.MustCompile(`Median Sale Price \$([0-9,]+)`).FindStringSubmatch(text); len(match) > 1 {
			stats.MedianSalePrice = parsePrice(match[1])
		}
		if match := regexp.MustCompile(`Average Price Per Sq Ft \$([0-9,]+)`).FindStringSubmatch(text); len(match) > 1 {
			stats.AveragePricePerSqFt = parsePrice(match[1])
		}
		if match := regexp.MustCompile(`([0-9]+) days on the market`).FindStringSubmatch(text); len(match) > 1 {
			stats.AverageDaysOnMarket = parseInt(match[1])
		}
		if match := regexp.MustCompile(`down ([0-9]+)%`).FindStringSubmatch(text); len(match) > 1 {
			stats.YearOverYearChange = -parsePercent(match[1] + "%")
		}
		if match := regexp.MustCompile(`([0-9]+) homes in Manoa`).FindStringSubmatch(text); len(match) > 1 {
			stats.HomesForSale = parseInt(match[1])
		}
		if match := regexp.MustCompile(`([0-9]+) home sales`).FindStringSubmatch(text); len(match) > 1 {
			stats.SalesLast12Months = parseInt(match[1])
		}
	}

	return stats, nil
}

// scrapePage scrapes a single page of properties
func (h *HomesScraper) scrapePage(url string) ([]Property, error) {
	log.Printf("Fetching properties from homes.com URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers to mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var properties []Property

	// Extract properties from the HTML content using text parsing
	// since homes.com has complex dynamic content
	pageText := doc.Text()

	// Find property blocks by looking for address patterns
	addressRegex := regexp.MustCompile(`(\d+\s+[A-Za-z\s]+(?:St|Ave|Rd|Dr|Pl|Place|Street|Avenue|Road|Drive).*?Honolulu,\s*HI\s*96822)`)
	priceRegex := regexp.MustCompile(`\$([0-9,]+)`)
	sqftRegex := regexp.MustCompile(`([0-9,]+)\s*Sq\s*Ft`)
	bedroomRegex := regexp.MustCompile(`(\d+)\s*Bed`)
	bathroomRegex := regexp.MustCompile(`([0-9.]+)\s*Bath`)
	soldRegex := regexp.MustCompile(`SOLD\s+([A-Z]{3}\s+\d{1,2},\s+\d{4})`)
	daysRegex := regexp.MustCompile(`(\d+)\s*Days\s*On\s*Market`)
	yearRegex := regexp.MustCompile(`Built\s+(\d{4})`)

	// Split text into chunks around addresses
	addresses := addressRegex.FindAllString(pageText, -1)

	for _, address := range addresses {
		if address == "" {
			continue
		}

		property := Property{
			Status: "sold",
		}

		// Parse address
		parts := strings.Split(address, ",")
		if len(parts) >= 3 {
			property.Address = strings.TrimSpace(parts[0])
			property.City = "Honolulu"
			property.State = "HI"
			property.ZipCode = "96822"
		}

		// Find the text chunk around this address
		addressIndex := strings.Index(pageText, address)
		if addressIndex == -1 {
			continue
		}

		start := addressIndex - 500
		if start < 0 {
			start = 0
		}
		end := addressIndex + 1000
		if end > len(pageText) {
			end = len(pageText)
		}

		chunk := pageText[start:end]

		// Extract price
		if priceMatch := priceRegex.FindStringSubmatch(chunk); len(priceMatch) > 1 {
			property.Price = parsePrice(priceMatch[1])
		}

		// Extract square feet
		if sqftMatch := sqftRegex.FindStringSubmatch(chunk); len(sqftMatch) > 1 {
			property.SquareFeet = parseInt(sqftMatch[1])
		}

		// Extract bedrooms
		if bedMatch := bedroomRegex.FindStringSubmatch(chunk); len(bedMatch) > 1 {
			property.Bedrooms = parseInt(bedMatch[1])
		}

		// Extract bathrooms
		if bathMatch := bathroomRegex.FindStringSubmatch(chunk); len(bathMatch) > 1 {
			property.Bathrooms = parseFloat(bathMatch[1])
		}

		// Extract sold date
		if soldMatch := soldRegex.FindStringSubmatch(chunk); len(soldMatch) > 1 {
			property.SoldDate = soldMatch[1]
		}

		// Extract days on market
		if daysMatch := daysRegex.FindStringSubmatch(chunk); len(daysMatch) > 1 {
			property.DaysOnMarket = parseInt(daysMatch[1])
		}

		// Extract year built
		if yearMatch := yearRegex.FindStringSubmatch(chunk); len(yearMatch) > 1 {
			property.YearBuilt = parseInt(yearMatch[1])
		}

		// Calculate price per sqft
		if property.Price > 0 && property.SquareFeet > 0 {
			property.PricePerSqFt = property.Price / property.SquareFeet
		}

		// Set property type
		property.PropertyType = "house"

		// Generate ID
		if property.Address != "" {
			property.ID = generatePropertyID(property.Address, property.City)
		}

		// Only add if we have basic info
		if property.Address != "" && property.Price > 0 {
			properties = append(properties, property)
		}
	}

	return properties, nil
}

// extractPropertyFromCard extracts property data from a property card element
func (h *HomesScraper) extractPropertyFromCard(s *goquery.Selection) Property {
	property := Property{}

	// Extract address
	addressText := strings.TrimSpace(s.Find("h3, .property-address, .listing-address").First().Text())
	if addressText != "" {
		parts := strings.Split(addressText, ",")
		if len(parts) >= 3 {
			property.Address = strings.TrimSpace(parts[0])
			property.City = strings.TrimSpace(parts[1])
			stateZip := strings.TrimSpace(parts[2])
			stateParts := strings.Fields(stateZip)
			if len(stateParts) >= 2 {
				property.State = stateParts[0]
				property.ZipCode = stateParts[1]
			}
		}
	}

	// Extract price
	priceText := s.Find(".price, .listing-price").First().Text()
	property.Price = parsePrice(priceText)

	// Extract bedrooms/bathrooms
	bedroomText := s.Find(".beds, .bedrooms").First().Text()
	property.Bedrooms = parseInt(bedroomText)

	bathroomText := s.Find(".baths, .bathrooms").First().Text()
	property.Bathrooms = parseFloat(bathroomText)

	// Extract square footage
	sqftText := s.Find(".sqft, .square-feet").First().Text()
	property.SquareFeet = parseInt(sqftText)

	// Calculate price per sqft if both are available
	if property.Price > 0 && property.SquareFeet > 0 {
		property.PricePerSqFt = property.Price / property.SquareFeet
	}

	// Extract year built
	yearText := s.Find(".year-built, .built").First().Text()
	property.YearBuilt = parseInt(yearText)

	// Extract days on market
	domText := s.Find(".days-on-market, .dom").First().Text()
	property.DaysOnMarket = parseInt(domText)

	// Extract description
	property.Description = strings.TrimSpace(s.Find(".description, .listing-description").First().Text())

	// Extract agent and brokerage
	property.Agent = strings.TrimSpace(s.Find(".agent-name, .listing-agent").First().Text())
	property.Brokerage = strings.TrimSpace(s.Find(".brokerage-name, .listing-brokerage").First().Text())

	// Set property type and status
	property.PropertyType = "house" // Default, could be enhanced
	property.Status = "for_sale"    // Default, could be enhanced

	// Generate an ID
	if property.Address != "" {
		property.ID = generatePropertyID(property.Address, property.City)
	}

	return property
}

// Helper functions for parsing
func parsePrice(s string) int {
	// Remove $ and commas, extract numbers
	re := regexp.MustCompile(`[0-9,]+`)
	match := re.FindString(s)
	if match == "" {
		return 0
	}

	cleaned := strings.ReplaceAll(match, ",", "")
	price, _ := strconv.Atoi(cleaned)
	return price
}

func parseInt(s string) int {
	re := regexp.MustCompile(`[0-9]+`)
	match := re.FindString(s)
	if match == "" {
		return 0
	}

	val, _ := strconv.Atoi(match)
	return val
}

func parseFloat(s string) float64 {
	re := regexp.MustCompile(`[0-9]+\.?[0-9]*`)
	match := re.FindString(s)
	if match == "" {
		return 0
	}

	val, _ := strconv.ParseFloat(match, 64)
	return val
}

func parsePercent(s string) float64 {
	re := regexp.MustCompile(`-?[0-9]+\.?[0-9]*`)
	match := re.FindString(s)
	if match == "" {
		return 0
	}

	val, _ := strconv.ParseFloat(match, 64)
	return val
}

func generatePropertyID(address, city string) string {
	// Generate a simple ID from address and city
	combined := strings.ToLower(address + city)
	cleaned := regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(combined, "")
	if len(cleaned) > 25 {
		cleaned = cleaned[:25]
	}
	if cleaned == "" {
		// Fallback to a simple hash-like string if cleaning results in empty string
		cleaned = "unknown"
	}
	return "prop_" + cleaned
}

// setBrowserHeaders sets comprehensive headers to mimic a real browser
func (h *HomesScraper) setBrowserHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Cache-Control", "max-age=0")
}

// ScrapePropertyDetail scrapes detailed information from a specific property page
func (h *HomesScraper) ScrapePropertyDetail(url string) (*Property, error) {
	log.Printf("Fetching property details from URL: %s", url)

	// Check if this is a Redfin URL and handle it differently
	if strings.Contains(url, "redfin.com") {
		return h.scrapeRedfinProperty(url)
	}

	// Original homes.com logic follows...
	// Extract address from URL first
	property := &Property{
		Status: "sold", // Default assumption for detail pages
	}

	// Parse address from URL
	urlParts := strings.Split(url, "/")
	if len(urlParts) >= 5 {
		addressSlug := urlParts[4]
		// Convert slug back to address format
		addressParts := strings.Split(addressSlug, "-")
		if len(addressParts) >= 3 {
			var addressComponents []string
			var city, state string

			// Look for state abbreviation (2 letters at the end)
			for i, part := range addressParts {
				if len(part) == 2 {
					state = strings.ToUpper(part)
					if i > 0 {
						city = strings.Title(addressParts[i-1])
					}
					// Everything before city is the address
					if i > 1 {
						for j := 0; j < i-1; j++ {
							if addressParts[j] != "" {
								addressComponents = append(addressComponents, strings.Title(addressParts[j]))
							}
						}
					}
					break
				}
			}

			if len(addressComponents) > 0 {
				property.Address = strings.Join(addressComponents, " ")
			}
			if city != "" {
				property.City = city
			}
			if state != "" {
				property.State = state
			}
		}
	}

	// If this is a Honolulu/Manoa property, try to find it in our existing sold data
	if strings.ToLower(property.City) == "honolulu" && strings.ToLower(property.State) == "hi" {
		log.Printf("Searching for property %s in existing sold data", property.Address)

		// Get all sold properties for Manoa
		soldProperties, err := h.ScrapeManoa("sold")
		if err == nil {
			// Search for matching address
			for _, soldProp := range soldProperties {
				if strings.Contains(strings.ToLower(soldProp.Address), strings.ToLower(property.Address)) {
					log.Printf("Found matching property in sold data: %s", soldProp.Address)
					return &soldProp, nil
				}
			}
		}
	}

	// If not found in sold data, try direct scraping (will likely fail with 403)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Use same simple headers as the working scrapePage function
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	resp, err := h.client.Do(req)
	if err != nil {
		// If direct scraping fails, return what we have from URL parsing
		log.Printf("Direct scraping failed, returning parsed URL data: %v", err)
		property.ID = generatePropertyID(property.Address, property.City)
		return property, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// If we get 403 or other error, return what we have from URL parsing
		log.Printf("HTTP %d: %s, returning parsed URL data", resp.StatusCode, resp.Status)
		property.ID = generatePropertyID(property.Address, property.City)
		return property, nil
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract property details from the page
	pageText := doc.Text()

	// Use regex patterns to extract detailed information from page content
	priceRegex := regexp.MustCompile(`\$([0-9,]+)`)
	sqftRegex := regexp.MustCompile(`([0-9,]+)\s*(?:sq|square)\s*(?:ft|feet)`)
	bedroomRegex := regexp.MustCompile(`(\d+)\s*(?:bed|bedroom)`)
	bathroomRegex := regexp.MustCompile(`([0-9.]+)\s*(?:bath|bathroom)`)
	yearRegex := regexp.MustCompile(`(?:built|year)\s*:?\s*(\d{4})`)
	soldRegex := regexp.MustCompile(`SOLD\s+([A-Z]{3}\s+\d{1,2},\s+\d{4})`)

	// Extract price
	if priceMatch := priceRegex.FindStringSubmatch(pageText); len(priceMatch) > 1 {
		property.Price = parsePrice(priceMatch[1])
	}

	// Extract square feet
	if sqftMatch := sqftRegex.FindStringSubmatch(pageText); len(sqftMatch) > 1 {
		property.SquareFeet = parseInt(sqftMatch[1])
	}

	// Extract bedrooms
	if bedMatch := bedroomRegex.FindStringSubmatch(pageText); len(bedMatch) > 1 {
		property.Bedrooms = parseInt(bedMatch[1])
	}

	// Extract bathrooms
	if bathMatch := bathroomRegex.FindStringSubmatch(pageText); len(bathMatch) > 1 {
		property.Bathrooms = parseFloat(bathMatch[1])
	}

	// Extract year built
	if yearMatch := yearRegex.FindStringSubmatch(pageText); len(yearMatch) > 1 {
		property.YearBuilt = parseInt(yearMatch[1])
	}

	// Extract sold date
	if soldMatch := soldRegex.FindStringSubmatch(pageText); len(soldMatch) > 1 {
		property.SoldDate = soldMatch[1]
	}

	// Extract description from meta tags or content
	doc.Find("meta[name='description']").Each(func(i int, s *goquery.Selection) {
		if content, exists := s.Attr("content"); exists {
			property.Description = content
		}
	})

	// Extract features from lists or bullet points
	var features []string
	doc.Find("li, .feature, .amenity").Each(func(i int, s *goquery.Selection) {
		feature := strings.TrimSpace(s.Text())
		if feature != "" && len(feature) < 100 { // Reasonable feature length
			features = append(features, feature)
		}
	})
	property.Features = features

	// Set property type
	property.PropertyType = "house"
	if strings.Contains(strings.ToLower(pageText), "condo") || strings.Contains(strings.ToLower(pageText), "condominium") {
		property.PropertyType = "condo"
	} else if strings.Contains(strings.ToLower(pageText), "townhouse") || strings.Contains(strings.ToLower(pageText), "townhome") {
		property.PropertyType = "townhouse"
	}

	// Generate ID
	if property.Address != "" {
		property.ID = generatePropertyID(property.Address, property.City)
	}

	// Set default zip code for Honolulu properties
	if property.City == "Honolulu" && property.State == "HI" {
		property.ZipCode = "96822"
	}

	// Calculate price per sqft
	if property.Price > 0 && property.SquareFeet > 0 {
		property.PricePerSqFt = property.Price / property.SquareFeet
	}

	return property, nil
}

// scrapeRedfinProperty scrapes property details from a Redfin URL
func (h *HomesScraper) scrapeRedfinProperty(url string) (*Property, error) {
	log.Printf("Fetching property details from Redfin URL: %s", url)

	property := &Property{
		Status: "sold", // Default assumption
	}

	// Parse address from Redfin URL format: /HI/Honolulu/2819-Poelua-St-96822/home/88513618
	urlParts := strings.Split(url, "/")
	for i, part := range urlParts {
		// Look for state (2 letters)
		if len(part) == 2 && strings.ToUpper(part) == part {
			property.State = part
			// Next part should be city
			if i+1 < len(urlParts) {
				property.City = urlParts[i+1]
			}
			// Next part should be address-zipcode
			if i+2 < len(urlParts) {
				addressPart := urlParts[i+2]
				// Split on last dash to separate address from zipcode
				lastDashIndex := strings.LastIndex(addressPart, "-")
				if lastDashIndex > 0 {
					address := addressPart[:lastDashIndex]
					zipCode := addressPart[lastDashIndex+1:]

					// Convert dashes to spaces and capitalize
					addressComponents := strings.Split(address, "-")
					for j, comp := range addressComponents {
						addressComponents[j] = strings.Title(comp)
					}
					property.Address = strings.Join(addressComponents, " ")
					property.ZipCode = zipCode
				}
			}
			break
		}
	}

	// Try to fetch the page
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// If we can't make the request, return what we parsed from URL
		property.ID = generatePropertyID(property.Address, property.City)
		return property, nil
	}

	// Use browser-like headers for Redfin
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	resp, err := h.client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch Redfin page: %v", err)
		property.ID = generatePropertyID(property.Address, property.City)
		return property, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Redfin returned HTTP %d: %s", resp.StatusCode, resp.Status)
		property.ID = generatePropertyID(property.Address, property.City)
		return property, nil
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse Redfin HTML: %v", err)
		property.ID = generatePropertyID(property.Address, property.City)
		return property, nil
	}

	// Extract data from Redfin page structure
	pageText := doc.Text()

	log.Printf("Extracting comprehensive property details from Redfin...")

	// Enhanced regex patterns for better extraction
	priceRegex := regexp.MustCompile(`\$([0-9,]+(?:\.[0-9]{2})?)`)
	sqftRegex := regexp.MustCompile(`([0-9,]+)\s*(?:sq\.?\s*ft\.?|square\s+feet)`)
	bedroomRegex := regexp.MustCompile(`(\d+)\s*(?:bed|bedroom)`)
	bathroomRegex := regexp.MustCompile(`([0-9.]+)\s*(?:bath|bathroom)`)
	yearRegex := regexp.MustCompile(`(?:built|year)\s*:?\s*(\d{4})`)
	soldRegex := regexp.MustCompile(`SOLD\s+([A-Z]{3}\s+\d{1,2},\s+\d{4})`)
	lotSizeRegex := regexp.MustCompile(`([0-9,.]+)\s*(?:acres?|ac\b)`)
	daysOnMarketRegex := regexp.MustCompile(`(\d+)\s*days?\s*on\s*market`)

	// School-related patterns
	schoolDistrictRegex := regexp.MustCompile(`(?:district|District):\s*([^,\n]+)`)
	elementaryRegex := regexp.MustCompile(`Elementary:\s*([^,\n]+)`)
	middleRegex := regexp.MustCompile(`Middle:\s*([^,\n]+)`)
	highRegex := regexp.MustCompile(`High:\s*([^,\n]+)`)

	// Property condition and features
	conditionRegex := regexp.MustCompile(`Condition:\s*([^,\n]+)`)
	propertyTaxRegex := regexp.MustCompile(`Property\s+Tax:\s*\$([0-9,]+)`)
	hoaRegex := regexp.MustCompile(`HOA:\s*\$([0-9,]+)`)
	parkingRegex := regexp.MustCompile(`(\d+)\s*car\s*garage|(\d+)\s*parking\s*space`)

	// Extract basic property details
	if priceMatch := priceRegex.FindStringSubmatch(pageText); len(priceMatch) > 1 {
		property.Price = parsePrice(priceMatch[1])
	}

	if sqftMatch := sqftRegex.FindStringSubmatch(pageText); len(sqftMatch) > 1 {
		property.SquareFeet = parseInt(sqftMatch[1])
	}

	if bedMatch := bedroomRegex.FindStringSubmatch(pageText); len(bedMatch) > 1 {
		property.Bedrooms = parseInt(bedMatch[1])
	}

	if bathMatch := bathroomRegex.FindStringSubmatch(pageText); len(bathMatch) > 1 {
		property.Bathrooms = parseFloat(bathMatch[1])
	}

	if yearMatch := yearRegex.FindStringSubmatch(pageText); len(yearMatch) > 1 {
		property.YearBuilt = parseInt(yearMatch[1])
	}

	if soldMatch := soldRegex.FindStringSubmatch(pageText); len(soldMatch) > 1 {
		property.SoldDate = soldMatch[1]
	}

	if lotMatch := lotSizeRegex.FindStringSubmatch(pageText); len(lotMatch) > 1 {
		property.LotSize = parseFloat(lotMatch[1])
	}

	if daysMatch := daysOnMarketRegex.FindStringSubmatch(pageText); len(daysMatch) > 1 {
		property.DaysOnMarket = parseInt(daysMatch[1])
	}

	// Extract school information
	if districtMatch := schoolDistrictRegex.FindStringSubmatch(pageText); len(districtMatch) > 1 {
		property.SchoolDistrict = strings.TrimSpace(districtMatch[1])
	}

	if elemMatch := elementaryRegex.FindStringSubmatch(pageText); len(elemMatch) > 1 {
		property.ElementarySchool = strings.TrimSpace(elemMatch[1])
	}

	if middleMatch := middleRegex.FindStringSubmatch(pageText); len(middleMatch) > 1 {
		property.MiddleSchool = strings.TrimSpace(middleMatch[1])
	}

	if highMatch := highRegex.FindStringSubmatch(pageText); len(highMatch) > 1 {
		property.HighSchool = strings.TrimSpace(highMatch[1])
	}

	// Extract property condition and financial details
	if condMatch := conditionRegex.FindStringSubmatch(pageText); len(condMatch) > 1 {
		property.PropertyCondition = strings.TrimSpace(condMatch[1])
	}

	if taxMatch := propertyTaxRegex.FindStringSubmatch(pageText); len(taxMatch) > 1 {
		property.PropertyTax = "$" + taxMatch[1]
	}

	if hoaMatch := hoaRegex.FindStringSubmatch(pageText); len(hoaMatch) > 1 {
		property.HOAFees = "$" + hoaMatch[1]
	}

	// Extract parking information
	if parkingMatch := parkingRegex.FindStringSubmatch(pageText); len(parkingMatch) > 1 {
		if parkingMatch[1] != "" {
			property.ParkingSpaces = parseInt(parkingMatch[1])
			property.Garage = parkingMatch[1] + " car garage"
		} else if parkingMatch[2] != "" {
			property.ParkingSpaces = parseInt(parkingMatch[2])
		}
	}

	// Extract property description from specific Redfin elements
	doc.Find(".remarks, .property-description, .listing-description, .public-remarks").Each(func(i int, s *goquery.Selection) {
		desc := strings.TrimSpace(s.Text())
		if desc != "" && len(desc) > len(property.Description) {
			property.Description = desc
		}
	})

	// If no structured description found, look for description patterns in text
	if property.Description == "" {
		descRegex := regexp.MustCompile(`(?i)description[:\s]*([^.]{50,500}\.?)`)
		if descMatch := descRegex.FindStringSubmatch(pageText); len(descMatch) > 1 {
			property.Description = strings.TrimSpace(descMatch[1])
		}
	}

	// Extract comprehensive features list
	var features []string

	// Look for structured feature lists
	doc.Find(".amenity-group li, .feature-list li, .amenities li, .features li").Each(func(i int, s *goquery.Selection) {
		feature := strings.TrimSpace(s.Text())
		if feature != "" && len(feature) < 100 {
			features = append(features, feature)
		}
	})

	// Extract features from common property details sections
	featurePatterns := []string{
		`Heating:\s*([^,\n]+)`,
		`Cooling:\s*([^,\n]+)`,
		`Flooring:\s*([^,\n]+)`,
		`Appliances:\s*([^,\n]+)`,
		`Roof:\s*([^,\n]+)`,
		`Foundation:\s*([^,\n]+)`,
		`View:\s*([^,\n]+)`,
		`Fireplace:\s*([^,\n]+)`,
	}

	for _, pattern := range featurePatterns {
		regex := regexp.MustCompile(pattern)
		if match := regex.FindStringSubmatch(pageText); len(match) > 1 {
			features = append(features, strings.TrimSpace(match[1]))
		}
	}

	// Add specific fields based on extracted features
	for _, feature := range features {
		lower := strings.ToLower(feature)
		if strings.Contains(lower, "heating") {
			property.Heating = feature
		}
		if strings.Contains(lower, "cooling") || strings.Contains(lower, "air") {
			property.Cooling = feature
		}
		if strings.Contains(lower, "floor") {
			property.Flooring = append(property.Flooring, feature)
		}
		if strings.Contains(lower, "appliance") {
			property.Appliances = append(property.Appliances, feature)
		}
	}

	property.Features = features

	// Set neighborhood to Manoa for Honolulu properties in 96822
	if property.City == "Honolulu" && property.ZipCode == "96822" {
		property.Neighborhood = "Manoa"
	}

	// Extract agent and brokerage information
	doc.Find(".agent-name, .listing-agent, .brokerage-name").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(strings.ToLower(text), "agent") && property.Agent == "" {
			property.Agent = text
		}
		if strings.Contains(strings.ToLower(text), "broker") && property.Brokerage == "" {
			property.Brokerage = text
		}
	})

	// Set property type with more specific detection
	property.PropertyType = "house"
	lowerText := strings.ToLower(pageText)
	if strings.Contains(lowerText, "condominium") || strings.Contains(lowerText, "condo") {
		property.PropertyType = "condo"
	} else if strings.Contains(lowerText, "townhouse") || strings.Contains(lowerText, "townhome") {
		property.PropertyType = "townhouse"
	} else if strings.Contains(lowerText, "apartment") {
		property.PropertyType = "apartment"
	}

	// Generate ID
	property.ID = generatePropertyID(property.Address, property.City)

	// Calculate price per sqft
	if property.Price > 0 && property.SquareFeet > 0 {
		property.PricePerSqFt = property.Price / property.SquareFeet
	}

	return property, nil
}

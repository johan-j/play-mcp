package scraper

import (
	"testing"
)

func TestHomesScraper(t *testing.T) {
	scraper := NewHomesScraper()

	// Test market stats scraping
	t.Run("ScrapeMarketStats", func(t *testing.T) {
		stats, err := scraper.ScrapeMarketStats("Honolulu", "HI")
		if err != nil {
			t.Fatalf("Failed to scrape market stats: %v", err)
		}

		if stats.Area == "" {
			t.Error("Area should not be empty")
		}

		t.Logf("Market Stats: %+v", stats)
	})

	// Test property scraping
	t.Run("ScrapeManoa", func(t *testing.T) {
		properties, err := scraper.ScrapeManoa("sold")
		if err != nil {
			t.Fatalf("Failed to scrape properties: %v", err)
		}

		if len(properties) == 0 {
			t.Error("Should have found some properties")
		}

		t.Logf("Found %d properties", len(properties))
		if len(properties) > 0 {
			t.Logf("First property: %+v", properties[0])
		}
	})
}

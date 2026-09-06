// Package parser parses AirlineSim market page HTML into structured aircraft data.
package parser

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// AircraftOffer represents a single aircraft listing on the market.
type AircraftOffer struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`       // e.g., "Airbus A320-200 heavy"
	Family    string    `json:"family"`     // e.g., "A320 / A321"
	Age       int       `json:"age"`        // in years
	Cycles    int       `json:"cycles"`     // flight cycles
	Condition float64   `json:"condition"`  // percentage 0-100
	Owner     string    `json:"owner"`
	Price     float64   `json:"price"`       // base price (AS$)
	LeaseRate float64   `json:"lease_rate"`  // weekly lease rate (AS$)
	LeaseDep  float64   `json:"lease_dep"`   // lease deposit (AS$)
	DownPmt   float64   `json:"down_pmt"`    // down payment (AS$)
	Install   float64   `json:"install"`     // weekly installment (AS$)
	HasBid    bool      `json:"has_bid"`     // whether there are existing bids
	OfferType string    `json:"offer_type"`  // "auction" or "immediate"
	Financing []string  `json:"financing"`   // available payment methods
	Location  string    `json:"location"`    // airport code
	BidCount  int       `json:"bid_count"`   // number of existing bids
	URL       string    `json:"url"`
	SeenAt    time.Time `json:"seen_at"`
}

// ParseResult contains the results of parsing a market page.
type ParseResult struct {
	Offers    []AircraftOffer `json:"offers"`
	PageCount int             `json:"page_count"`
	Error     string          `json:"error,omitempty"`
}

// Parser parses Wicket market page HTML into aircraft offers.
type Parser struct{}

// New creates a new market page parser.
func New() *Parser {
	return &Parser{}
}

// Parse parses the Wicket market page HTML into structured aircraft data.
func (p *Parser) Parse(body []byte) (*ParseResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	result := &ParseResult{
		Offers: make([]AircraftOffer, 0),
	}

	// The market page uses Wicket panels with repeating rows.
	// Each aircraft listing is in a table row (tr) inside a table.
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			// Skip header rows (they have th elements)
			if row.Find("th").Length() > 0 {
				return
			}
			// Skip rows with only checkbox/empty cells
			cells := row.Find("td")
			if cells.Length() < 5 {
				return
			}
			offer := p.parseRow(row)
			if offer != nil {
				result.Offers = append(result.Offers, *offer)
			}
		})
	})

	return result, nil
}

// parseRow parses a single table row from the Wicket market table.
func (p *Parser) parseRow(row *goquery.Selection) *AircraftOffer {
	cells := row.Find("td")
	if cells.Length() < 6 {
		return nil
	}

	offer := &AircraftOffer{
		SeenAt: time.Now(),
	}

	// The Wicket market table has these columns (based on actual market page):
	// 0: Aircraft name/type (with link to detail)
	// 1: Base price (基準價格)
	// 2: Down payment (頭期款)
	// 3: Installment (分期付款)
	// 4: Lease deposit (租賃押金)
	// 5: Lease rate (租賃費率)
	// Additional columns may include: age, condition, cycles, location, bid status

	cells.Each(func(i int, cell *goquery.Selection) {
		text := strings.TrimSpace(cell.Text())

		switch i {
		case 0:
			// Aircraft type name
			offer.Type = cleanText(text)
			// Extract URL from the link
			if link := cell.Find("a"); link.Length() > 0 {
				if href, exists := link.Attr("href"); exists {
					offer.URL = href
				}
			}
		case 1:
			// Base price
			offer.Price = extractPrice(text)
		case 2:
			// Down payment
			offer.DownPmt = extractPrice(text)
		case 3:
			// Installment (weekly)
			offer.Install = extractPrice(text)
		case 4:
			// Lease deposit
			offer.LeaseDep = extractPrice(text)
		case 5:
			// Lease rate (weekly)
			offer.LeaseRate = extractPrice(text)
		case 6:
			// Age (if present)
			age := extractNumber(text)
			if a, err := strconv.Atoi(age); err == nil {
				offer.Age = a
			}
		case 7:
			// Condition (if present)
			cond := extractNumber(text)
			if c, err := strconv.ParseFloat(cond, 64); err == nil {
				offer.Condition = math.Min(c, 100)
			}
		case 8:
			// Cycles (if present)
			cycles := extractNumber(text)
			if c, err := strconv.Atoi(cycles); err == nil {
				offer.Cycles = c
			}
		case 9:
			// Location
			offer.Location = cleanText(text)
		}
	})

	// Check if the row has bid-related indicators
	// The market page shows "投標" (bid) or "offer" for auction items
	rowHTML, _ := row.Html()
	offer.HasBid = strings.Contains(rowHTML, "投標") || strings.Contains(rowHTML, "bid")
	if offer.HasBid {
		offer.OfferType = "auction"
	} else {
		offer.OfferType = "immediate"
	}

	// Generate a unique ID
	offer.ID = offer.GenerateID()

	if offer.Type == "" {
		return nil
	}

	return offer
}

// cleanText removes extra whitespace and special characters.
func cleanText(s string) string {
	s = strings.TrimSpace(s)
	// Remove non-breaking spaces and other special chars
	s = strings.ReplaceAll(s, "\u00a0", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	// Collapse multiple spaces
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.TrimSpace(s)
}

// extractNumber extracts the first numeric value from a string.
func extractNumber(s string) string {
	s = strings.TrimSpace(s)
	var num strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' || r == '.' || r == ',' {
			if r == ',' {
				r = '.'
			}
			num.WriteRune(r)
		} else if num.Len() > 0 {
			break
		}
	}
	return num.String()
}

// extractPrice parses a price string like "AS$ 1,234,567" or "1.234.567,00" into a float.
func extractPrice(s string) float64 {
	s = strings.TrimSpace(s)
	// Remove currency symbols
	s = strings.ReplaceAll(s, "AS$", "")
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, "\u00a0", " ")
	s = strings.TrimSpace(s)

	// Handle different number formats
	// The game might use European format: 1.234.567,89
	// Or US format: 1,234,567.89
	hasDot := strings.Contains(s, ".")
	hasComma := strings.Contains(s, ",")

	if hasDot && hasComma {
		// Mixed format - determine which is thousands separator
		lastDot := strings.LastIndex(s, ".")
		lastComma := strings.LastIndex(s, ",")
		if lastDot > lastComma {
			// European: 1.234.567,89 -> remove dots, replace comma with dot
			s = strings.ReplaceAll(s, ".", "")
			s = strings.Replace(s, ",", ".", 1)
		} else {
			// US: 1,234,567.89 -> remove commas
			s = strings.ReplaceAll(s, ",", "")
		}
	} else if hasComma {
		// Only commas - could be European or US
		// Check if comma is used as decimal separator
		lastComma := strings.LastIndex(s, ",")
		restAfterComma := s[lastComma+1:]
		if len(restAfterComma) <= 2 && !strings.Contains(restAfterComma, ",") {
			// European format: 1.234.567,89 -> but we only have commas
			// This is likely a thousands separator or decimal
			s = strings.ReplaceAll(s, ",", "")
		} else {
			s = strings.ReplaceAll(s, ",", "")
		}
	} else if hasDot {
		// Only dots - could be thousands separator or decimal
		lastDot := strings.LastIndex(s, ".")
		restAfterDot := s[lastDot+1:]
		if len(restAfterDot) <= 2 && !strings.Contains(restAfterDot, ".") {
			// Decimal point - keep as is
		} else {
			// Thousands separator - remove
			s = strings.ReplaceAll(s, ".", "")
		}
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

// GenerateID generates a unique identifier for an aircraft offer.
func (o *AircraftOffer) GenerateID() string {
	return fmt.Sprintf("%s-%s", o.Type, o.SeenAt.Format("150405"))
}
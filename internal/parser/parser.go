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
// The market page uses a per-aircraft card layout with each aircraft in its own container.
func (p *Parser) Parse(body []byte) (*ParseResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	result := &ParseResult{
		Offers: make([]AircraftOffer, 0),
	}

	// Find the offers container and iterate over each aircraft listing
	// Each aircraft is in a div with class "even" or "odd" inside the offers container
	doc.Find("div.offers > div").Each(func(i int, card *goquery.Selection) {
		// Skip non-aircraft divs
		class, _ := card.Attr("class")
		if class != "even" && class != "odd" {
			return
		}

		offer := p.parseAircraftCard(card)
		if offer != nil {
			result.Offers = append(result.Offers, *offer)
		}
	})

	return result, nil
}

// parseAircraftCard parses a single aircraft listing card.
func (p *Parser) parseAircraftCard(card *goquery.Selection) *AircraftOffer {
	offer := &AircraftOffer{
		SeenAt: time.Now(),
	}

	// 1. Extract aircraft type name and URL from <a class="type">
	typeLink := card.Find("a.type")
	if typeLink.Length() == 0 {
		return nil
	}
	offer.Type = cleanText(typeLink.Text())
	if href, exists := typeLink.Attr("href"); exists {
		offer.URL = href
	}

	// 2. Extract offer type from label
	label := card.Find("span.label")
	if label.Length() > 0 {
		labelText := cleanText(label.Text())
		if strings.Contains(labelText, "官方") || strings.Contains(labelText, "offer") {
			offer.OfferType = "immediate"
		}
	}

	// 3. Extract details from .col-md-4
	detailCol := card.Find("div.col-md-4").First()
	if detailCol.Length() > 0 {
		// Owner
		ownerEl := detailCol.Find("div.owner a")
		if ownerEl.Length() > 0 {
			offer.Owner = cleanText(ownerEl.Text())
		}

		// Registration (not stored in struct, but we can extract it)

		// Age
		ageEl := detailCol.Find("div.age span")
		if ageEl.Length() > 0 {
			ageText := cleanText(ageEl.Text())
			ageText = strings.TrimSuffix(ageText, "年")
			ageText = strings.TrimSuffix(ageText, "years")
			ageText = strings.TrimSpace(ageText)
			if age, err := strconv.ParseFloat(ageText, 64); err == nil {
				offer.Age = int(math.Round(age))
			}
		}

		// Condition
		condEl := detailCol.Find("div.condition span")
		if condEl.Length() > 0 {
			condText := cleanText(condEl.Text())
			condText = strings.TrimSuffix(condText, "%")
			condText = strings.TrimSpace(condText)
			if cond, err := strconv.ParseFloat(condText, 64); err == nil {
				offer.Condition = math.Min(cond, 100)
			}
		}

		// Location
		locEl := detailCol.Find("div.location a, div.location span")
		if locEl.Length() > 0 {
			offer.Location = cleanText(locEl.Text())
		}
	}

	// 4. Extract pricing from the table
	table := card.Find("div.as-table-well table")
	if table.Length() == 0 {
		// Fallback: try any table in the card
		table = card.Find("table")
	}
	if table.Length() > 0 {
		p.parsePricingTable(table, offer)
	}

	// 5. Determine offer type and bid status from the bid button
	bidBtn := card.Find("a.btn-warning")
	if bidBtn.Length() > 0 {
		offer.OfferType = "auction"
		// Check for bid indicator
		btnHTML, _ := bidBtn.Html()
		if strings.Contains(btnHTML, "投標") || strings.Contains(btnHTML, "bid") {
			offer.HasBid = true
		}
	}
	// Check for "立即購買" (immediate buy) button
	buyBtn := card.Find("a.btn-success")
	if buyBtn.Length() > 0 {
		if offer.OfferType == "" {
			offer.OfferType = "immediate"
		}
	}

	// 6. Extract financing options from the bid/buy dropdown
	card.Find("ul.dropdown-menu a").Each(func(i int, a *goquery.Selection) {
		text := cleanText(a.Text())
		if text != "" {
			financing := p.classifyFinancing(text)
			if financing != "" {
				offer.Financing = append(offer.Financing, financing)
			}
		}
	})

	// Generate a unique ID
	if offer.Type != "" {
		offer.ID = offer.GenerateID()
		return offer
	}

	return nil
}

// parsePricingTable parses the pricing table for an aircraft.
// Table structure:
//   Row 0 header: empty, 基準價格, 頭期款, 分期付款, 租賃押金, 租賃費率
//   Row 1: 目前出價, bid_value (or 不適用)
//   Row 2: 下一出價, base_price, down_pmt, installment, lease_dep, lease_rate
//   Row 3: 立即購買, base_price, down_pmt, installment, lease_dep, lease_rate
func (p *Parser) parsePricingTable(table *goquery.Selection, offer *AircraftOffer) {
	table.Find("tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 2 {
			return
		}

		label := cleanText(cells.First().Text())

		switch label {
		case "下一出價", "next bid":
			if cells.Length() >= 6 {
				offer.Price = extractPrice(cleanText(cells.Eq(1).Text()))
				offer.DownPmt = extractPrice(cleanText(cells.Eq(2).Text()))
				offer.Install = extractPrice(cleanText(cells.Eq(3).Text()))
				offer.LeaseDep = extractPrice(cleanText(cells.Eq(4).Text()))
				offer.LeaseRate = extractPrice(cleanText(cells.Eq(5).Text()))
			}
		case "立即購買", "immediate buy", "buy now":
			if cells.Length() >= 6 {
				if offer.Price == 0 {
					offer.Price = extractPrice(cleanText(cells.Eq(1).Text()))
				}
				if offer.DownPmt == 0 {
					offer.DownPmt = extractPrice(cleanText(cells.Eq(2).Text()))
				}
				if offer.Install == 0 {
					offer.Install = extractPrice(cleanText(cells.Eq(3).Text()))
				}
				if offer.LeaseDep == 0 {
					offer.LeaseDep = extractPrice(cleanText(cells.Eq(4).Text()))
				}
				if offer.LeaseRate == 0 {
					offer.LeaseRate = extractPrice(cleanText(cells.Eq(5).Text()))
				}
			}
		case "目前出價", "current bid":
			bidText := cleanText(cells.Eq(1).Text())
			if bidText != "不適用" && bidText != "n/a" && bidText != "" {
				offer.HasBid = true
			}
		}
	})
}

// classifyFinancing categorizes a financing option text into a type.
func (p *Parser) classifyFinancing(text string) string {
	text = strings.ToLower(text)
	switch {
	case strings.Contains(text, "租賃"), strings.Contains(text, "lease"):
		return "lease"
	case strings.Contains(text, "頭期"), strings.Contains(text, "down"), strings.Contains(text, "cash"):
		return "cash"
	case strings.Contains(text, "分期"), strings.Contains(text, "credit"), strings.Contains(text, "install"):
		return "credit"
	}
	return ""
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
	hasDot := strings.Contains(s, ".")
	hasComma := strings.Contains(s, ",")

	if hasDot && hasComma {
		lastDot := strings.LastIndex(s, ".")
		lastComma := strings.LastIndex(s, ",")
		if lastDot > lastComma {
			s = strings.ReplaceAll(s, ".", "")
			s = strings.Replace(s, ",", ".", 1)
		} else {
			s = strings.ReplaceAll(s, ",", "")
		}
	} else if hasComma {
		s = strings.ReplaceAll(s, ",", "")
	} else if hasDot {
		lastDot := strings.LastIndex(s, ".")
		restAfterDot := s[lastDot+1:]
		if len(restAfterDot) > 2 || strings.Contains(restAfterDot, ".") {
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
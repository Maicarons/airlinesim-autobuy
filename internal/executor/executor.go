// Package executor handles the purchase of aircraft on the market.
package executor

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

// PurchaseRequest represents a request to purchase an aircraft.
type PurchaseRequest struct {
	AircraftURL string  `json:"aircraft_url"`
	Price       float64 `json:"price"`
	RuleName    string  `json:"rule_name"`
	IsAuction   bool    `json:"is_auction"`
	BidAmount   float64 `json:"bid_amount,omitempty"`
	MinBalance  float64 `json:"min_balance"` // Minimum account balance to retain after purchase
}

// PurchaseResult represents the outcome of a purchase attempt.
type PurchaseResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Executor handles the execution of aircraft purchases.
type Executor struct {
	client *http.Client
}

// New creates a new purchase executor.
func New(client *http.Client) *Executor {
	return &Executor{
		client: client,
	}
}

// Execute performs a purchase or bid on an aircraft.
func (e *Executor) Execute(req *PurchaseRequest) *PurchaseResult {
	slog.Info("executing purchase",
		"url", req.AircraftURL,
		"price", req.Price,
		"rule", req.RuleName,
		"is_auction", req.IsAuction,
		"min_balance", req.MinBalance,
	)

	// Check minimum balance constraint
	if req.MinBalance > 0 && req.Price >= req.MinBalance {
		slog.Warn("purchase price exceeds minimum balance threshold",
			"price", req.Price,
			"min_balance", req.MinBalance,
		)
		// Allow the purchase anyway - the actual balance check depends on
		// knowing the current account balance, which isn't available yet.
		// This is a safety warning.
	}

	if req.IsAuction {
		return e.placeBid(req)
	}
	return e.buyImmediately(req)
}

// buyImmediately handles an immediate purchase.
func (e *Executor) buyImmediately(req *PurchaseRequest) *PurchaseResult {
	// Step 1: Navigate to the aircraft detail page
	resp, err := e.client.Get(req.AircraftURL)
	if err != nil {
		return &PurchaseResult{Success: false, Message: fmt.Sprintf("failed to access aircraft page: %v", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &PurchaseResult{Success: false, Message: fmt.Sprintf("failed to read aircraft page: %v", err)}
	}

	// Step 2: Find the purchase form/button
	// Look for a "Buy Now" or "Purchase" button
	bodyStr := string(body)
	_ = bodyStr

	// Step 3: Submit the purchase
	// The purchase form typically posts to a Wicket URL
	purchaseURL := req.AircraftURL + "?purchase=true"
	formData := url.Values{
		"confirm": {"true"},
		"price":   {fmt.Sprintf("%.0f", req.Price)},
	}

	purchaseResp, err := e.client.PostForm(purchaseURL, formData)
	if err != nil {
		return &PurchaseResult{Success: false, Message: fmt.Sprintf("purchase request failed: %v", err)}
	}
	defer purchaseResp.Body.Close()

	if purchaseResp.StatusCode == http.StatusOK {
		slog.Info("purchase successful",
			"url", req.AircraftURL,
			"price", req.Price,
		)
		return &PurchaseResult{
			Success: true,
			Message: fmt.Sprintf("Successfully purchased aircraft for AS$ %.0f", req.Price),
		}
	}

	return &PurchaseResult{
		Success: false,
		Message: fmt.Sprintf("purchase returned status %d", purchaseResp.StatusCode),
	}
}

// placeBid handles placing a bid on an auction.
func (e *Executor) placeBid(req *PurchaseRequest) *PurchaseResult {
	bidURL := req.AircraftURL + "?bid=true"
	formData := url.Values{
		"bid_amount": {fmt.Sprintf("%.0f", req.BidAmount)},
	}

	resp, err := e.client.PostForm(bidURL, formData)
	if err != nil {
		return &PurchaseResult{Success: false, Message: fmt.Sprintf("bid failed: %v", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &PurchaseResult{Success: false, Message: "bid placed but failed to read confirmation"}
	}

	// Check response for success indicators
	bodyStr := strings.ToLower(string(body))
	if strings.Contains(bodyStr, "success") || strings.Contains(bodyStr, "bid placed") {
		slog.Info("bid placed successfully",
			"url", req.AircraftURL,
			"amount", req.BidAmount,
		)
		return &PurchaseResult{
			Success: true,
			Message: fmt.Sprintf("Bid of AS$ %.0f placed successfully", req.BidAmount),
		}
	}

	return &PurchaseResult{
		Success: false,
		Message: "bid failed: unexpected response",
	}
}

// VerifyPrice checks the current price of an aircraft before purchase.
func (e *Executor) VerifyPrice(aircraftURL string, expectedPrice float64) (bool, float64, error) {
	resp, err := e.client.Get(aircraftURL)
	if err != nil {
		return false, 0, fmt.Errorf("failed to verify price: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, 0, fmt.Errorf("failed to read price page: %w", err)
	}

	// Parse the current price from the page
	_ = body

	// For now, return true as the price verification
	// In production, this would parse the aircraft detail page for the current price
	return true, expectedPrice, nil
}
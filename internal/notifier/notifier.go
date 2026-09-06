// Package notifier handles sending notifications about events.
package notifier

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Maicarons/airlinesim-autobuy/internal/parser"
)

// EventType represents the type of notification event.
type EventType string

const (
	EventAircraftFound   EventType = "aircraft_found"
	EventPurchaseMade    EventType = "purchase_made"
	EventPurchaseFailed  EventType = "purchase_failed"
	EventBidPlaced       EventType = "bid_placed"
	EventError           EventType = "error"
	EventInfo            EventType = "info"
)

// Event represents a notification event.
type Event struct {
	Type    EventType `json:"type"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// Notifier handles sending notifications through configured channels.
type Notifier struct {
	console bool
}

// New creates a new notifier.
func New(console bool) *Notifier {
	return &Notifier{
		console: console,
	}
}

// Notify sends an event through all configured channels.
func (n *Notifier) Notify(event Event) {
	if n.console {
		n.notifyConsole(event)
	}
}

// NotifyAircraftFound sends a notification about a found aircraft.
func (n *Notifier) NotifyAircraftFound(aircraft *parser.AircraftOffer, ruleName string) {
	var details strings.Builder
	details.WriteString(fmt.Sprintf("Type: %s | Price: AS$ %.0f | Age: %dy | Cond: %.0f%%",
		aircraft.Type, aircraft.Price, aircraft.Age, aircraft.Condition))
	if aircraft.OfferType != "" {
		details.WriteString(fmt.Sprintf(" | Type: %s", aircraft.OfferType))
	}

	n.Notify(Event{
		Type:    EventAircraftFound,
		Message: fmt.Sprintf("Aircraft found matching rule '%s': %s", ruleName, aircraft.Type),
		Details: details.String(),
	})
}

// NotifyPurchaseMade sends a notification about a successful purchase.
func (n *Notifier) NotifyPurchaseMade(aircraftType string, price float64, ruleName string) {
	n.Notify(Event{
		Type:    EventPurchaseMade,
		Message: fmt.Sprintf("PURCHASED %s for AS$ %.0f (rule: %s)", aircraftType, price, ruleName),
	})
}

// NotifyPurchaseFailed sends a notification about a failed purchase.
func (n *Notifier) NotifyPurchaseFailed(aircraftType string, reason string) {
	n.Notify(Event{
		Type:    EventPurchaseFailed,
		Message: fmt.Sprintf("Failed to purchase %s: %s", aircraftType, reason),
	})
}

// NotifyError sends an error notification.
func (n *Notifier) NotifyError(err error) {
	n.Notify(Event{
		Type:    EventError,
		Message: fmt.Sprintf("Error: %v", err),
	})
}

// NotifyInfo sends an informational notification.
func (n *Notifier) NotifyInfo(msg string) {
	n.Notify(Event{
		Type:    EventInfo,
		Message: msg,
	})
}

// notifyConsole logs the event to the console.
func (n *Notifier) notifyConsole(event Event) {
	switch event.Type {
	case EventPurchaseMade:
		slog.Info("🛒 " + event.Message)
	case EventPurchaseFailed:
		slog.Warn("❌ " + event.Message)
	case EventError:
		slog.Error("⚠️ " + event.Message)
	case EventAircraftFound:
		slog.Info("🔍 " + event.Message, "details", event.Details)
	case EventBidPlaced:
		slog.Info("💰 " + event.Message)
	default:
		slog.Info("ℹ️ " + event.Message)
	}
}
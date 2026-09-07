// Package notifier handles sending notifications about events.
package notifier

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

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
	console         bool
	dingtalkWebhook string
	dingtalkSecret  string
	httpClient      *http.Client
}

// New creates a new notifier.
func New(console bool, dingtalkWebhook, dingtalkSecret string) *Notifier {
	return &Notifier{
		console:         console,
		dingtalkWebhook: dingtalkWebhook,
		dingtalkSecret:  dingtalkSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Notify sends an event through all configured channels.
func (n *Notifier) Notify(event Event) {
	if n.console {
		n.notifyConsole(event)
	}
	if n.dingtalkWebhook != "" {
		n.notifyDingTalk(event)
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

// dingtalkMessage represents the JSON payload for a DingTalk custom robot message.
type dingtalkMessage struct {
	MsgType string       `json:"msgtype"`
	Text    dingtalkText `json:"text"`
}

type dingtalkText struct {
	Content string `json:"content"`
}

// dingtalkResponse represents the response from DingTalk webhook API.
type dingtalkResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// notifyDingTalk sends an event to DingTalk via custom robot webhook.
func (n *Notifier) notifyDingTalk(event Event) {
	// Build the message text
	var content string
	switch event.Type {
	case EventPurchaseMade:
		content = "🛒 " + event.Message
	case EventPurchaseFailed:
		content = "❌ " + event.Message
	case EventError:
		content = "⚠️ " + event.Message
	case EventAircraftFound:
		content = "🔍 " + event.Message
		if event.Details != "" {
			content += "\n" + event.Details
		}
	case EventBidPlaced:
		content = "💰 " + event.Message
	default:
		content = "ℹ️ " + event.Message
	}

	msg := dingtalkMessage{
		MsgType: "text",
		Text:    dingtalkText{Content: content},
	}

	body, err := json.Marshal(msg)
	if err != nil {
		slog.Error("failed to marshal DingTalk message", "error", err)
		return
	}

	// Build the webhook URL
	webhookURL := n.dingtalkWebhook

	// If a secret is configured, append timestamp and sign to the URL
	if n.dingtalkSecret != "" {
		timestamp := time.Now().UnixMilli()
		sign := n.signDingTalk(n.dingtalkSecret, timestamp)
		sep := "&"
		if !strings.Contains(webhookURL, "?") {
			sep = "?"
		}
		webhookURL = fmt.Sprintf("%s%stimestamp=%d&sign=%s", webhookURL, sep, timestamp, url.QueryEscape(sign))
	}

	resp, err := n.httpClient.Post(webhookURL, "application/json; charset=utf-8", bytes.NewReader(body))
	if err != nil {
		slog.Error("failed to send DingTalk notification", "error", err)
		return
	}
	defer resp.Body.Close()

	var result dingtalkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("failed to decode DingTalk response", "error", err)
		return
	}

	if result.ErrCode != 0 {
		slog.Error("DingTalk API error", "errcode", result.ErrCode, "errmsg", result.ErrMsg)
	}
}

// signDingTalk generates the HMAC-SHA256 signature for DingTalk webhook authentication.
func (n *Notifier) signDingTalk(secret string, timestamp int64) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
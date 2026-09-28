// Package payment provides multi-provider payment gateway integration for Everato.
// Supported payment providers include Razorpay, Stripe, and Cashfree.
package payment

import (
	"context"
	"fmt"
	"strings"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/utils"
)

// ProviderType identifies the active payment gateway.
type ProviderType string

const (
	ProviderRazorpay ProviderType = "razorpay"
	ProviderStripe   ProviderType = "stripe"
	ProviderCashfree ProviderType = "cashfree"
)

// CreateOrderRequest is the parameters for initializing a checkout order.
type CreateOrderRequest struct {
	BookingID     string            `json:"booking_id"`
	Amount        float64           `json:"amount"` // in major currency units (e.g. 100.50 INR)
	Currency      string            `json:"currency"`
	CustomerName  string            `json:"customer_name"`
	CustomerEmail string            `json:"customer_email"`
	CustomerPhone string            `json:"customer_phone"`
	Receipt       string            `json:"receipt"`
	Notes         map[string]string `json:"notes,omitempty"`
}

// CreateOrderResponse holds the order initialization returned to the client.
type CreateOrderResponse struct {
	OrderID      string            `json:"order_id"`
	Provider     string            `json:"provider"`
	Amount       float64           `json:"amount"`
	Currency     string            `json:"currency"`
	KeyID        string            `json:"key_id,omitempty"`        // Client public key / App ID
	ClientSecret string            `json:"client_secret,omitempty"` // For Stripe PaymentIntent
	SessionID    string            `json:"session_id,omitempty"`    // For Cashfree payment_session_id
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// VerifyPaymentRequest contains credentials returned from frontend checkout.
type VerifyPaymentRequest struct {
	BookingID string `json:"booking_id"`
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Signature string `json:"signature,omitempty"`
}

// VerifyPaymentResponse indicates payment verification result.
type VerifyPaymentResponse struct {
	Success   bool   `json:"success"`
	PaymentID string `json:"payment_id"`
	OrderID   string `json:"order_id"`
	Message   string `json:"message"`
}

// WebhookEvent represents an incoming webhook notification from the provider.
type WebhookEvent struct {
	EventID   string
	EventType string
	BookingID string
	OrderID   string
	PaymentID string
	Status    string // "SUCCESS", "FAILED"
	RawBody   []byte
}

// PaymentGateway defines the common contract implemented by all payment providers.
type PaymentGateway interface {
	GetProvider() ProviderType
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error)
	VerifyPayment(ctx context.Context, req VerifyPaymentRequest) (*VerifyPaymentResponse, error)
	HandleWebhook(ctx context.Context, body []byte, signature string) (*WebhookEvent, error)
	GetClientConfig() map[string]string
}

// NewPaymentGateway is a factory that reads configuration and creates the chosen payment gateway.
// It validates that required environment variables are present for the chosen provider.
func NewPaymentGateway(cfg *config.Config) (PaymentGateway, error) {
	providerName := strings.ToLower(strings.TrimSpace(cfg.Payment.Provider))
	if providerName == "" {
		providerName = strings.ToLower(strings.TrimSpace(utils.GetEnv("PAYMENT_PROVIDER", "razorpay")))
	}

	currency := strings.ToUpper(strings.TrimSpace(cfg.Payment.Currency))
	if currency == "" {
		currency = strings.ToUpper(strings.TrimSpace(utils.GetEnv("PAYMENT_CURRENCY", "INR")))
	}

	switch ProviderType(providerName) {
	case ProviderRazorpay:
		return NewRazorpayProvider(currency)

	case ProviderStripe:
		return NewStripeProvider(currency)

	case ProviderCashfree:
		return NewCashfreeProvider(currency)

	default:
		return nil, fmt.Errorf(
			"unsupported payment provider '%s'. Allowed options: 'razorpay', 'stripe', 'cashfree'",
			providerName,
		)
	}
}

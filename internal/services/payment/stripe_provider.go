package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// StripeProvider implements PaymentGateway for Stripe.
type StripeProvider struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
	Currency       string
	HTTPClient     *http.Client
}

// NewStripeProvider initializes a Stripe client, validating required environment variables.
func NewStripeProvider(currency string) (*StripeProvider, error) {
	secretKey := utils.GetEnv("STRIPE_SECRET_KEY", "")
	publishableKey := utils.GetEnv("STRIPE_PUBLISHABLE_KEY", "")
	webhookSecret := utils.GetEnv("STRIPE_WEBHOOK_SECRET", "")

	var missing []string
	if secretKey == "" {
		missing = append(missing, "STRIPE_SECRET_KEY")
	}
	if publishableKey == "" {
		missing = append(missing, "STRIPE_PUBLISHABLE_KEY")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf(
			"Stripe selected as payment provider, but required environment variables are missing: %v. Please add them to your .env file or environment",
			missing,
		)
	}

	if currency == "" {
		currency = "USD"
	}

	return &StripeProvider{
		SecretKey:      secretKey,
		PublishableKey: publishableKey,
		WebhookSecret:  webhookSecret,
		Currency:       currency,
		HTTPClient:     &http.Client{Timeout: 15 * time.Second},
	}, nil
}

func (s *StripeProvider) GetProvider() ProviderType {
	return ProviderStripe
}

// CreateOrder creates a Stripe PaymentIntent: POST https://api.stripe.com/v1/payment_intents
func (s *StripeProvider) CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Amount in cents / minor currency units
	amountInCents := int(math.Round(req.Amount * 100))

	currency := strings.ToLower(req.Currency)
	if currency == "" {
		currency = strings.ToLower(s.Currency)
	}

	data := url.Values{}
	data.Set("amount", strconv.Itoa(amountInCents))
	data.Set("currency", currency)
	data.Set("metadata[booking_id]", req.BookingID)
	data.Set("metadata[customer_name]", req.CustomerName)
	data.Set("metadata[customer_email]", req.CustomerEmail)
	data.Add("payment_method_types[]", "card")

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.stripe.com/v1/payment_intents",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("create stripe http request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+s.SecretKey)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute stripe payment intent request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read stripe response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.Error("Stripe PaymentIntent creation failed", "status", resp.StatusCode, "body", string(respBytes))
		return nil, fmt.Errorf("stripe error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var parsed struct {
		ID           string `json:"id"`
		ClientSecret string `json:"client_secret"`
		Amount       int    `json:"amount"`
		Currency     string `json:"currency"`
		Status       string `json:"status"`
	}
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal stripe response: %w", err)
	}

	return &CreateOrderResponse{
		OrderID:      parsed.ID,
		Provider:     string(ProviderStripe),
		Amount:       req.Amount,
		Currency:     strings.ToUpper(parsed.Currency),
		KeyID:        s.PublishableKey,
		ClientSecret: parsed.ClientSecret,
		Metadata: map[string]string{
			"status": parsed.Status,
		},
	}, nil
}

// VerifyPayment checks the status of a Stripe PaymentIntent.
func (s *StripeProvider) VerifyPayment(ctx context.Context, req VerifyPaymentRequest) (*VerifyPaymentResponse, error) {
	intentID := req.PaymentID
	if intentID == "" {
		intentID = req.OrderID
	}
	if intentID == "" {
		return nil, fmt.Errorf("payment_id or order_id (PaymentIntent ID) is required for Stripe verification")
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.stripe.com/v1/payment_intents/"+intentID,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create stripe verify request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+s.SecretKey)

	resp, err := s.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute stripe verify request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read stripe verify response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &VerifyPaymentResponse{
			Success:   false,
			PaymentID: intentID,
			Message:   fmt.Sprintf("Stripe payment intent lookup failed with status %d", resp.StatusCode),
		}, fmt.Errorf("stripe status %d: %s", resp.StatusCode, string(respBytes))
	}

	var parsed struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal stripe verify response: %w", err)
	}

	if parsed.Status != "succeeded" && parsed.Status != "processing" {
		return &VerifyPaymentResponse{
			Success:   false,
			PaymentID: parsed.ID,
			Message:   fmt.Sprintf("Payment intent status is '%s', not succeeded", parsed.Status),
		}, fmt.Errorf("payment not succeeded: status=%s", parsed.Status)
	}

	return &VerifyPaymentResponse{
		Success:   true,
		PaymentID: parsed.ID,
		OrderID:   parsed.ID,
		Message:   "Stripe PaymentIntent confirmed successfully",
	}, nil
}

// HandleWebhook decodes Stripe webhook event notifications.
func (s *StripeProvider) HandleWebhook(ctx context.Context, body []byte, signature string) (*WebhookEvent, error) {
	var payload struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Data struct {
			Object struct {
				ID       string            `json:"id"`
				Status   string            `json:"status"`
				Metadata map[string]string `json:"metadata"`
			} `json:"object"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal stripe webhook: %w", err)
	}

	status := "FAILED"
	if payload.Type == "payment_intent.succeeded" || payload.Data.Object.Status == "succeeded" {
		status = "SUCCESS"
	}

	return &WebhookEvent{
		EventID:   payload.ID,
		EventType: payload.Type,
		PaymentID: payload.Data.Object.ID,
		BookingID: payload.Data.Object.Metadata["booking_id"],
		Status:    status,
		RawBody:   body,
	}, nil
}

func (s *StripeProvider) GetClientConfig() map[string]string {
	return map[string]string{
		"provider":        string(ProviderStripe),
		"publishable_key": s.PublishableKey,
		"currency":        s.Currency,
	}
}

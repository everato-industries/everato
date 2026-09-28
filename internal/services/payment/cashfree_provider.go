package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// CashfreeProvider implements PaymentGateway for Cashfree PG.
type CashfreeProvider struct {
	AppID      string
	SecretKey  string
	Env        string // "TEST" or "PRODUCTION"
	BaseURL    string
	Currency   string
	HTTPClient *http.Client
}

// NewCashfreeProvider initializes a Cashfree client, validating required environment variables.
func NewCashfreeProvider(currency string) (*CashfreeProvider, error) {
	appID := utils.GetEnv("CASHFREE_APP_ID", "")
	secretKey := utils.GetEnv("CASHFREE_SECRET_KEY", "")
	env := strings.ToUpper(utils.GetEnv("CASHFREE_ENV", "TEST"))

	var missing []string
	if appID == "" {
		missing = append(missing, "CASHFREE_APP_ID")
	}
	if secretKey == "" {
		missing = append(missing, "CASHFREE_SECRET_KEY")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf(
			"Cashfree selected as payment provider, but required environment variables are missing: %v. Please add them to your .env file or environment",
			missing,
		)
	}

	baseURL := "https://sandbox.cashfree.com/pg"
	if env == "PRODUCTION" || env == "PROD" {
		baseURL = "https://api.cashfree.com/pg"
	}

	if currency == "" {
		currency = "INR"
	}

	return &CashfreeProvider{
		AppID:      appID,
		SecretKey:  secretKey,
		Env:        env,
		BaseURL:    baseURL,
		Currency:   currency,
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
	}, nil
}

func (c *CashfreeProvider) GetProvider() ProviderType {
	return ProviderCashfree
}

// CreateOrder calls Cashfree PG API to create an order and retrieve payment_session_id.
// POST /pg/orders
func (c *CashfreeProvider) CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error) {
	logger := pkg.NewLogger()
	defer logger.Close()

	currency := strings.ToUpper(req.Currency)
	if currency == "" {
		currency = c.Currency
	}

	// Order ID must be alphanumeric and up to 50 chars
	orderID := fmt.Sprintf("order_%s_%d", req.BookingID[:8], time.Now().Unix())

	customerID := req.CustomerEmail
	if customerID == "" {
		customerID = "cust_" + req.BookingID[:8]
	}
	// Sanitize customer ID
	customerID = strings.ReplaceAll(customerID, "@", "_at_")
	customerID = strings.ReplaceAll(customerID, ".", "_")

	customerPhone := req.CustomerPhone
	if customerPhone == "" {
		customerPhone = "9999999999" // Fallback valid 10-digit phone for sandbox requirements
	}

	customerName := req.CustomerName
	if customerName == "" {
		customerName = "Guest Attendee"
	}

	payload := map[string]any{
		"order_id":       orderID,
		"order_amount":   req.Amount,
		"order_currency": currency,
		"customer_details": map[string]string{
			"customer_id":    customerID,
			"customer_name":  customerName,
			"customer_email": req.CustomerEmail,
			"customer_phone": customerPhone,
		},
		"order_note": fmt.Sprintf("Booking %s", req.BookingID),
		"order_meta": map[string]string{
			"return_url": fmt.Sprintf("http://localhost:8080/my-tickets?order_id=%s", orderID),
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal cashfree order payload: %w", err)
	}

	url := c.BaseURL + "/orders"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create cashfree http request: %w", err)
	}

	httpReq.Header.Set("x-client-id", c.AppID)
	httpReq.Header.Set("x-client-secret", c.SecretKey)
	httpReq.Header.Set("x-api-version", "2023-08-01")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute cashfree order request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read cashfree response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.Error("Cashfree order creation failed", "status", resp.StatusCode, "body", string(respBytes))
		return nil, fmt.Errorf("cashfree error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var parsed struct {
		CFOrderID       string  `json:"cf_order_id"`
		OrderID         string  `json:"order_id"`
		PaymentSessionID string `json:"payment_session_id"`
		OrderStatus     string  `json:"order_status"`
		OrderAmount     float64 `json:"order_amount"`
		OrderCurrency   string  `json:"order_currency"`
	}
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal cashfree response: %w", err)
	}

	return &CreateOrderResponse{
		OrderID:   parsed.OrderID,
		Provider:  string(ProviderCashfree),
		Amount:    parsed.OrderAmount,
		Currency:  parsed.OrderCurrency,
		KeyID:     c.AppID,
		SessionID: parsed.PaymentSessionID,
		Metadata: map[string]string{
			"cf_order_id": parsed.CFOrderID,
			"env":         c.Env,
		},
	}, nil
}

// VerifyPayment checks the order status via Cashfree PG API: GET /pg/orders/{order_id}
func (c *CashfreeProvider) VerifyPayment(ctx context.Context, req VerifyPaymentRequest) (*VerifyPaymentResponse, error) {
	orderID := req.OrderID
	if orderID == "" {
		orderID = req.PaymentID
	}
	if orderID == "" {
		return nil, fmt.Errorf("order_id is required for Cashfree verification")
	}

	url := fmt.Sprintf("%s/orders/%s", c.BaseURL, orderID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create cashfree verify request: %w", err)
	}

	httpReq.Header.Set("x-client-id", c.AppID)
	httpReq.Header.Set("x-client-secret", c.SecretKey)
	httpReq.Header.Set("x-api-version", "2023-08-01")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute cashfree verify request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read cashfree verify response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &VerifyPaymentResponse{
			Success: false,
			OrderID: orderID,
			Message: fmt.Sprintf("Cashfree order verification returned HTTP %d", resp.StatusCode),
		}, fmt.Errorf("cashfree status %d: %s", resp.StatusCode, string(respBytes))
	}

	var parsed struct {
		OrderID     string `json:"order_id"`
		OrderStatus string `json:"order_status"`
	}
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal cashfree status: %w", err)
	}

	if parsed.OrderStatus != "PAID" {
		return &VerifyPaymentResponse{
			Success: false,
			OrderID: parsed.OrderID,
			Message: fmt.Sprintf("Cashfree order status is '%s', expected 'PAID'", parsed.OrderStatus),
		}, fmt.Errorf("order not paid: %s", parsed.OrderStatus)
	}

	return &VerifyPaymentResponse{
		Success: true,
		OrderID: parsed.OrderID,
		Message: "Cashfree payment verified successfully",
	}, nil
}

// HandleWebhook decodes Cashfree webhook payloads.
func (c *CashfreeProvider) HandleWebhook(ctx context.Context, body []byte, signature string) (*WebhookEvent, error) {
	var payload struct {
		Data struct {
			Order struct {
				OrderID string `json:"order_id"`
			} `json:"order"`
			Payment struct {
				PaymentStatus string `json:"payment_status"`
				CFPaymentID   string `json:"cf_payment_id"`
			} `json:"payment"`
		} `json:"data"`
		Type string `json:"type"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal cashfree webhook: %w", err)
	}

	status := "FAILED"
	if payload.Data.Payment.PaymentStatus == "SUCCESS" {
		status = "SUCCESS"
	}

	return &WebhookEvent{
		EventType: payload.Type,
		OrderID:   payload.Data.Order.OrderID,
		PaymentID: payload.Data.Payment.CFPaymentID,
		Status:    status,
		RawBody:   body,
	}, nil
}

func (c *CashfreeProvider) GetClientConfig() map[string]string {
	return map[string]string{
		"provider": string(ProviderCashfree),
		"app_id":   c.AppID,
		"env":      c.Env,
		"currency": c.Currency,
	}
}

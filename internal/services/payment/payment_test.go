package payment_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/services/payment"
)

func TestPaymentGatewayProviderValidation(t *testing.T) {
	// Clean env before test
	os.Unsetenv("RAZORPAY_KEY_ID")
	os.Unsetenv("RAZORPAY_KEY_SECRET")
	os.Unsetenv("STRIPE_SECRET_KEY")
	os.Unsetenv("STRIPE_PUBLISHABLE_KEY")
	os.Unsetenv("CASHFREE_APP_ID")
	os.Unsetenv("CASHFREE_SECRET_KEY")

	// 1. Test unsupported provider
	cfg := &config.Config{
		Payment: config.PaymentConfig{
			Provider: "paypal_unsupported",
		},
	}
	_, err := payment.NewPaymentGateway(cfg)
	if err == nil || !strings.Contains(err.Error(), "unsupported payment provider") {
		t.Errorf("Expected unsupported provider error, got: %v", err)
	}

	// 2. Test Razorpay missing env vars
	cfg.Payment.Provider = "razorpay"
	_, err = payment.NewPaymentGateway(cfg)
	if err == nil {
		t.Fatal("Expected error when Razorpay credentials missing")
	}
	if !strings.Contains(err.Error(), "RAZORPAY_KEY_ID") || !strings.Contains(err.Error(), "RAZORPAY_KEY_SECRET") {
		t.Errorf("Expected message specifying missing Razorpay keys, got: %v", err)
	}

	// 3. Test Stripe missing env vars
	cfg.Payment.Provider = "stripe"
	_, err = payment.NewPaymentGateway(cfg)
	if err == nil {
		t.Fatal("Expected error when Stripe credentials missing")
	}
	if !strings.Contains(err.Error(), "STRIPE_SECRET_KEY") || !strings.Contains(err.Error(), "STRIPE_PUBLISHABLE_KEY") {
		t.Errorf("Expected message specifying missing Stripe keys, got: %v", err)
	}

	// 4. Test Cashfree missing env vars
	cfg.Payment.Provider = "cashfree"
	_, err = payment.NewPaymentGateway(cfg)
	if err == nil {
		t.Fatal("Expected error when Cashfree credentials missing")
	}
	if !strings.Contains(err.Error(), "CASHFREE_APP_ID") || !strings.Contains(err.Error(), "CASHFREE_SECRET_KEY") {
		t.Errorf("Expected message specifying missing Cashfree keys, got: %v", err)
	}
}

func TestRazorpaySignatureVerification(t *testing.T) {
	keyID := "rzp_test_sample123"
	keySecret := "sample_secret_key_abc456"

	t.Setenv("RAZORPAY_KEY_ID", keyID)
	t.Setenv("RAZORPAY_KEY_SECRET", keySecret)

	provider, err := payment.NewRazorpayProvider("INR")
	if err != nil {
		t.Fatalf("Failed to initialize Razorpay provider: %v", err)
	}

	orderID := "order_EKwxwAgItmmXdp"
	paymentID := "pay_29QQoUBi66xm2f"

	// Calculate valid HMAC-SHA256 signature
	data := orderID + "|" + paymentID
	mac := hmac.New(sha256.New, []byte(keySecret))
	mac.Write([]byte(data))
	validSignature := hex.EncodeToString(mac.Sum(nil))

	// Test valid verification
	resp, err := provider.VerifyPayment(context.Background(), payment.VerifyPaymentRequest{
		OrderID:   orderID,
		PaymentID: paymentID,
		Signature: validSignature,
	})
	if err != nil || !resp.Success {
		t.Fatalf("Expected valid signature to verify successfully, err: %v", err)
	}

	// Test tampered signature
	_, err = provider.VerifyPayment(context.Background(), payment.VerifyPaymentRequest{
		OrderID:   orderID,
		PaymentID: paymentID,
		Signature: "tampered_fake_signature_hex",
	})
	if err == nil {
		t.Error("Expected tampered signature to fail verification")
	}

	// Test client config
	clientConfig := provider.GetClientConfig()
	if clientConfig["provider"] != "razorpay" || clientConfig["key_id"] != keyID {
		t.Errorf("Unexpected client config: %+v", clientConfig)
	}
}

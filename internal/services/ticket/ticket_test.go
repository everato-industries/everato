package ticket_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/dtg-lucifer/everato/internal/services/ticket"
)

func TestBuildAndVerifyQRData(t *testing.T) {
	secret := "test-secret-key-12345"
	ticketID := "550e8400-e29b-41d4-a716-446655440000"
	bookingID := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	eventID := "6ba7b811-9dad-11d1-80b4-00c04fd430c8"
	userID := "6ba7b812-9dad-11d1-80b4-00c04fd430c8"
	ticketNum := "TECHCONF-20260929-550E84"

	// 1. Build signed QR data
	qrData, err := ticket.BuildSignedQRData(ticketID, bookingID, eventID, userID, ticketNum, secret)
	if err != nil {
		t.Fatalf("BuildSignedQRData failed: %v", err)
	}
	if qrData == "" {
		t.Fatal("Expected non-empty QR data")
	}

	// 2. Verify with correct secret
	payload, err := ticket.VerifyQRData(qrData, secret)
	if err != nil {
		t.Fatalf("VerifyQRData with correct secret failed: %v", err)
	}
	if payload.TicketID != ticketID {
		t.Errorf("Expected TicketID %s, got %s", ticketID, payload.TicketID)
	}
	if payload.TicketNum != ticketNum {
		t.Errorf("Expected TicketNum %s, got %s", ticketNum, payload.TicketNum)
	}
	if payload.EventID != eventID {
		t.Errorf("Expected EventID %s, got %s", eventID, payload.EventID)
	}

	// 3. Verify with wrong secret must fail
	_, err = ticket.VerifyQRData(qrData, "wrong-secret-key")
	if err == nil {
		t.Error("Expected verification to fail with wrong secret, but succeeded")
	}

	// 4. Tampered data must fail
	tampered := strings.Replace(qrData, ticketNum, "HACKED-TICKET-000", 1)
	_, err = ticket.VerifyQRData(tampered, secret)
	if err == nil {
		t.Error("Expected verification to fail on tampered data, but succeeded")
	}
}

func TestGenerateQRImage(t *testing.T) {
	content := "https://everato.io/tickets/verify/test-123"
	b64Image, err := ticket.GenerateQRImage(content, 200)
	if err != nil {
		t.Fatalf("GenerateQRImage failed: %v", err)
	}
	if b64Image == "" {
		t.Fatal("Expected non-empty base64 string")
	}

	// Verify it's valid base64
	decoded, err := base64.StdEncoding.DecodeString(b64Image)
	if err != nil {
		t.Fatalf("Invalid base64 encoding: %v", err)
	}
	// PNG magic bytes: 0x89 'P' 'N' 'G'
	if len(decoded) < 8 || decoded[0] != 0x89 || decoded[1] != 'P' || decoded[2] != 'N' || decoded[3] != 'G' {
		t.Fatal("Decoded image is not a valid PNG")
	}

	dataURL, err := ticket.DataURLForQR(content, 200)
	if err != nil {
		t.Fatalf("DataURLForQR failed: %v", err)
	}
	if !strings.HasPrefix(dataURL, "data:image/png;base64,") {
		t.Errorf("Data URL missing expected prefix: %s", dataURL[:30])
	}
}

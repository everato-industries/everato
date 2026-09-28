package checkin_test

import (
	"testing"
	"time"

	"github.com/dtg-lucifer/everato/internal/services/ticket"
)

func TestQRCheckinPayloadIntegrity(t *testing.T) {
	secret := "super-secure-production-key"
	ticketID := "e3b0c442-98fc-1c14-9afbf4c8996fb924"
	bookingID := "b1a7b810-9dad-11d1-80b4-00c04fd430c8"
	eventID := "c2a7b811-9dad-11d1-80b4-00c04fd430c8"
	userID := "d3a7b812-9dad-11d1-80b4-00c04fd430c8"
	ticketNum := "TECHFEST-2026-T1001"

	qrData, err := ticket.BuildSignedQRData(ticketID, bookingID, eventID, userID, ticketNum, secret)
	if err != nil {
		t.Fatalf("Failed to build signed QR data: %v", err)
	}

	// Verify valid scan
	payload, err := ticket.VerifyQRData(qrData, secret)
	if err != nil {
		t.Fatalf("Expected valid QR scan verification, got: %v", err)
	}

	if payload.TicketID != ticketID {
		t.Errorf("Expected ticket ID %s, got %s", ticketID, payload.TicketID)
	}

	if payload.IssuedAt <= 0 || payload.IssuedAt > time.Now().Unix()+60 {
		t.Errorf("Invalid IssuedAt timestamp: %d", payload.IssuedAt)
	}
}

// Package ticket provides services for QR-based ticket generation and management
// in the Everato platform. It handles ticket creation after booking confirmation,
// QR code generation, and ticket validation.
package ticket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ─── DTOs ────────────────────────────────────────────────────────────────────

// GenerateTicketsRequest is the input for generating tickets for a confirmed booking.
type GenerateTicketsRequest struct {
	BookingID     string
	EventID       string
	EventSlug     string
	TicketTypeID  string
	TicketTypeName string
	Quantity      int
	UserID        string
	HMACSecret    string // from config
}

// TicketResponse represents a single generated ticket returned to the caller.
type TicketResponse struct {
	ID           string `json:"id"`
	BookingID    string `json:"booking_id"`
	TicketNumber string `json:"ticket_number"`
	QRCodeData   string `json:"qr_code_data"`
	QRCodeImage  string `json:"qr_code_image,omitempty"` // base64 PNG
	IsCheckedIn  bool   `json:"is_checked_in"`
	CreatedAt    string `json:"created_at"`
}

// QRPayload is the JSON data embedded into the QR code.
// It is HMAC-signed to prevent forgery.
type QRPayload struct {
	TicketID   string `json:"tid"`
	BookingID  string `json:"bid"`
	EventID    string `json:"eid"`
	UserID     string `json:"uid"`
	TicketNum  string `json:"num"`
	IssuedAt   int64  `json:"iat"`
	Sig        string `json:"sig,omitempty"` // appended after signing
}

// signPayload creates an HMAC-SHA256 signature over the canonical payload fields.
// We sign: tid|bid|eid|uid|num|iat
func signPayload(p QRPayload, secret string) string {
	msg := fmt.Sprintf("%s|%s|%s|%s|%s|%d",
		p.TicketID, p.BookingID, p.EventID, p.UserID, p.TicketNum, p.IssuedAt)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

// BuildSignedQRData builds the final QR string (signed JSON compact).
func BuildSignedQRData(ticketID, bookingID, eventID, userID, ticketNum, secret string) (string, error) {
	payload := QRPayload{
		TicketID:  ticketID,
		BookingID: bookingID,
		EventID:   eventID,
		UserID:    userID,
		TicketNum: ticketNum,
		IssuedAt:  time.Now().Unix(),
	}
	payload.Sig = signPayload(payload, secret)

	b, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal QR payload: %w", err)
	}
	return string(b), nil
}

// VerifyQRData parses and HMAC-verifies a QR code string.
// Returns the decoded payload if valid, or an error if tampered.
func VerifyQRData(qrData, secret string) (*QRPayload, error) {
	var payload QRPayload
	if err := json.Unmarshal([]byte(qrData), &payload); err != nil {
		return nil, fmt.Errorf("invalid QR format: %w", err)
	}

	sig := payload.Sig
	payload.Sig = ""
	expected := signPayload(payload, secret)

	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return nil, fmt.Errorf("invalid QR signature — ticket may be forged")
	}

	payload.Sig = sig
	return &payload, nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// pgTextToString safely extracts a Go string from pgtype.Text.
func pgTextToString(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}

// pgTimeToString safely formats a pgtype.Timestamptz.
func pgTimeToString(t pgtype.Timestamptz) string {
	if t.Valid {
		return t.Time.UTC().Format(time.RFC3339)
	}
	return ""
}

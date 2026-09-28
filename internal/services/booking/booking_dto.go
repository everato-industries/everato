// Package booking provides services for booking management in the Everato platform.
// It handles the creation, retrieval, and cancellation of bookings.
package booking

import (
	"fmt"
	"math/big"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

// ValidationError represents a validation error for a specific field.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// numericToFloat converts pgtype.Numeric to float64
func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	// Convert big.Int to float64 using the exponent
	f := new(big.Float).SetInt(n.Int)
	factor := new(big.Float).SetFloat64(1)
	for i := int32(0); i < -n.Exp; i++ {
		factor.Mul(factor, new(big.Float).SetFloat64(0.1))
	}
	result, _ := new(big.Float).Mul(f, factor).Float64()
	return result
}

// CreateBookingRequest represents the request body for creating a new booking.
// It contains the event ID, ticket selections, and optional coupon code.
type CreateBookingRequest struct {
	EventID    string             `json:"event_id"`
	Tickets    []TicketSelection  `json:"tickets"`
	CouponCode string             `json:"coupon_code,omitempty"`
}

// TicketSelection represents a single ticket type selection in a booking request.
// It contains the ticket type ID and the quantity requested.
type TicketSelection struct {
	TicketTypeID string `json:"ticket_type_id"`
	Quantity     int    `json:"quantity"`
}

// Validate validates the CreateBookingRequest.
// It checks that required fields are present and valid.
func (r *CreateBookingRequest) Validate() error {
	// Validate event ID is a valid UUID
	if r.EventID == "" {
		return &ValidationError{Field: "event_id", Message: "event_id is required"}
	}

	// Validate at least one ticket selection
	if len(r.Tickets) == 0 {
		return &ValidationError{Field: "tickets", Message: "at least one ticket selection is required"}
	}

	// Validate each ticket selection
	for i, ticket := range r.Tickets {
		if ticket.TicketTypeID == "" {
			return &ValidationError{Field: fmt.Sprintf("tickets[%d].ticket_type_id", i), Message: "ticket_type_id is required"}
		}
		if ticket.Quantity <= 0 {
			return &ValidationError{Field: fmt.Sprintf("tickets[%d].quantity", i), Message: "quantity must be greater than 0"}
		}
	}

	return nil
}

// BookingResponse represents the response for a booking operation.
// It contains the booking details along with related information.
type BookingResponse struct {
	Booking      repository.Booking      `json:"booking"`
	Tickets      []BookingTicketResponse `json:"tickets,omitempty"`
	Event        *EventInfo              `json:"event,omitempty"`
	TotalAmount  float64                 `json:"total_amount"`
	Discount     float64                 `json:"discount"`
	FinalAmount  float64                 `json:"final_amount"`
}

// BookingTicketResponse represents a ticket within a booking response.
// It contains the ticket details along with the ticket type information.
type BookingTicketResponse struct {
	BookingTicket repository.BookingTicket `json:"booking_ticket"`
	TicketType    TicketTypeInfo           `json:"ticket_type"`
}

// TicketTypeInfo contains summary information about a ticket type.
// It's used in booking responses to provide context about the tickets.
type TicketTypeInfo struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
}

// EventInfo contains summary information about an event.
// It's used in booking responses to provide context about the event.
type EventInfo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Location  string `json:"location"`
}

// AvailabilityResponse represents the response for ticket availability check.
// It contains the availability information for each ticket type.
type AvailabilityResponse struct {
	EventID       string               `json:"event_id"`
	EventTitle    string               `json:"event_title"`
	TicketTypes   []TicketAvailability `json:"ticket_types"`
	TotalSeats    int                  `json:"total_seats"`
	AvailableSeats int                 `json:"available_seats"`
}

// TicketAvailability represents the availability for a single ticket type.
// It contains the ticket type details and current availability.
type TicketAvailability struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	TotalAvailable  int     `json:"total_available"`
	Remaining       int     `json:"remaining"`
	MaxPerUser      int     `json:"max_per_user"`
}

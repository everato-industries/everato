// Package event provides services for event management in the Everato platform.
// It handles the creation, updating, deletion, and querying of events.
package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

// Location constants define the valid location types for events
// These are used to validate location input and provide consistent options
const (
	LocationOnline   = "online"    // Online/virtual events held remotely
	LocationInPerson = "in-person" // Physical events held at a specific venue
)

// stringToPgText converts a string to pgtype.Text, handling empty strings
func stringToPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	text, _ := utils.StringToText(s)
	return text
}

// TicketTypeDTO represents the data for creating a ticket type
type TicketTypeDTO struct {
	Name             string  `json:"name" validate:"required,min=2,max=100"`                // Ticket type name (e.g., "VIP", "General", "Student")
	Price            float64 `json:"price" validate:"required,min=0,max=100000"`            // Ticket price
	AvailableTickets int     `json:"available_tickets" validate:"required,min=1,max=10000"` // Number of tickets available for this type
}

// CouponDTO represents the data for creating a coupon
type CouponDTO struct {
	Code               string  `json:"code" validate:"required,min=3,max=50"`                 // Coupon code (e.g., "EARLY20", "STUDENT")
	DiscountPercentage float64 `json:"discount_percentage" validate:"required,min=1,max=100"` // Discount percentage (1-100%)
	ValidFrom          string  `json:"valid_from" validate:"required,datetime"`               // When coupon becomes valid (ISO 8601)
	ValidUntil         string  `json:"valid_until" validate:"required,datetime"`              // When coupon expires (ISO 8601)
	UsageLimit         int     `json:"usage_limit" validate:"required,min=1,max=10000"`       // Maximum number of times coupon can be used
}

// CreateEventDTO represents the data transfer object for creating a new event.
// It handles the JSON deserialization and validation of event creation requests.
//
// This DTO provides several key functions:
//   - Defines the structure for incoming JSON event data
//   - Specifies validation rules for each field
//   - Provides methods to convert the validated data to repository formats
//   - Transforms string representations to appropriate database types (UUIDs, timestamps)
//   - Includes ticket types and coupons for comprehensive event creation
type CreateEventDTO struct {
	Title          string          `json:"title" validate:"required,min=2,max=100"`                               // Event title (2-100 chars)
	Description    string          `json:"description" validate:"required,min=10,max=500"`                        // Event description (10-500 chars)
	StartTime      string          `json:"start_time" validate:"required,datetime"`                               // Event start time in ISO 8601 format
	EndTime        string          `json:"end_time" validate:"required,datetime"`                                 // Event end time in ISO 8601 format
	Location       string          `json:"location" validate:"required,min=2,max=100"`                            // Event location (online or in-person)
	AdminID        string          `json:"admin_id" validate:"required,uuid"`                                     // UUID of the event administrator
	BannerURL      string          `json:"banner_url" validate:"omitempty,url"`                                   // Optional URL to event banner image
	IconURL        string          `json:"icon_url" validate:"omitempty,url"`                                     // Optional URL to event icon image
	TotalSeats     int             `json:"total_seats" validate:"required,min=1,max=10000"`                       // Total capacity of the event (1-10000)
	AvailableSeats int             `json:"available_seats" validate:"required,min=0,max=10000"`                   // Initially available seats for booking
	Status         string          `json:"status" validate:"omitempty,oneof=CREATED STARTED COMPLETED CANCELLED"` // Event status (defaults to CREATED)
	TicketTypes    []TicketTypeDTO `json:"ticket_types" validate:"required,min=1,dive"`                           // Ticket types for the event (at least 1 required)
	Coupons        []CouponDTO     `json:"coupons" validate:"omitempty,dive"`                                     // Optional coupons for the event

	// Organizer Information
	OrganizerName  string `json:"organizer_name" validate:"omitempty,min=2,max=255"`  // Name of the event organizer
	OrganizerEmail string `json:"organizer_email" validate:"omitempty,email,max=255"` // Email of the event organizer
	OrganizerPhone string `json:"organizer_phone" validate:"omitempty,max=50"`        // Phone number of the event organizer
	Organization   string `json:"organization" validate:"omitempty,max=255"`          // Organization hosting the event

	// Contact Information
	ContactEmail string `json:"contact_email" validate:"omitempty,email,max=255"` // Contact email for event inquiries
	ContactPhone string `json:"contact_phone" validate:"omitempty,max=50"`        // Contact phone for event inquiries

	// Event Details
	EventType          string `json:"event_type" validate:"omitempty,oneof=CONFERENCE WORKSHOP SEMINAR MEETUP FESTIVAL CONCERT EXHIBITION OTHER"`          // Type of event
	Category           string `json:"category" validate:"omitempty,oneof=TECHNOLOGY BUSINESS EDUCATION HEALTH ARTS SPORTS ENTERTAINMENT NETWORKING OTHER"` // Event category
	MaxTicketsPerUser  int    `json:"max_tickets_per_user" validate:"omitempty,min=1,max=100"`                                                             // Maximum tickets per user
	BookingStartTime   string `json:"booking_start_time" validate:"omitempty,datetime"`                                                                    // When booking opens
	BookingEndTime     string `json:"booking_end_time" validate:"omitempty,datetime"`                                                                      // When booking closes
	RefundPolicy       string `json:"refund_policy" validate:"omitempty,max=2000"`                                                                         // Event refund policy
	TermsAndConditions string `json:"terms_and_conditions" validate:"omitempty,max=5000"`                                                                  // Event terms and conditions
	Tags               string `json:"tags" validate:"omitempty,max=500"`                                                                                   // Comma-separated tags

	// Social Links
	WebsiteURL   string `json:"website_url" validate:"omitempty,url,max=500"`   // Event website URL
	FacebookURL  string `json:"facebook_url" validate:"omitempty,url,max=500"`  // Facebook page URL
	TwitterURL   string `json:"twitter_url" validate:"omitempty,url,max=500"`   // Twitter profile URL
	InstagramURL string `json:"instagram_url" validate:"omitempty,url,max=500"` // Instagram profile URL
	LinkedinURL  string `json:"linkedin_url" validate:"omitempty,url,max=500"`  // LinkedIn profile URL

	// Venue Details
	VenueName    string  `json:"venue_name" validate:"omitempty,max=255"`         // Name of the venue
	AddressLine1 string  `json:"address_line1" validate:"omitempty,max=255"`      // Primary address line
	AddressLine2 string  `json:"address_line2" validate:"omitempty,max=255"`      // Secondary address line
	City         string  `json:"city" validate:"omitempty,max=100"`               // City
	State        string  `json:"state" validate:"omitempty,max=100"`              // State/Province
	PostalCode   string  `json:"postal_code" validate:"omitempty,max=20"`         // Postal/ZIP code
	Country      string  `json:"country" validate:"omitempty,max=100"`            // Country
	Latitude     float64 `json:"latitude" validate:"omitempty,min=-90,max=90"`    // Latitude coordinate
	Longitude    float64 `json:"longitude" validate:"omitempty,min=-180,max=180"` // Longitude coordinate

	// Internal parsed time fields
	TIME_Start        time.Time `json:"-" validate:"-"` // Internal field for parsed start time (not exposed in JSON)
	TIME_End          time.Time `json:"-" validate:"-"` // Internal field for parsed end time (not exposed in JSON)
	TIME_BookingStart time.Time `json:"-" validate:"-"` // Internal field for parsed booking start time
	TIME_BookingEnd   time.Time `json:"-" validate:"-"` // Internal field for parsed booking end time
}

// time_parser is a custom validator function for validating datetime strings.
// It ensures that date/time fields conform to RFC3339/ISO8601 format.
//
// Parameters:
//   - fl: The validator's field level with access to the field being validated
//
// Returns:
//   - true if the string can be parsed as a valid RFC3339 timestamp, false otherwise
func time_parser(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}

// uuid_parser is a custom validator function for validating UUID strings.
// It ensures that UUID fields contain valid UUID format strings.
//
// Parameters:
//   - fl: The validator's field level with access to the field being validated
//
// Returns:
//   - true if the string can be parsed as a valid UUID, false otherwise
func uuid_parser(fl validator.FieldLevel) bool {
	_, err := utils.StringToUUID(fl.Field().String())
	return err == nil
}

// Validate checks the CreateEventDTO for required fields and formats.
// It registers and applies custom validators for datetime and UUID fields,
// then performs full validation of the DTO structure.
//
// The validation ensures:
// - All required fields are present
// - String lengths are within specified bounds
// - Dates are in proper RFC3339 format
// - UUIDs are valid
// - URLs are properly formatted
// - Numerical ranges are appropriate
//
// Returns:
//   - nil if validation passes
//   - A validation error detailing which fields failed validation and why
func (c *CreateEventDTO) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	_ = v.RegisterValidation("datetime", time_parser) // Register custom datetime validator
	_ = v.RegisterValidation("uuid", uuid_parser)     // Register custom UUID validator

	// Validate main event datetime fields
	st, err := time.Parse(time.RFC3339, c.StartTime)
	if err != nil {
		return fmt.Errorf("invalid start_time format: %w", err)
	}
	et, err := time.Parse(time.RFC3339, c.EndTime)
	if err != nil {
		return fmt.Errorf("invalid end_time format: %w", err)
	}

	// Validate coupon datetime fields
	for i, coupon := range c.Coupons {
		validFrom, err := time.Parse(time.RFC3339, coupon.ValidFrom)
		if err != nil {
			return fmt.Errorf("invalid valid_from format in coupon %d: %w", i, err)
		}
		validUntil, err := time.Parse(time.RFC3339, coupon.ValidUntil)
		if err != nil {
			return fmt.Errorf("invalid valid_until format in coupon %d: %w", i, err)
		}

		// Ensure coupon validity period makes sense
		if validUntil.Before(validFrom) {
			return fmt.Errorf("coupon %d: valid_until must be after valid_from", i)
		}
	}

	// Validate that ticket types don't exceed total seats
	totalTicketCapacity := 0
	for _, ticketType := range c.TicketTypes {
		totalTicketCapacity += ticketType.AvailableTickets
	}
	if totalTicketCapacity > c.TotalSeats {
		return fmt.Errorf("total ticket capacity (%d) exceeds event total seats (%d)", totalTicketCapacity, c.TotalSeats)
	}

	// Validate booking datetime fields if provided
	if c.BookingStartTime != "" {
		bst, err := time.Parse(time.RFC3339, c.BookingStartTime)
		if err != nil {
			return fmt.Errorf("invalid booking_start_time format: %w", err)
		}
		c.TIME_BookingStart = bst
	}

	if c.BookingEndTime != "" {
		bet, err := time.Parse(time.RFC3339, c.BookingEndTime)
		if err != nil {
			return fmt.Errorf("invalid booking_end_time format: %w", err)
		}
		c.TIME_BookingEnd = bet

		// Ensure booking end is after booking start if both are provided
		if c.BookingStartTime != "" && bet.Before(c.TIME_BookingStart) {
			return fmt.Errorf("booking_end_time must be after booking_start_time")
		}
	}

	c.TIME_Start = st // Set parsed start time for internal use
	c.TIME_End = et   // Set parsed end time for internal use
	return v.Struct(c)
}

// ToCreateEventParams converts the validated CreateEventDTO to a CreateEventParams struct
// for database operations. This method transforms the DTO's string and primitive types
// into the PostgreSQL-specific types required by the repository.
//
// The method performs the following conversions:
// - String timestamps to pgtype.Timestamptz
// - String location to pgtype.Text
// - String UUID to pgtype.UUID
// - Integer seat counts to int32 (PostgreSQL integer)
//
// Note: This method assumes validation has already been performed, so it ignores
// potential errors from the conversion functions.
//
// Returns:
//   - A repository.CreateEventParams struct ready for database insertion
func (c CreateEventDTO) ToCreateEventParams(slug string) (repository.CreateEventParams, error) {
	// Parse start time
	startTime, err := time.Parse(time.RFC3339, c.StartTime)
	if err != nil {
		return repository.CreateEventParams{}, fmt.Errorf("invalid start_time format: %w", err)
	}

	// Parse end time
	endTime, err := time.Parse(time.RFC3339, c.EndTime)
	if err != nil {
		return repository.CreateEventParams{}, fmt.Errorf("invalid end_time format: %w", err)
	}

	// Create properly initialized Timestamptz values with Valid=true
	st := pgtype.Timestamptz{Time: startTime, InfinityModifier: pgtype.Finite, Valid: true}
	et := pgtype.Timestamptz{Time: endTime, InfinityModifier: pgtype.Finite, Valid: true}

	location, err := utils.StringToText(c.Location)
	if err != nil {
		return repository.CreateEventParams{}, fmt.Errorf("invalid location format: %w", err)
	}

	adminUUID, err := utils.StringToUUID(c.AdminID)
	if err != nil {
		return repository.CreateEventParams{}, fmt.Errorf("invalid admin ID: %w", err)
	}

	// Set default status if not provided
	status := repository.EventStatus(c.Status)
	if status == "" {
		status = repository.EventStatusCREATED
	}

	// Handle optional booking times
	var bookingStartTime, bookingEndTime pgtype.Timestamptz
	if c.BookingStartTime != "" {
		bookingStartTime = pgtype.Timestamptz{Time: c.TIME_BookingStart, InfinityModifier: pgtype.Finite, Valid: true}
	}
	if c.BookingEndTime != "" {
		bookingEndTime = pgtype.Timestamptz{Time: c.TIME_BookingEnd, InfinityModifier: pgtype.Finite, Valid: true}
	}

	// Convert tags string to array
	var tagsArray []string
	if c.Tags != "" {
		tagsArray = strings.Split(c.Tags, ",")
		for i, tag := range tagsArray {
			tagsArray[i] = strings.TrimSpace(tag)
		}
	}

	// Convert coordinates to pgtype.Numeric
	var latitude, longitude pgtype.Numeric
	if c.Latitude != 0 {
		latitude = pgtype.Numeric{Valid: true}
		latitude.Scan(fmt.Sprintf("%.8f", c.Latitude))
	}
	if c.Longitude != 0 {
		longitude = pgtype.Numeric{Valid: true}
		longitude.Scan(fmt.Sprintf("%.8f", c.Longitude))
	}

	return repository.CreateEventParams{
		Title:              c.Title,
		Description:        c.Description,
		Slug:               slug,
		StartTime:          st,
		EndTime:            et,
		Location:           location,
		AdminID:            adminUUID,
		Banner:             c.BannerURL,
		Icon:               c.IconURL,
		TotalSeats:         int32(c.TotalSeats),
		AvailableSeats:     int32(c.AvailableSeats),
		Status:             status,
		OrganizerName:      stringToPgText(c.OrganizerName),
		OrganizerEmail:     stringToPgText(c.OrganizerEmail),
		OrganizerPhone:     stringToPgText(c.OrganizerPhone),
		Organization:       stringToPgText(c.Organization),
		ContactEmail:       stringToPgText(c.ContactEmail),
		ContactPhone:       stringToPgText(c.ContactPhone),
		RefundPolicy:       stringToPgText(c.RefundPolicy),
		TermsAndConditions: stringToPgText(c.TermsAndConditions),
		EventType:          stringToPgText(c.EventType),
		Category:           stringToPgText(c.Category),
		MaxTicketsPerUser:  pgtype.Int4{Int32: int32(c.MaxTicketsPerUser), Valid: c.MaxTicketsPerUser > 0},
		BookingStartTime:   bookingStartTime,
		BookingEndTime:     bookingEndTime,
		Tags:               tagsArray,
		WebsiteUrl:         stringToPgText(c.WebsiteURL),
		FacebookUrl:        stringToPgText(c.FacebookURL),
		TwitterUrl:         stringToPgText(c.TwitterURL),
		InstagramUrl:       stringToPgText(c.InstagramURL),
		LinkedinUrl:        stringToPgText(c.LinkedinURL),
		VenueName:          stringToPgText(c.VenueName),
		AddressLine1:       stringToPgText(c.AddressLine1),
		AddressLine2:       stringToPgText(c.AddressLine2),
		City:               stringToPgText(c.City),
		State:              stringToPgText(c.State),
		PostalCode:         stringToPgText(c.PostalCode),
		Country:            stringToPgText(c.Country),
		Latitude:           latitude,
		Longitude:          longitude,
	}, nil
}

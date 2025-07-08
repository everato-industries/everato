// Package event provides services for event management in the Everato platform.
// It handles the creation, updating, deletion, and querying of events.
package event

import (
	"time"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/go-playground/validator/v10"
)

// Location constants define the valid location types for events
// These are used to validate location input and provide consistent options
const (
	LocationOnline   = "online"    // Online/virtual events held remotely
	LocationInPerson = "in-person" // Physical events held at a specific venue
)

// CreateEventDTO represents the data transfer object for creating a new event.
// It handles the JSON deserialization and validation of event creation requests.
//
// This DTO provides several key functions:
//   - Defines the structure for incoming JSON event data
//   - Specifies validation rules for each field
//   - Provides methods to convert the validated data to repository formats
//   - Transforms string representations to appropriate database types (UUIDs, timestamps)
type CreateEventDTO struct {
	Title          string `json:"title" validate:"required,min=2,max=100"`             // Event title (2-100 chars)
	Description    string `json:"description" validate:"required,min=10,max=500"`      // Event description (10-500 chars)
	StartTime      string `json:"start_time" validate:"required,datetime"`             // Event start time in ISO 8601 format
	EndTime        string `json:"end_time" validate:"required,datetime"`               // Event end time in ISO 8601 format
	Location       string `json:"location" validate:"required,min=2,max=100"`          // Event location (online or in-person)
	AdminID        string `json:"admin_id" validate:"required,uuid"`                   // UUID of the event administrator
	BannerURL      string `json:"banner_url" validate:"omitempty,url"`                 // Optional URL to event banner image
	IconURL        string `json:"icon_url" validate:"omitempty,url"`                   // Optional URL to event icon image
	TotalSeats     int    `json:"total_seats" validate:"required,min=1,max=10000"`     // Total capacity of the event (1-10000)
	AvailableSeats int    `json:"available_seats" validate:"required,min=0,max=10000"` // Initially available seats for booking
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
func (c CreateEventDTO) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	_ = v.RegisterValidation("datetime", time_parser) // Register custom datetime validator
	_ = v.RegisterValidation("uuid", uuid_parser)     // Register custom UUID validator
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
func (c CreateEventDTO) ToCreateEventParams() repository.CreateEventParams {
	start_time, _ := utils.StringToTime(c.StartTime) // Convert ISO8601 string to PostgreSQL timestamp
	end_time, _ := utils.StringToTime(c.EndTime)     // Convert ISO8601 string to PostgreSQL timestamp
	location, _ := utils.StringToText(c.Location)    // Convert string to PostgreSQL text
	adminUUID, _ := utils.StringToUUID(c.AdminID)    // Convert string to PostgreSQL UUID

	return repository.CreateEventParams{
		Title:          c.Title,
		Description:    c.Description,
		StartTime:      start_time,
		EndTime:        end_time,
		Location:       location,
		AdminID:        adminUUID,
		Banner:         c.BannerURL,
		Icon:           c.IconURL,
		TotalSeats:     int32(c.TotalSeats),
		AvailableSeats: int32(c.AvailableSeats),
	}
}

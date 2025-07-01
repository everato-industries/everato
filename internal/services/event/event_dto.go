package event

import (
	"time"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/go-playground/validator/v10"
)

// Location ENUM provides the valid options for creating an event
const (
	LocationOnline   = "online"    // ONLINE, i.e. event organisation is TODO
	LocationInPerson = "in-person" // IN-PERSON
)

// CreateEventDTO represents the data transfer object for creating a new event
//
// This will also include some additional functionalities
// i.e.
//   - parsing the raw request body to the DTO object it self
//   - transforming the string UUID to a pgtype.UUID format
//
// and etc.
type CreateEventDTO struct {
	Title          string `json:"title" validate:"required,min=2,max=100"`             // Name of the event
	Description    string `json:"description" validate:"required,min=10,max=500"`      // Description of the event
	StartTime      string `json:"start_time" validate:"required,datetime"`             // Start time of the event, ISO 8061 format
	EndTime        string `json:"end_time" validate:"required,datetime"`               // End time of the event, ISO 8061 format
	Location       string `json:"location" validate:"required,min=2,max=100"`          // Location of the event, i.e it can be online or in-person
	AdminID        string `json:"admin_id" validate:"required,uuid"`                   // UUID format for admin ID
	BannerURL      string `json:"banner_url" validate:"omitempty,url"`                 // Optional URL for event banner
	IconURL        string `json:"icon_url" validate:"omitempty,url"`                   // Optional URL for event icon
	TotalSeats     int    `json:"total_seats" validate:"required,min=1,max=10000"`     // Total number of seats available for the event
	AvailableSeats int    `json:"available_seats" validate:"required,min=0,max=10000"` // Number of seats available for booking
}

// Custom datetime validator
func time_parser(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}

// Custome UUID validator
func uuid_parser(fl validator.FieldLevel) bool {
	_, err := utils.StringToUUID(fl.Field().String())
	return err == nil
}

// Validate checks the CreateEventDTO for required fields and formats.
// it returns an error if there is some issue with the data
func (c CreateEventDTO) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	_ = v.RegisterValidation("datetime", time_parser) // Register custom datetime validator
	_ = v.RegisterValidation("uuid", uuid_parser)     // Register custom UUID validator
	return v.Struct(c)
}

// ToCreateEventParams converts the CreateEventDTO to a CreateEventParams for database operations.
func (c CreateEventDTO) ToCreateEventParams() repository.CreateEventParams {

	start_time, _ := utils.StringToTime(c.StartTime) // Parse to time
	end_time, _ := utils.StringToTime(c.EndTime)     // Parse to time
	location, _ := utils.StringToText(c.Location)    // Parse to text
	adminUUID, _ := utils.StringToUUID(c.AdminID)    // Parse to UUID

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

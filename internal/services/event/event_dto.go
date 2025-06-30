package event

import "github.com/go-playground/validator/v10"

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
	Title          string `json:"title" validate:"required,min=2,max=100"`                           // Name of the event
	Description    string `json:"description" validate:"required,min=10,max=500"`                    // Description of the event
	StartTime      string `json:"start_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"` // Start time of the event, ISO 8061 format
	EndTime        string `json:"end_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`   // End time of the event, ISO 8061 format
	Location       string `json:"location" validate:"required,min=2,max=100"`                        // Location of the event, i.e it can be online or in-person
	AdminID        string `json:"admin_id" validate:"required,uuid"`                                 // UUID format for admin ID
	BannerURL      string `json:"banner_url" validate:"omitempty,url"`                               // Optional URL for event banner
	IconURL        string `json:"icon_url" validate:"omitempty,url"`                                 // Optional URL for event icon
	TotalSeats     int    `json:"total_seats" validate:"required,min=1,max=10000"`                   // Total number of seats available for the event
	AvailableSeats int    `json:"available_seats" validate:"required,min=0,max=10000"`               // Number of seats available for booking
}

// Validate checks the CreateEventDTO for required fields and formats.
// it returns an error if there is some issue with the data
func (c CreateEventDTO) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	return v.Struct(c)
}

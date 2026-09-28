package booking

import (
	"net/http"

	"github.com/jackc/pgx/v5"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// CheckAvailability handles the check for ticket availability for an event.
//
// This function performs the following operations:
// 1. Extracts the event ID from the URL path
// 2. Validates the event ID format
// 3. Retrieves the event and ticket availability from the database
// 4. Returns the availability information or appropriate error responses
//
// Parameters:
//   - wr: Custom HTTP writer for response handling
//   - repo: Database repository for booking operations
//   - conn: Database connection for transaction management
func CheckAvailability(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Extract event ID from URL path
	eventID := utils.GetParam(wr.R, "eventId")
	if eventID == "" {
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Event ID is required.",
			},
		)
		return
	}

	// Parse event ID to UUID
	eventUUID, err := utils.StringToUUID(eventID)
	if err != nil {
		logger.StdoutLogger.Error("Invalid event UUID", "eventID", eventID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid event ID format.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Get event from database
	event, err := repo.GetEventByID(wr.R.Context(), eventUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			wr.Status(http.StatusNotFound).Json(
				utils.M{
					"message": "Event not found.",
				},
			)
			return
		}
		logger.StdoutLogger.Error("Failed to get event", "eventID", eventID, "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to retrieve event.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Get ticket availability for the event
	ticketTypes, err := repo.GetEventTicketAvailability(wr.R.Context(), eventUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No ticket types found
			ticketTypes = []repository.GetEventTicketAvailabilityRow{}
		} else {
			logger.StdoutLogger.Error("Failed to get ticket availability", "eventID", eventID, "err", err.Error())
			wr.Status(http.StatusInternalServerError).Json(
				utils.M{
					"message": "Failed to retrieve ticket availability.",
					"err":     err.Error(),
				},
			)
			return
		}
	}

	// Build availability response
	var ticketAvailability []TicketAvailability
	for _, tt := range ticketTypes {
		ticketAvailability = append(ticketAvailability, TicketAvailability{
			ID:             tt.ID.String(),
			Name:           tt.Name,
			Price:          tt.Price,
			TotalAvailable: int(tt.TotalAvailable),
			Remaining:      int(tt.Remaining),
			MaxPerUser:     int(event.MaxTicketsPerUser.Int32),
		})
	}

	response := AvailabilityResponse{
		EventID:        event.ID.String(),
		EventTitle:     event.Title,
		TicketTypes:    ticketAvailability,
		TotalSeats:     int(event.TotalSeats),
		AvailableSeats: int(event.AvailableSeats),
	}

	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Availability retrieved successfully.",
			"data":    response,
		},
	)
}

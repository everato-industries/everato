package booking

import (
	"net/http"

	"github.com/jackc/pgx/v5"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// GetUserBookings handles the retrieval of all bookings for a specific user.
//
// This function performs the following operations:
// 1. Extracts the user ID from the URL path
// 2. Validates the user ID format
// 3. Retrieves all bookings for the user from the database
// 4. Returns the list of bookings or appropriate error responses
//
// Parameters:
//   - wr: Custom HTTP writer for response handling
//   - repo: Database repository for booking operations
//   - conn: Database connection for transaction management
func GetUserBookings(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Extract user ID from URL path
	userID := utils.GetParam(wr.R, "userId")
	if userID == "" {
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "User ID is required.",
			},
		)
		return
	}

	// Parse user ID to UUID
	userUUID, err := utils.StringToUUID(userID)
	if err != nil {
		logger.StdoutLogger.Error("Invalid user UUID", "userID", userID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid user ID format.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Get user bookings from database
	bookings, err := repo.GetUserBookings(wr.R.Context(), userUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			wr.Status(http.StatusOK).Json(
				utils.M{
					"message": "No bookings found.",
					"data":    []repository.Booking{},
				},
			)
			return
		}
		logger.StdoutLogger.Error("Failed to get user bookings", "userID", userID, "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to retrieve bookings.",
				"err":     err.Error(),
			},
		)
		return
	}

	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Bookings retrieved successfully.",
			"data":    bookings,
			"count":   len(bookings),
		},
	)
}

// GetBookingDetails handles the retrieval of details for a specific booking.
//
// This function performs the following operations:
// 1. Extracts the booking ID from the URL path
// 2. Validates the booking ID format
// 3. Retrieves the booking and associated tickets from the database
// 4. Returns the booking details or appropriate error responses
//
// Parameters:
//   - wr: Custom HTTP writer for response handling
//   - repo: Database repository for booking operations
//   - conn: Database connection for transaction management
func GetBookingDetails(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Extract booking ID from URL path
	bookingID := utils.GetParam(wr.R, "bookingId")
	if bookingID == "" {
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Booking ID is required.",
			},
		)
		return
	}

	// Parse booking ID to UUID
	bookingUUID, err := utils.StringToUUID(bookingID)
	if err != nil {
		logger.StdoutLogger.Error("Invalid booking UUID", "bookingID", bookingID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid booking ID format.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Get booking from database
	booking, err := repo.GetBookingByID(wr.R.Context(), bookingUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			wr.Status(http.StatusNotFound).Json(
				utils.M{
					"message": "Booking not found.",
				},
			)
			return
		}
		logger.StdoutLogger.Error("Failed to get booking", "bookingID", bookingID, "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to retrieve booking.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Get booking tickets
	tickets, err := repo.GetBookingTicketsByBookingID(wr.R.Context(), bookingUUID)
	if err != nil && err != pgx.ErrNoRows {
		logger.StdoutLogger.Error("Failed to get booking tickets", "bookingID", bookingID, "err", err.Error())
		// Continue without tickets, don't fail the request
	}

	// Build response
	response := BookingResponse{
		Booking:     booking,
		TotalAmount: numericToFloat(booking.TotalAmount),
		Discount:    numericToFloat(booking.DiscountAmount),
		FinalAmount: numericToFloat(booking.FinalAmount),
	}

	// Add ticket details to response
	for _, ticket := range tickets {
		response.Tickets = append(response.Tickets, BookingTicketResponse{
			BookingTicket: repository.BookingTicket{
				ID:           ticket.ID,
				BookingID:    ticket.BookingID,
				TicketTypeID: ticket.TicketTypeID,
				Quantity:     ticket.Quantity,
				PricePerUnit: ticket.PricePerUnit,
				Subtotal:     ticket.Subtotal,
				CreatedAt:    ticket.CreatedAt,
			},
			TicketType: TicketTypeInfo{
				ID:    ticket.TicketTypeID.String(),
				Name:  ticket.TicketTypeName,
				Price: numericToFloat(ticket.PricePerUnit),
			},
		})
	}

	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Booking retrieved successfully.",
			"data":    response,
		},
	)
}

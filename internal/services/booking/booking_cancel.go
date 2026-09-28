package booking

import (
	"net/http"

	"github.com/jackc/pgx/v5"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// CancelBooking handles the cancellation of a booking.
//
// This function performs the following operations:
// 1. Extracts the booking ID from the URL path
// 2. Validates the booking ID format
// 3. Retrieves the booking and associated tickets
// 4. Cancels the booking in a transaction
// 5. Restores ticket availability
// 6. Returns the cancellation confirmation or appropriate error responses
//
// Parameters:
//   - wr: Custom HTTP writer for response handling
//   - repo: Database repository for booking operations
//   - conn: Database connection for transaction management
func CancelBooking(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
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

	// Get the user ID from the context (set by auth middleware)
	userID, ok := wr.R.Context().Value("user_id").(string)
	if !ok {
		logger.StdoutLogger.Error("Failed to get user ID from context")
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication required.",
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

	// Start a database transaction
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Failed to begin transaction", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error.",
				"err":     err.Error(),
			},
		)
		return
	}
	defer tx.Rollback(wr.R.Context())

	// Get booking from database
	booking, err := repo.WithTx(tx).GetBookingByID(wr.R.Context(), bookingUUID)
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

	// Verify the user owns this booking
	if booking.UserID != userUUID {
		wr.Status(http.StatusForbidden).Json(
			utils.M{
				"message": "You are not authorized to cancel this booking.",
			},
		)
		return
	}

	// Check if booking can be cancelled (only PENDING or PENDING_PAYMENT bookings)
	if booking.Status != "PENDING" && booking.Status != "PENDING_PAYMENT" {
		wr.Status(http.StatusConflict).Json(
			utils.M{
				"message": "Booking cannot be cancelled in its current status.",
				"status":  booking.Status,
			},
		)
		return
	}

	// Get booking tickets to restore availability
	tickets, err := repo.WithTx(tx).GetBookingTicketsByBookingID(wr.R.Context(), bookingUUID)
	if err != nil && err != pgx.ErrNoRows {
		logger.StdoutLogger.Error("Failed to get booking tickets", "bookingID", bookingID, "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to retrieve booking tickets.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Restore ticket availability
	for _, ticket := range tickets {
		_, err := repo.WithTx(tx).IncrementTicketAvailability(wr.R.Context(), repository.IncrementTicketAvailabilityParams{
			ID:               ticket.TicketTypeID,
			AvailableTickets: ticket.Quantity,
		})
		if err != nil {
			logger.StdoutLogger.Error("Failed to restore ticket availability", "err", err.Error())
			wr.Status(http.StatusInternalServerError).Json(
				utils.M{
					"message": "Failed to restore ticket availability.",
					"err":     err.Error(),
				},
			)
			return
		}
	}

	// Cancel the booking
	cancelledBooking, err := repo.WithTx(tx).CancelBooking(wr.R.Context(), bookingUUID)
	if err != nil {
		logger.StdoutLogger.Error("Failed to cancel booking", "bookingID", bookingID, "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to cancel booking.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Commit the transaction
	if err := tx.Commit(wr.R.Context()); err != nil {
		logger.StdoutLogger.Error("Failed to commit transaction", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to complete cancellation.",
				"err":     err.Error(),
			},
		)
		return
	}

	logger.StdoutLogger.Info("Booking cancelled successfully",
		"booking_id", booking.ID,
		"user_id", userID,
	)

	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Booking cancelled successfully.",
			"data":    cancelledBooking,
		},
	)
}

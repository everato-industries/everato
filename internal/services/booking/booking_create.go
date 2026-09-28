package booking

import (
	"math/big"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/services/ticket"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// uuidToPgtype converts google/uuid.UUID to pgtype.UUID
func uuidToPgtype(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

// generateUUID generates a new UUID and returns it as pgtype.UUID
func generateUUID() pgtype.UUID {
	return uuidToPgtype(uuid.New())
}

// floatToNumeric converts a float64 to pgtype.Numeric
func floatToNumeric(f float64) pgtype.Numeric {
	// Convert float to big.Int with appropriate exponent
	// For example, 19.99 becomes Int=1999, Exp=2
	intVal := new(big.Int)
	intVal.SetInt64(int64(f * 100))
	return pgtype.Numeric{
		Int:   intVal,
		Exp:   -2,
		Valid: true,
	}
}

// CreateBooking handles the creation of a new booking in the system.
//
// This function performs the following operations:
// 1. Parses and validates the booking request
// 2. Verifies the user exists and is authenticated
// 3. Checks ticket availability for the event
// 4. Validates the coupon if provided
// 5. Creates the booking in a transaction for data consistency
// 6. Returns the created booking data or appropriate error responses
//
// Parameters:
//   - wr: Custom HTTP writer for response handling
//   - repo: Database repository for booking operations
//   - conn: Database connection for transaction management
//   - cfg: Application configuration
func CreateBooking(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Parse the request body to the CreateBookingRequest
	request := &CreateBookingRequest{}
	if err := wr.ParseBody(request); err != nil {
		logger.StdoutLogger.Error("Failed to parse request body", "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Please send the proper data :(",
				"err":     err.Error(),
			},
		)
		return
	}

	// Validate the request data
	if err := request.Validate(); err != nil {
		logger.StdoutLogger.Error("Validation failed for CreateBookingRequest", "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Error parsing the provided data :(",
				"err":     err.Error(),
			},
		)
		return
	}

	// Start a database transaction to ensure ACID properties
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Failed to begin transaction", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error, please try again later.",
				"err":     err.Error(),
			},
		)
		return
	}
	defer tx.Rollback(wr.R.Context())

	// Get the user ID from the context (set by auth middleware)
	userID, ok := wr.R.Context().Value("uid").(string)
	if !ok {
		logger.StdoutLogger.Error("Failed to get user ID from context")
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication required, user ID not found in request context.",
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

	// Parse event ID to UUID
	eventUUID, err := utils.StringToUUID(request.EventID)
	if err != nil {
		logger.StdoutLogger.Error("Invalid event UUID", "eventID", request.EventID, "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid event ID format.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Get event details and validate it exists
	event, err := repo.WithTx(tx).GetEventByID(wr.R.Context(), eventUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			wr.Status(http.StatusNotFound).Json(
				utils.M{
					"message": "Event not found.",
				},
			)
			return
		}
		logger.StdoutLogger.Error("Failed to get event", "eventID", request.EventID, "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to retrieve event information.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Check if event is available for booking (status should be CREATED or STARTED)
	if event.Status != "CREATED" && event.Status != "STARTED" {
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Event is not available for booking.",
				"status":  event.Status,
			},
		)
		return
	}

	// Check user's max tickets per event limit
	userBookingCount, err := repo.WithTx(tx).GetUserBookingCountForEvent(wr.R.Context(), repository.GetUserBookingCountForEventParams{
		UserID:  userUUID,
		EventID: eventUUID,
	})
	if err != nil {
		logger.StdoutLogger.Error("Failed to get user booking count", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to check booking limit.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Calculate total tickets requested
	totalTicketsRequested := 0
	for _, ticket := range request.Tickets {
		totalTicketsRequested += ticket.Quantity
	}

	// Check if user exceeds max tickets per user limit
	maxTickets := int(event.MaxTicketsPerUser.Int32)
	if int(userBookingCount)+totalTicketsRequested > maxTickets {
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message":   "You have exceeded the maximum tickets per user limit.",
				"max":       maxTickets,
				"current":   userBookingCount,
				"requested": totalTicketsRequested,
			},
		)
		return
	}

	// Validate ticket availability and calculate totals
	var totalAmount float64
	type ticketInfo struct {
		uuid   pgtype.UUID
		qty    int32
		price  float64
		subtotal float64
	}
	var ticketInfos []ticketInfo

	for _, ticketSelection := range request.Tickets {
		// Parse ticket type ID
		ticketTypeUUID, err := utils.StringToUUID(ticketSelection.TicketTypeID)
		if err != nil {
			logger.StdoutLogger.Error("Invalid ticket type UUID", "ticketTypeID", ticketSelection.TicketTypeID, "err", err.Error())
			wr.Status(http.StatusBadRequest).Json(
				utils.M{
					"message": "Invalid ticket type ID format.",
					"err":     err.Error(),
				},
			)
			return
		}

		// Get ticket availability
		availability, err := repo.WithTx(tx).GetTicketAvailability(wr.R.Context(), ticketTypeUUID)
		if err != nil {
			if err == pgx.ErrNoRows {
				wr.Status(http.StatusNotFound).Json(
					utils.M{
						"message": "Ticket type not found.",
					},
				)
				return
			}
			logger.StdoutLogger.Error("Failed to get ticket availability", "err", err.Error())
			wr.Status(http.StatusInternalServerError).Json(
				utils.M{
					"message": "Failed to check ticket availability.",
					"err":     err.Error(),
				},
			)
			return
		}

		// Check if enough tickets are available
		if availability.Remaining < int32(ticketSelection.Quantity) {
			wr.Status(http.StatusConflict).Json(
				utils.M{
					"message":     "Not enough tickets available.",
					"ticket_type": availability.Name,
					"available":   availability.Remaining,
					"requested":   ticketSelection.Quantity,
				},
			)
			return
		}

		// Calculate subtotal for this ticket type
		subtotal := availability.Price * float64(ticketSelection.Quantity)
		totalAmount += subtotal

		// Decrement ticket availability
		_, err = repo.WithTx(tx).DecrementTicketAvailability(wr.R.Context(), repository.DecrementTicketAvailabilityParams{
			ID:               ticketTypeUUID,
			AvailableTickets: int32(ticketSelection.Quantity),
		})
		if err != nil {
			logger.StdoutLogger.Error("Failed to decrement ticket availability", "err", err.Error())
			wr.Status(http.StatusInternalServerError).Json(
				utils.M{
					"message": "Failed to reserve tickets.",
					"err":     err.Error(),
				},
			)
			return
		}

		ticketInfos = append(ticketInfos, ticketInfo{
			uuid:     ticketTypeUUID,
			qty:      int32(ticketSelection.Quantity),
			price:    availability.Price,
			subtotal: subtotal,
		})
	}

	// Calculate discount if coupon is provided
	var discountAmount float64
	var couponCode *string

	if request.CouponCode != "" {
		coupon, err := repo.WithTx(tx).GetCouponByCode(wr.R.Context(), repository.GetCouponByCodeParams{
			Code:    request.CouponCode,
			EventID: eventUUID,
		})
		if err == nil {
			// Check if coupon is valid (within date range and usage limit)
			now := pgtype.Timestamptz{Time: time.Now(), InfinityModifier: pgtype.Finite, Valid: true}
			if coupon.ValidFrom.Valid && coupon.ValidUntil.Valid &&
				now.Time.After(coupon.ValidFrom.Time) && now.Time.Before(coupon.ValidUntil.Time) &&
				coupon.UsageCount < coupon.UsageLimit {
				// Apply discount
				discountAmount = totalAmount * (float64(coupon.DiscountPercentage) / 100)
				couponCode = &request.CouponCode

				// Increment coupon usage
				_, err = repo.WithTx(tx).IncrementCouponUsage(wr.R.Context(), coupon.ID)
				if err != nil {
					logger.StdoutLogger.Error("Failed to increment coupon usage", "err", err.Error())
					// Continue without failing the booking
				}
			}
		}
	}

	// Convert couponCode to pgtype.Text
	var couponCodeText pgtype.Text
	if couponCode != nil {
		couponCodeText = pgtype.Text{String: *couponCode, Valid: true}
	}

	// Calculate final amount
	finalAmount := totalAmount - discountAmount
	if finalAmount < 0 {
		finalAmount = 0
	}

	// Create the booking first (without tickets)
	bookingID := generateUUID()
	booking, err := repo.WithTx(tx).CreateBooking(wr.R.Context(), repository.CreateBookingParams{
		ID:             bookingID,
		UserID:         userUUID,
		EventID:        eventUUID,
		TotalAmount:    floatToNumeric(totalAmount),
		Status:         "CONFIRMED",
		CouponCode:     couponCodeText,
		DiscountAmount: floatToNumeric(discountAmount),
		FinalAmount:    floatToNumeric(finalAmount),
	})
	if err != nil {
		logger.StdoutLogger.Error("Failed to create booking", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to create booking.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Now create the booking tickets with the actual booking ID
	var createdTickets []repository.BookingTicket
	for _, ti := range ticketInfos {
		bookingTicketID := generateUUID()
		bookingTicket, err := repo.WithTx(tx).CreateBookingTicket(wr.R.Context(), repository.CreateBookingTicketParams{
			ID:           bookingTicketID,
			BookingID:    booking.ID,
			TicketTypeID: ti.uuid,
			Quantity:     ti.qty,
			PricePerUnit: floatToNumeric(ti.price),
			Subtotal:     floatToNumeric(ti.subtotal),
		})
		if err != nil {
			logger.StdoutLogger.Error("Failed to create booking ticket", "err", err.Error())
			wr.Status(http.StatusInternalServerError).Json(
				utils.M{
					"message": "Failed to create booking ticket.",
					"err":     err.Error(),
				},
			)
			return
		}
		createdTickets = append(createdTickets, bookingTicket)

		// Generate individual QR tickets
		_, err = ticket.GenerateTicketsForBooking(wr.R.Context(), repo.WithTx(tx), ticket.GenerateTicketsRequest{
			BookingID:    booking.ID.String(),
			EventID:      event.ID.String(),
			EventSlug:    event.Slug,
			TicketTypeID: ti.uuid.String(),
			Quantity:     int(ti.qty),
			UserID:       userID,
			HMACSecret:   utils.GetEnv("JWT_SECRET", "SUPER_SECRET_KEY"),
		}, cfg)
		if err != nil {
			logger.StdoutLogger.Error("Failed to generate tickets for booking", "err", err.Error())
			wr.Status(http.StatusInternalServerError).Json(
				utils.M{
					"message": "Failed to generate tickets.",
					"err":     err.Error(),
				},
			)
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit(wr.R.Context()); err != nil {
		logger.StdoutLogger.Error("Failed to commit transaction", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to complete booking.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Build response
	response := BookingResponse{
		Booking:     booking,
		TotalAmount: totalAmount,
		Discount:    discountAmount,
		FinalAmount: finalAmount,
		Event: &EventInfo{
			ID:        event.ID.String(),
			Title:     event.Title,
			StartTime: event.StartTime.Time.Format("2006-01-02T15:04:05Z"),
			EndTime:   event.EndTime.Time.Format("2006-01-02T15:04:05Z"),
			Location:  event.Location.String,
		},
	}

	// Add ticket details to response
	for _, ticket := range createdTickets {
		ticketTypes, _ := repo.GetTicketTypesByEventID(wr.R.Context(), eventUUID)
		for _, tt := range ticketTypes {
			if tt.ID == ticket.TicketTypeID {
				response.Tickets = append(response.Tickets, BookingTicketResponse{
					BookingTicket: ticket,
					TicketType: TicketTypeInfo{
						ID:    tt.ID.String(),
						Name:  tt.Name,
						Price: tt.Price,
					},
				})
				break
			}
		}
	}

	logger.StdoutLogger.Info("Booking created successfully",
		"booking_id", booking.ID,
		"user_id", userID,
		"event_id", request.EventID,
		"total_amount", totalAmount,
		"final_amount", finalAmount,
	)

	wr.Status(http.StatusCreated).Json(
		utils.M{
			"message": "Booking created successfully!",
			"data":    response,
		},
	)
}

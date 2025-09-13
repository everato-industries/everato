// Package event provides services for event management in the Everato platform.
// It handles the creation, updating, deletion and querying of events.
package event

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// CreateEvent handles the creation of a new event in the system.
//
// This function performs the following operations:
// 1. Parses and validates the event data from the request body
// 2. Generates a unique slug for the event based on its title
// 3. Verifies the slug is unique in the database
// 4. Creates the event record in a transaction for data consistency
// 5. Returns the created event data or appropriate error responses
//
// Parameters:
//   - wr: Custom HTTP writer for response handling
//   - repo: Database repository for event operations
//   - conn: Database connection for transaction management
func CreateEvent(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Parse the request body to the CreateEventDTO
	eventDTO := &CreateEventDTO{}
	// Parse the request body into the CreateEventDTO struct
	// This extracts all the event details from the JSON request
	if err := wr.ParseBody(eventDTO); err != nil {
		logger.StdoutLogger.Error("Failed to parse request body", "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Please send the proper data :(",
				"err":     err.Error(),
			},
		)
		return
	}

	// Validate the event data against the defined rules
	// This ensures all required fields are present and formatted correctly
	if err := eventDTO.Validate(); err != nil {
		logger.StdoutLogger.Error("Validation failed for CreateEventDTO", "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Error parsing the provided data :(",
				"err":     err.Error(),
			},
		)
		return
	}

	// Start a database transaction to ensure ACID properties
	// This guarantees that all operations either complete successfully or fail without side effects
	tx, err := conn.Begin(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Failed to begin transaction", "err", err.Error())
		logger.FileLogger.Error("Failed to begin transaction", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error, please try again later.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Get the admin details from the context
	// then look if the admin has the right permissions or not
	adminID, ok := wr.R.Context().Value("admin_id").(string)
	if !ok {
		logger.StdoutLogger.Error("Failed to get admin ID from context")
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication required, admin ID not found in request context.",
				"err":     "Admin ID not found in request context",
			},
		)
		return
	}

	// Get the admin details from the db
	admin_uuid, err := utils.StringToUUID(adminID)
	if err != nil {
		logger.StdoutLogger.Error("Invalid admin UUID", "adminID", adminID, "err", err.Error())
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid admin ID format.",
				"err":     err.Error(),
			},
		)
		return
	}

	admin, err := repo.WithTx(tx).GetAdminById(wr.R.Context(), admin_uuid)
	if err != nil {
		logger.StdoutLogger.Error("Failed to retrieve admin data", "adminID", adminID, "err", err.Error())
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		wr.Status(http.StatusUnauthorized).Json(
			utils.M{
				"message": "Authentication failed, unable to retrieve admin data.",
				"err":     err.Error(),
			},
		)
		return
	}

	has_permission := false // Flag to check if the admin has permission to create events
	for perm := range admin.Permissions {
		if repository.Permissions(admin.Permissions[perm]) == repository.PermissionsCREATEEVENT {
			has_permission = true // Admin has permission to create events
			break
		}
	}

	if !has_permission {
		logger.StdoutLogger.Error("Admin does not have permission to create events", "adminID", adminID)
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		wr.Status(http.StatusForbidden).Json(
			utils.M{
				"message": "You do not have permission to create events.",
				"err":     "Permission denied",
			},
		)
		return
	}

	// Generate a URL-friendly slug from the event title
	// This will be used in event URLs and must be unique across all events
	slug, err := utils.GenerateSlug(eventDTO.Title)
	if err != nil {
		logger.StdoutLogger.Error("Failed to generate slug", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error, please try again later.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Check if the generated slug already exists in the database
	// If it does, we need to return an error as slugs must be unique
	if _, err := repo.WithTx(tx).SearchSlug(wr.R.Context(), slug); err == nil {
		logger.StdoutLogger.Error("Slug already exists", "slug", slug)
		wr.Status(http.StatusConflict).Json(
			utils.M{
				"message": "Slug already exists, please try a different one.",
				"slug":    slug,
			},
		)
		return
	}

	// // Handle case where the slug search failed for a reason other than "not found"
	// // This indicates a database error rather than a slug uniqueness issue
	// if err != pgx.ErrNoRows {
	// 	logger.StdoutLogger.Error("Failed to create event", "err", err.Error())
	// 	wr.Status(http.StatusInternalServerError).Json(
	// 		utils.M{
	// 			"message": "Internal server error, please try again later.",
	// 			"err":     err.Error(),
	// 		},
	// 	)
	// 	return
	// }

	// Create the event in the database using the transaction
	// This converts our DTO to the format required by the repository
	params, err := eventDTO.ToCreateEventParams(slug)
	if err != nil {
		logger.StdoutLogger.Error("Failed to prepare event parameters", "err", err.Error())
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid event data format",
				"err":     err.Error(),
			},
		)
		return
	}

	// Create the main event first
	event, err := repo.WithTx(tx).CreateEvent(wr.R.Context(), params)
	if err != nil {
		logger.StdoutLogger.Error("Failed to create event in database", "err", err.Error())
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Failed to create event, please try again later.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Create ticket types for the event
	var createdTicketTypes []repository.TicketType
	for i, ticketTypeDTO := range eventDTO.TicketTypes {
		ticketTypeParams := repository.CreateTicketTypeParams{
			Name:             ticketTypeDTO.Name,
			EventID:          event.ID,
			Price:            ticketTypeDTO.Price,
			AvailableTickets: int32(ticketTypeDTO.AvailableTickets),
		}

		ticketType, err := repo.WithTx(tx).CreateTicketType(wr.R.Context(), ticketTypeParams)
		if err != nil {
			logger.StdoutLogger.Error("Failed to create ticket type", "ticketTypeIndex", i, "err", err.Error())
			tx.Rollback(wr.R.Context()) // Rollback the transaction
			wr.Status(http.StatusInternalServerError).Json(
				utils.M{
					"message": fmt.Sprintf("Failed to create ticket type '%s'", ticketTypeDTO.Name),
					"err":     err.Error(),
				},
			)
			return
		}
		createdTicketTypes = append(createdTicketTypes, ticketType)
	}

	// Create coupons for the event (if any)
	var createdCoupons []repository.Coupon
	for i, couponDTO := range eventDTO.Coupons {
		// Parse coupon dates manually to ensure proper pgtype.Timestamptz format
		validFromTime, err := time.Parse(time.RFC3339, couponDTO.ValidFrom)
		if err != nil {
			logger.StdoutLogger.Error("Failed to parse coupon valid_from", "couponIndex", i, "err", err.Error())
			tx.Rollback(wr.R.Context())
			wr.Status(http.StatusBadRequest).Json(
				utils.M{
					"message": fmt.Sprintf("Invalid valid_from date in coupon '%s'", couponDTO.Code),
					"err":     err.Error(),
				},
			)
			return
		}

		validUntilTime, err := time.Parse(time.RFC3339, couponDTO.ValidUntil)
		if err != nil {
			logger.StdoutLogger.Error("Failed to parse coupon valid_until", "couponIndex", i, "err", err.Error())
			tx.Rollback(wr.R.Context())
			wr.Status(http.StatusBadRequest).Json(
				utils.M{
					"message": fmt.Sprintf("Invalid valid_until date in coupon '%s'", couponDTO.Code),
					"err":     err.Error(),
				},
			)
			return
		}

		// Create properly initialized Timestamptz values with Valid=true
		validFrom := pgtype.Timestamptz{Time: validFromTime, InfinityModifier: pgtype.Finite, Valid: true}
		validUntil := pgtype.Timestamptz{Time: validUntilTime, InfinityModifier: pgtype.Finite, Valid: true}

		couponParams := repository.CreateCouponParams{
			EventID:            event.ID,
			Code:               couponDTO.Code,
			DiscountPercentage: couponDTO.DiscountPercentage,
			ValidFrom:          validFrom,
			ValidUntil:         validUntil,
			UsageLimit:         int32(couponDTO.UsageLimit),
		}

		coupon, err := repo.WithTx(tx).CreateCoupon(wr.R.Context(), couponParams)
		if err != nil {
			logger.StdoutLogger.Error("Failed to create coupon", "couponIndex", i, "err", err.Error())
			tx.Rollback(wr.R.Context()) // Rollback the transaction
			wr.Status(http.StatusInternalServerError).Json(
				utils.M{
					"message": fmt.Sprintf("Failed to create coupon '%s'", couponDTO.Code),
					"err":     err.Error(),
				},
			)
			return
		}
		createdCoupons = append(createdCoupons, coupon)
	}

	// Commit the transaction to finalize the event creation
	// This makes all changes permanent in the database
	err = tx.Commit(wr.R.Context())
	if err != nil {
		logger.StdoutLogger.Error("Failed to commit transaction", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error, please try again later.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Return a successful response with HTTP 201 Created status
	// Include the created event data along with ticket types and coupons
	wr.Status(http.StatusCreated).Json(
		utils.M{
			"message": "Event created successfully with ticket types and coupons!",
			"data": utils.M{
				"event":        event,
				"ticket_types": createdTicketTypes,
				"coupons":      createdCoupons,
				"stats": utils.M{
					"ticket_types_created": len(createdTicketTypes),
					"coupons_created":      len(createdCoupons),
				},
			},
		},
	)
}

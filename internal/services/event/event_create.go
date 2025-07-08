// Package event provides services for event management in the Everato platform.
// It handles the creation, updating, deletion and querying of events.
package event

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
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

	// Generate a URL-friendly slug from the event title
	// This will be used in event URLs and must be unique across all events
	slug, err := utils.GenerateSlug(eventDTO.Title)
	if slug == "" && err != nil {
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

	// Handle case where the slug search failed for a reason other than "not found"
	// This indicates a database error rather than a slug uniqueness issue
	if err != pgx.ErrNoRows {
		logger.StdoutLogger.Error("Failed to create event", "err", err.Error())
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error, please try again later.",
				"err":     err.Error(),
			},
		)
		return
	}

	// Create the event in the database using the transaction
	// This converts our DTO to the format required by the repository
	event, err := repo.WithTx(tx).CreateEvent(wr.R.Context(), eventDTO.ToCreateEventParams())
	// Handle any other errors that might have occurred during event creation
	// This is a catch-all for unexpected issues
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

	if err != nil {
		logger.StdoutLogger.Error("Failed to create the event", "err", err.Error())
		tx.Rollback(wr.R.Context()) // Rollback the transaction
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Internal server error, please try again later.",
				"err":     err.Error(),
			},
		)
		return
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
	// Include the created event data in the response body
	wr.Status(http.StatusCreated).Json(
		utils.M{
			"message": "Event created successfully!",
			"data":    event,
		},
	)
}

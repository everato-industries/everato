package event

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
)

// This function will handle the creation of an event.
func CreateEvent(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Parse the request body to the CreateEventDTO
	eventDTO := &CreateEventDTO{}
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

	// Validate whether the data is valid or not
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

	// Start a transaction for ATOMICITY
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

	// Generate a slug from the title then search if that is unique or not
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

	// Create the event in the database
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

	if _, err := repo.WithTx(tx).CreateEvent(wr.R.Context(), eventDTO.ToCreateEventParams()); err != nil {
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
}

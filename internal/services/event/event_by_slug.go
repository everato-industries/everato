package event

import (
	"github.com/jackc/pgx/v5"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
)

// GetEventBySlug retrieves an event by its slug from the database and writes the response.
//
// Parameters:
//   - wr: An HttpWriter to write the HTTP response
//   - repo: The database repository for event operations
//   - conn: The database connection (not used in this function but included for consistency)
//
// Behavior:
//   - Extracts the 'slug' parameter from the URL
//   - Queries the database for an event matching the slug
//   - If found, responds with the event data and HTTP 200 status
//   - If not found, responds with HTTP 404 status and an error message
//   - If any other error occurs, responds with HTTP 500 status and an error message
func GetEventBySlug(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	slug := utils.GetParam(wr.R, "slug")

	if slug == "" {
		wr.Status(400).Json(
			utils.M{
				"message": "Slug parameter is required",
			},
		)
		return
	}

	event, err := repo.GetEventBySlug(wr.R.Context(), slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			wr.Status(404).Json(
				utils.M{
					"message": "Event not found",
				},
			)
			return
		}
		wr.Status(500).Json(
			utils.M{
				"message": "Error while fetching event",
				"error":   err.Error(),
			},
		)
		return
	}

	wr.Status(200).Json(
		utils.M{
			"message": "Event fetched successfully",
			"event":   event,
		},
	)
}

package event

import (
	"net/http"
	"strconv"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/jackc/pgx/v5"
)

func GetAllEvents(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	query := wr.R.URL.Query()

	limit, offset := "", ""
	if query.Has("limit") {
		limit = query.Get("limit")
	} else {
		limit = utils.GetEnv("DEFAULT_LIMIT", "10")
	}

	if query.Has("offset") {
		offset = query.Get("offset")
	} else {
		offset = utils.GetEnv("DEFAULT_OFFSET", "0")
	}

	// Convert limit and offset to integers
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Passed limit value is not good :(",
				"error":   err.Error(),
			},
		)
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Passed limit value is not good :(",
				"error":   err.Error(),
			},
		)
		return
	}

	// Get total count of events for pagination
	totalCount, err := repo.CountTotalEvents(wr.R.Context())
	if err != nil {
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Error while counting events",
				"error":   err.Error(),
			},
		)
		return
	}

	// Get paginated events
	events, err := repo.ListEvents(wr.R.Context(), repository.ListEventsParams{
		Limit:  int32(limitInt),
		Offset: int32(offsetInt),
	})
	if err != nil {
		wr.Status(http.StatusInternalServerError).Json(
			utils.M{
				"message": "Error while fetching events",
				"error":   err.Error(),
			},
		)
		return
	}

	// Calculate pagination metadata
	totalPages := (int(totalCount) + limitInt - 1) / limitInt
	currentPage := (offsetInt / limitInt) + 1

	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Events fetched successfully",
			"data":    events,
			"pagination": utils.M{
				"total_count":  totalCount,
				"total_pages":  totalPages,
				"current_page": currentPage,
				"limit":        limitInt,
				"offset":       offsetInt,
				"has_next":     currentPage < totalPages,
				"has_previous": currentPage > 1,
			},
		},
	)
}

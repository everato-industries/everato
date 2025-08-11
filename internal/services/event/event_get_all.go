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

	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Events fetched successfully",
			"data":    events,
		},
	)
}

package event

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
)

func GetRecentEvents(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn) {
	// Get limit parameter from query (default: 10, max: 50)
	limitStr := wr.R.URL.Query().Get("limit")
	limit := int32(10) // default limit

	if limitStr != "" {
		parsedLimit, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || parsedLimit < 1 || parsedLimit > 50 {
			wr.Status(http.StatusBadRequest).Json(utils.M{
				"message": "Invalid limit parameter. Must be between 1 and 50",
				"error":   "Limit must be a number between 1 and 50",
			})
			return
		}
		limit = int32(parsedLimit)
	}

	// Get recent events from database
	recentEvents, err := repo.GetRecentEvents(wr.R.Context(), limit)
	if err != nil {
		wr.Status(http.StatusInternalServerError).Json(utils.M{
			"message": "Failed to fetch recent events",
			"error":   err.Error(),
		})
		return
	}

	// Format response data
	var events []map[string]interface{}
	for _, event := range recentEvents {
		eventData := map[string]interface{}{
			"id":              event.ID,
			"title":           event.Title,
			"description":     event.Description,
			"banner":          event.Banner,
			"icon":            event.Icon,
			"start_time":      event.StartTime,
			"end_time":        event.EndTime,
			"location":        event.Location,
			"status":          string(event.Status),
			"slug":            event.Slug,
			"total_seats":     event.TotalSeats,
			"available_seats": event.AvailableSeats,
			"created_at":      event.CreatedAt,
			"updated_at":      event.UpdatedAt,
		}

		// Add admin ID if available
		if event.AdminID.Valid {
			eventData["admin_id"] = event.AdminID.String()
		}

		events = append(events, eventData)
	}

	wr.Status(http.StatusOK).Json(utils.M{
		"message": "Recent events fetched successfully",
		"data": map[string]interface{}{
			"events": events,
			"count":  len(events),
			"limit":  limit,
		},
	})
}

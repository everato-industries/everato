package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

type DashboardHandler struct {
	Repo     *repository.Queries // Database repository for dashboard operations
	Conn     *pgx.Conn           // Database connection
	BasePath string              // Base URL path for dashboard endpoint
	Config   *config.Config      // Application configuration
}

var _ handlers.Handler = (*DashboardHandler)(nil)

func NewDashboardHandler(cfg *config.Config) *DashboardHandler {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Establish connection to the PostgreSQL database
	conn, err := pgx.Connect(
		context.Background(),
		utils.GetEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato?ssl_mode=disable"),
	)
	if err != nil {
		logger.StdoutLogger.Error("Error connecting to the postgres db", "err", err.Error())
		return &DashboardHandler{
			Repo: nil,
		}
	}

	// Initialize repository with database connection
	repo := repository.New(conn)
	return &DashboardHandler{
		Repo:     repo,
		Conn:     conn,
		BasePath: "/dashboard",
		Config:   cfg,
	}
}

func (h *DashboardHandler) RegisterRoutes(router *mux.Router) {
	r := router.PathPrefix(h.BasePath).Subrouter()

	// GET /dashboard/stats
	r.HandleFunc("/stats", h.Stats).Methods(http.MethodGet)
}

// Stats route displays the stats for the dashboard
//
// HTTP Method: GET
// Route: /api/v1/dashboard/stats
//
// Response:
//   - 200 OK with JSON containing dashboard statistics
func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)
	logger := pkg.NewLogger()
	defer logger.Close()

	// Check if repository is available
	if h.Repo == nil {
		logger.Error("Database repository not available")
		wr.Status(http.StatusInternalServerError).Json(utils.M{
			"error": "Database connection not available",
		})
		return
	}

	ctx := context.Background()

	// Get dashboard stats from database
	eventStats, err := h.Repo.GetDashboardStats(ctx)
	if err != nil {
		logger.Error("Failed to get event stats", "error", err)
		wr.Status(http.StatusInternalServerError).Json(utils.M{
			"error": "Failed to fetch event statistics",
		})
		return
	}

	// Get total users count
	userCount, err := h.Repo.CountTotalUsers(ctx)
	if err != nil {
		logger.Error("Failed to get user count", "error", err)
		wr.Status(http.StatusInternalServerError).Json(utils.M{
			"error": "Failed to fetch user statistics",
		})
		return
	}

	// Prepare response
	response := utils.M{
		"total_events":       eventStats.TotalEvents,
		"total_users":        userCount,
		"upcoming_events":    eventStats.UpcomingEvents,
		"active_events":      eventStats.ActiveEvents,
		"completed_events":   eventStats.CompletedEvents,
		"created_events":     eventStats.CreatedEvents,
		"cancelled_events":   eventStats.CancelledEvents,
		"total_tickets_sold": 0,        // Keep as dummy for now
		"total_revenue":      12345.67, // Keep as dummy for now
	}

	wr.Status(http.StatusOK).Json(response)
}

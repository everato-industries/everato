// Package api provides HTTP handlers for the Everato platform REST API.
package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/middlewares"
	"github.com/dtg-lucifer/everato/internal/services/checkin"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// CheckinHandler manages check-in and attendance HTTP endpoints in the API.
// It allows event administrators to scan QR codes, mark attendance, and view stats.
//
// Route prefix:
//   - `/api/v1/checkin`
//
// Routes:
//   - POST /checkin/scan - Scan QR code and check in attendee (POST)
//   - POST /checkin/validate - Validate QR code and check in attendee (alias) (POST)
//   - GET /checkin/events/{eventID}/stats - Get attendance stats for an event (GET)
//   - GET /checkin/events/{eventID}/attendees - Get attendee list for an event (GET)
type CheckinHandler struct {
	Repo     *repository.Queries
	Conn     *pgx.Conn
	BasePath string
	Cfg      *config.Config
}

// Asserting the implementation of the handler interface
var _ handlers.Handler = (*CheckinHandler)(nil)

// NewCheckinHandler creates and initializes a new CheckinHandler instance.
func NewCheckinHandler(cfg *config.Config) *CheckinHandler {
	logger := pkg.NewLogger()
	defer logger.Close()

	conn, err := pgx.Connect(
		context.Background(),
		utils.GetEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato?ssl_mode=disable"),
	)
	if err != nil {
		logger.StdoutLogger.Error("Error connecting to the postgres db", "err", err.Error())
		return &CheckinHandler{
			Repo:     nil,
			BasePath: "/checkin",
			Cfg:      cfg,
		}
	}

	repo := repository.New(conn)
	return &CheckinHandler{
		Repo:     repo,
		Conn:     conn,
		BasePath: "/checkin",
		Cfg:      cfg,
	}
}

// RegisterRoutes registers all check-in routes with the provided router.
func (h *CheckinHandler) RegisterRoutes(router *mux.Router) {
	checkinGroup := router.PathPrefix(h.BasePath).Subrouter()

	guard := middlewares.NewAdminMiddleware(h.Repo, h.Conn, false)
	protected := checkinGroup.NewRoute().Subrouter()
	protected.Use(guard.Guard)

	protected.HandleFunc("/scan", h.ScanTicket).Methods(http.MethodPost)
	protected.HandleFunc("/validate", h.ScanTicket).Methods(http.MethodPost)
	protected.HandleFunc("/events/{eventID}/stats", h.GetStats).Methods(http.MethodGet)
	protected.HandleFunc("/events/{eventID}/attendees", h.GetAttendees).Methods(http.MethodGet)
}

// ScanTicket processes a scanned QR code to validate and record check-in.
func (h *CheckinHandler) ScanTicket(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)
	if h.Repo == nil || h.Conn == nil {
		wr.Status(http.StatusBadGateway).Json(utils.M{"error": "No database connection"})
		return
	}
	secret := utils.GetEnv("JWT_SECRET", "SUPER_SECRET_KEY")
	checkin.ValidateAndCheckin(wr, h.Repo, h.Conn, secret)
}

// GetStats returns check-in statistics for a specific event.
func (h *CheckinHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(utils.M{"error": "No database connection"})
		return
	}
	checkin.GetEventAttendanceStats(wr, h.Repo)
}

// GetAttendees returns the full list of attendees for an event.
func (h *CheckinHandler) GetAttendees(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(utils.M{"error": "No database connection"})
		return
	}
	checkin.GetEventAttendees(wr, h.Repo)
}

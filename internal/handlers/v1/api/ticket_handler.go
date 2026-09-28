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
	"github.com/dtg-lucifer/everato/internal/services/ticket"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// TicketHandler manages ticket-related HTTP endpoints in the API.
// It handles retrieval of tickets and QR code data for bookings.
//
// Route prefix:
//   - `/api/v1/tickets`
//
// Routes:
//   - GET /tickets/booking/{bookingID} - Get all tickets for a booking
//   - GET /tickets/{ticketID} - Get ticket details by ID
type TicketHandler struct {
	Repo     *repository.Queries
	Conn     *pgx.Conn
	BasePath string
	Cfg      *config.Config
}

// Asserting the implementation of the handler interface
var _ handlers.Handler = (*TicketHandler)(nil)

// NewTicketHandler creates and initializes a new TicketHandler instance.
func NewTicketHandler(cfg *config.Config) *TicketHandler {
	logger := pkg.NewLogger()
	defer logger.Close()

	conn, err := pgx.Connect(
		context.Background(),
		utils.GetEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato?ssl_mode=disable"),
	)
	if err != nil {
		logger.StdoutLogger.Error("Error connecting to the postgres db", "err", err.Error())
		return &TicketHandler{
			Repo:     nil,
			BasePath: "/tickets",
			Cfg:      cfg,
		}
	}

	repo := repository.New(conn)
	return &TicketHandler{
		Repo:     repo,
		Conn:     conn,
		BasePath: "/tickets",
		Cfg:      cfg,
	}
}

// RegisterRoutes registers all ticket-related routes with the provided router.
func (h *TicketHandler) RegisterRoutes(router *mux.Router) {
	tickets := router.PathPrefix(h.BasePath).Subrouter()

	guard := middlewares.NewAuthMiddleware(h.Repo, h.Conn, false)
	protected := tickets.NewRoute().Subrouter()
	protected.Use(guard.Guard)

	protected.HandleFunc("/booking/{bookingID}", h.GetTicketsByBooking).Methods(http.MethodGet)
	protected.HandleFunc("/{ticketID}", h.GetTicketByID).Methods(http.MethodGet)
}

// GetTicketsByBooking handles retrieval of all tickets associated with a booking.
func (h *TicketHandler) GetTicketsByBooking(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(utils.M{"error": "No database connection"})
		return
	}
	ticket.GetTicketsByBooking(wr, h.Repo, h.Conn)
}

// GetTicketByID handles retrieval of a specific ticket by its ID.
func (h *TicketHandler) GetTicketByID(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(utils.M{"error": "No database connection"})
		return
	}
	ticket.GetTicketByID(wr, h.Repo, h.Conn)
}

package v1

import (
	"context"
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/services/event"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// EventHandler is the handler for event-related operations.
type EventHandler struct {
	Repo     *repository.Queries
	Conn     *pgx.Conn
	BasePath string
}

// -----------------------------------------------------
var _ handlers.Handler = (*EventHandler)(nil) // Assert the interface implementation to catch errors
// -----------------------------------------------------

func NewEventHandler() *EventHandler {
	logger := pkg.NewLogger()
	defer logger.Close()

	conn, err := pgx.Connect(
		context.Background(),
		utils.GetEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato?ssl_mode=disable"),
	)
	if err != nil {
		logger.StdoutLogger.Error("Error connecting to the postgres db", "err", err.Error())
		return &EventHandler{
			Repo: nil,
		}
	}

	repo := repository.New(conn)
	return &EventHandler{
		Repo:     repo,
		Conn:     conn,
		BasePath: "/events",
	}
}

// RegisterRoutes registers the event-related routes with the given router.
func (h *EventHandler) RegisterRoutes(router *mux.Router) {
	// Create a subrouter
	events := router.PathPrefix(h.BasePath).Subrouter()

	// Routes
	events.HandleFunc("/create", h.CreateEvent).Methods(http.MethodPost) // Create a new event
	events.HandleFunc("/update", h.UpdateEvent).Methods(http.MethodPut)  // Update an event
	events.HandleFunc("/all", h.GetAllEvents).Methods(http.MethodPost)   // Get all events
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	// If no repo is set then there must be a error, to be sure ABORT!
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(
			utils.M{
				"message": "BAD_GATEWAY, No database connection, Oops!",
			},
		)
		return
	}

	// Call the service
	event.CreateEvent(wr, h.Repo, h.Conn)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {}

func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {}

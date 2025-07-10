// Package api provides handlers for the REST API endpoints of the Everato application.
// This package contains all the HTTP handlers for API routes under /api/v1/
package api

import (
	"context"
	"net/http"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/middlewares"
	"github.com/dtg-lucifer/everato/internal/services/event"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// EventHandler manages event-related HTTP endpoints in the API.
// It handles operations like event creation, updating, listing, and searching.
//
// Route prefix:
//   - `/api/vX/events/`
//
// Routes:
//   - /create - Create a new event (POST)
//   - /update - Update an existing event (PUT)
//   - /all    - List all events with pagination (POST)
type EventHandler struct {
	Repo     *repository.Queries // Database repository for event operations
	Conn     *pgx.Conn           // Database connection for transactions
	BasePath string              // Base URL path for event endpoints
}

// -----------------------------------------------------
var _ handlers.Handler = (*EventHandler)(nil) // Assert the interface implementation to catch errors
// -----------------------------------------------------


// NewEventHandler creates and initializes a new EventHandler instance.
// It establishes a database connection and initializes the repository for event operations.
//
// Returns:
//   - A fully initialized EventHandler, or partially initialized handler if DB connection fails
//     (in which case the Repo field will be nil)
func NewEventHandler() *EventHandler {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Establish connection to the PostgreSQL database
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

	// Initialize repository with database connection
	repo := repository.New(conn)
	return &EventHandler{
		Repo:     repo,
		Conn:     conn,
		BasePath: "/events",
		Cfg:      cfg,
	}
}

// RegisterRoutes registers all event-related routes with the provided router.
// It creates a subrouter with the base path, applies authentication middleware,
// and maps HTTP methods to handler functions.
//
// All event endpoints require authentication, enforced by the AuthGuard middleware.
// The middleware expects a valid JWT token in either cookies or Authorization header.
//
// Parameters:
//   - router: The main router to attach event routes to
func (h *EventHandler) RegisterRoutes(router *mux.Router) {
	// Create a subrouter for all event routes
	events := router.PathPrefix(h.BasePath).Subrouter()

	// Apply authentication middleware to all event routes
	// Uses JSON responses for auth failures instead of redirects
	guard := middlewares.NewAuthMiddleware(h.Repo, h.Conn, false)
	events.Use(guard.Guard) // Guard the whole route group

	// Register individual route handlers
	events.HandleFunc("/create", h.CreateEvent).Methods(http.MethodPost) // Create a new event
	events.HandleFunc("/update", h.UpdateEvent).Methods(http.MethodPut)  // Update an existing event
	events.HandleFunc("/all", h.GetAllEvents).Methods(http.MethodPost)   // Get all events with filtering
}

// CreateEvent handles requests to create a new event in the system.
// It validates the request and delegates the business logic to the event service.
//
// This handler expects a JSON request body containing event details including:
// title, description, start time, end time, location, and capacity information.
//
// HTTP Method: POST
// Route: /api/v1/events/create
//
// Request: JSON event data
// Response:
//   - 201 Created with event details on success
//   - 400 Bad Request if validation fails
//   - 401 Unauthorized if user is not authenticated
//   - 409 Conflict if event slug already exists
//   - 502 Bad Gateway if database connection fails
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	// Validate database repository connectivity
	// This ensures the handler can't proceed without a valid database connection
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(
			utils.M{
				"message": "BAD_GATEWAY, No database connection, Oops!",
			},
		)
		return
	}

	// Delegate event creation to the service layer
	// This separates HTTP handling from business logic
	event.CreateEvent(wr, h.Repo, h.Conn)
}

// UpdateEvent handles requests to update an existing event.
// It validates the request and delegates the update logic to the event service.
//
// HTTP Method: PUT
// Route: /api/v1/events/update
//
// Request: JSON with event ID and fields to update
// Response:
//   - 200 OK with updated event details on success
//   - 400 Bad Request if validation fails
//   - 401 Unauthorized if user is not authenticated
//   - 403 Forbidden if user doesn't have permission to update the event
//   - 404 Not Found if event doesn't exist
//   - 502 Bad Gateway if database connection fails
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement event update functionality
}

// GetAllEvents handles requests to list events with optional filtering and pagination.
// It supports filtering by date range, location, and search terms.
//
// HTTP Method: POST
// Route: /api/v1/events/all
//
// Request: JSON with optional filter criteria and pagination parameters
// Response:
//   - 200 OK with paginated list of events
//   - 400 Bad Request if filter parameters are invalid
//   - 401 Unauthorized if user is not authenticated
//   - 502 Bad Gateway if database connection fails
func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement event listing functionality
}

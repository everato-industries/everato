// Package api provides handlers for the REST API endpoints of the Everato application.
// This package contains all the HTTP handlers for API routes under /api/v1/
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
	"github.com/dtg-lucifer/everato/internal/services/booking"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
)

// BookingHandler manages booking-related HTTP endpoints in the API.
// It handles operations like booking creation, retrieval, and cancellation.
//
// Route prefix:
//   - `/api/v1/bookings/`
//
// Routes:
//   - POST /bookings/create - Create a new booking (POST)
//   - GET /bookings/user/{userId} - Get user bookings (GET)
//   - GET /bookings/{bookingId} - Get booking details (GET)
//   - PUT /bookings/{bookingId}/status - Update booking status (PUT)
//   - DELETE /bookings/{bookingId} - Cancel booking (DELETE)
//   - GET /events/{eventId}/availability - Check ticket availability (GET)
type BookingHandler struct {
	Repo     *repository.Queries // Database repository for booking operations
	Conn     *pgx.Conn           // Database connection for transactions
	BasePath string              // Base URL path for booking endpoints
	Cfg      *config.Config      // Application configuration
}

// Asserting the implementation of the handler interface
var _ handlers.Handler = (*BookingHandler)(nil)

// NewBookingHandler creates and initializes a new BookingHandler instance.
// It establishes a database connection and initializes the repository.
//
// Returns:
//   - A fully initialized BookingHandler, or partially initialized handler if DB connection fails
func NewBookingHandler(cfg *config.Config) *BookingHandler {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Establish connection to the PostgreSQL database
	conn, err := pgx.Connect(
		context.Background(),
		utils.GetEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato?ssl_mode=disable"),
	)
	if err != nil {
		logger.StdoutLogger.Error("Error connecting to the postgres db", "err", err.Error())
		return &BookingHandler{
			Repo: nil,
		}
	}

	// Initialize repository with database connection
	repo := repository.New(conn)
	return &BookingHandler{
		Repo:     repo,
		Conn:     conn,
		BasePath: "/bookings",
		Cfg:      cfg,
	}
}

// RegisterRoutes registers all booking-related routes with the provided router.
// It creates a subrouter with the base path and maps HTTP methods to handler functions.
//
// Public endpoints (no authentication required):
//   - GET /events/{eventId}/availability - Check ticket availability
//
// Protected endpoints (require user authentication):
//   - POST /bookings/create - Create a new booking
//   - GET /bookings/user/{userId} - Get user's bookings
//   - GET /bookings/{bookingId} - Get booking details
//   - DELETE /bookings/{bookingId} - Cancel booking
//
// Parameters:
//   - router: The main router to attach booking routes to
func (h *BookingHandler) RegisterRoutes(router *mux.Router) {
	// Create a subrouter for booking routes
	bookings := router.PathPrefix(h.BasePath).Subrouter()

	// Public endpoints (no authentication required)
	// These are mounted on the events path for availability checking
	router.HandleFunc("/events/{eventId}/availability", h.CheckAvailability).Methods(http.MethodGet)

	// Create the AuthGuard for protected routes
	guard := middlewares.NewAuthMiddleware(h.Repo, h.Conn, false)
	protected := bookings.NewRoute().Subrouter()
	protected.Use(guard.Guard)

	// Register protected route handlers (authenticated users only)
	protected.HandleFunc("/create", h.CreateBooking).Methods(http.MethodPost)
	protected.HandleFunc("/user/{userId}", h.GetUserBookings).Methods(http.MethodGet)
	protected.HandleFunc("/{bookingId}", h.GetBookingDetails).Methods(http.MethodGet)
	protected.HandleFunc("/{bookingId}", h.CancelBooking).Methods(http.MethodDelete)
}

// CreateBooking handles requests to create a new booking.
// It validates the request and delegates the business logic to the booking service.
//
// HTTP Method: POST
// Route: /api/v1/bookings/create
//
// Request: JSON with event ID, ticket selections, and optional coupon code
// Response:
//   - 201 Created with booking details on success
//   - 400 Bad Request if validation fails
//   - 401 Unauthorized if user is not authenticated
//   - 409 Conflict if tickets no longer available
//   - 502 Bad Gateway if database connection fails
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	// Validate database repository connectivity
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(
			utils.M{
				"message": "BAD_GATEWAY, No database connection, Oops!",
			},
		)
		return
	}

	// Delegate booking creation to the service layer
	booking.CreateBooking(wr, h.Repo, h.Conn, h.Cfg)
}

// GetUserBookings handles requests to retrieve all bookings for a specific user.
// It validates the request and delegates the retrieval logic to the booking service.
//
// HTTP Method: GET
// Route: /api/v1/bookings/user/{userId}
//
// Path Parameters:
//   - userId: The user ID to retrieve bookings for
//
// Response:
//   - 200 OK with list of bookings on success
//   - 400 Bad Request if userId is missing or invalid
//   - 401 Unauthorized if user is not authenticated
//   - 403 Forbidden if user tries to access another user's bookings
//   - 502 Bad Gateway if database connection fails
func (h *BookingHandler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	// Validate database repository connectivity
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(
			utils.M{
				"message": "BAD_GATEWAY, No database connection, Oops!",
			},
		)
		return
	}

	// Delegate to the service layer
	booking.GetUserBookings(wr, h.Repo, h.Conn)
}

// GetBookingDetails handles requests to retrieve details for a specific booking.
// It validates the request and delegates the retrieval logic to the booking service.
//
// HTTP Method: GET
// Route: /api/v1/bookings/{bookingId}
//
// Path Parameters:
//   - bookingId: The booking ID to retrieve
//
// Response:
//   - 200 OK with booking details on success
//   - 400 Bad Request if bookingId is missing or invalid
//   - 401 Unauthorized if user is not authenticated
//   - 403 Forbidden if user tries to access another user's booking
//   - 404 Not Found if booking doesn't exist
//   - 502 Bad Gateway if database connection fails
func (h *BookingHandler) GetBookingDetails(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	// Validate database repository connectivity
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(
			utils.M{
				"message": "BAD_GATEWAY, No database connection, Oops!",
			},
		)
		return
	}

	// Delegate to the service layer
	booking.GetBookingDetails(wr, h.Repo, h.Conn)
}

// CancelBooking handles requests to cancel a booking.
// It validates the request and delegates the cancellation logic to the booking service.
//
// HTTP Method: DELETE
// Route: /api/v1/bookings/{bookingId}
//
// Path Parameters:
//   - bookingId: The booking ID to cancel
//
// Response:
//   - 200 OK with cancellation confirmation on success
//   - 400 Bad Request if bookingId is missing or invalid
//   - 401 Unauthorized if user is not authenticated
//   - 403 Forbidden if user tries to cancel another user's booking
//   - 404 Not Found if booking doesn't exist
//   - 409 Conflict if booking cannot be cancelled (e.g., already confirmed)
//   - 502 Bad Gateway if database connection fails
func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	// Validate database repository connectivity
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(
			utils.M{
				"message": "BAD_GATEWAY, No database connection, Oops!",
			},
		)
		return
	}

	// Delegate to the service layer
	booking.CancelBooking(wr, h.Repo, h.Conn)
}

// CheckAvailability handles requests to check ticket availability for an event.
// It validates the request and delegates the availability check to the booking service.
//
// HTTP Method: GET
// Route: /api/v1/events/{eventId}/availability
//
// Path Parameters:
//   - eventId: The event ID to check availability for
//
// Response:
//   - 200 OK with availability information on success
//   - 400 Bad Request if eventId is missing or invalid
//   - 404 Not Found if event doesn't exist
//   - 502 Bad Gateway if database connection fails
func (h *BookingHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	// Validate database repository connectivity
	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(
			utils.M{
				"message": "BAD_GATEWAY, No database connection, Oops!",
			},
		)
		return
	}

	// Delegate to the service layer
	booking.CheckAvailability(wr, h.Repo, h.Conn)
}

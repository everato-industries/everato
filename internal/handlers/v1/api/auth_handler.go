// Package api provides handlers for the REST API endpoints of the Everato application.
// This package contains all the HTTP handlers for API routes under /api/v1/
package api

import (
	"context"
	"net/http"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/services/user"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// AuthHandler handles authentication-related routes and logic.
//
// It manages user registration, login, email verification, token refreshing,
// and password management operations through RESTful endpoints.
//
// Route prefix:
//   - `/api/vX/auth/`
//
// Routes:
//   - /register      - Create a new user account (POST)
//   - /login         - Authenticate and receive JWT token (POST)
//   - /verify-email  - Verify user email address (GET)
//   - /refresh       - Refresh an expired JWT token (POST)
//   - /reset-password - Request password reset (POST)
//   - /change-password - Change user password with token (POST)
type AuthHandler struct {
	Repo     *repository.Queries // Database repository for user operations
	Conn     *pgx.Conn           // Database connection
	BasePath string              // Base URL path for auth endpoints
	Cfg      *config.Config      // Application configuration
}

// Asserting the implementation of the handler interface
var _ handlers.Handler = (*AuthHandler)(nil)

// NewAuthHandler creates and initializes a new AuthHandler instance.
// It establishes a database connection and initializes the repository.
//
// Returns:
//   - A fully initialized AuthHandler, or a partially initialized handler if DB connection fails
//     (in which case the Repo field will be nil)
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Establish connection to the PostgreSQL database using connection string from environment
	conn, err := pgx.Connect(
		context.Background(),
		utils.GetEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato?ssl_mode=disable"),
	)
	if err != nil {
		logger.StdoutLogger.Error("Error connecting to the postgres db", "err", err.Error())
		return &AuthHandler{
			Repo: nil,
		}
	}

	// Initialize the repository with the database connection
	repo := repository.New(conn)

	return &AuthHandler{
		Repo:     repo,
		Conn:     conn,
		BasePath: "/auth",
		Cfg:      cfg,
	}
}

// RegisterRoutes registers all authentication routes with the provided router.
// It creates a subrouter with the base path and maps HTTP methods to handler functions.
//
// Parameters:
//   - router: The main router to attach auth routes to
func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	auth_router := router.PathPrefix(h.BasePath).Subrouter()

	auth_router.HandleFunc("/register", h.Register).Methods(http.MethodPost)              // Register a new user
	auth_router.HandleFunc("/login", h.Login).Methods(http.MethodPost)                    // Login an existing user (returning the jwt token)
	auth_router.HandleFunc("/verify-email", h.VerifyEmail).Methods(http.MethodGet)        // Verify the user's email address
	auth_router.HandleFunc("/refresh", h.Refresh).Methods(http.MethodPost)                // Refresh the JWT token @TODO
	auth_router.HandleFunc("/reset-password", h.ResetPassword).Methods(http.MethodPost)   // Reset the user's password @TODO
	auth_router.HandleFunc("/change-password", h.ChangePassword).Methods(http.MethodPost) // Change the user's password @TODO
}

// Register handles new user registration.
// It creates a new user account and sends a verification email.
//
// HTTP Method: POST
// Route: /api/v1/auth/register
//
// Request body: JSON with user registration details
// Response:
//   - 201 Created on success with user details
//   - 400 Bad Request if request data is invalid
//   - 409 Conflict if email already exists
//   - 502 Bad Gateway if database connection fails
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
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

	// Create a new user in the database
	user.CreateUser(wr, h.Repo, h.Conn, h.Cfg)
}

// Login authenticates a user and returns a JWT token.
// On successful authentication, it sets a JWT cookie and returns user info with the token.
//
// HTTP Method: POST
// Route: /api/v1/auth/login
//
// Request body: JSON with email and password
// Response:
//   - 200 OK on success with JWT token and user details
//   - 401 Unauthorized if credentials are invalid
//   - 502 Bad Gateway if database connection fails
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	user.LoginUser(wr, h.Repo, h.Conn)
}

// VerifyEmail handles email verification requests.
// It processes the verification token from the URL and updates the user's verification status.
//
// HTTP Method: GET
// Route: /api/v1/auth/verify-email?uid={user_id}
//
// URL Parameters:
//   - uid: User ID to verify
//
// Response:
//   - HTML page showing verification success or failure
//   - 502 Bad Gateway if database connection fails
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
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

	// Upon clicking on the URL of the verify email this will set the verified = True of the user
	user.VerifyUser(wr, h.Repo, h.Conn)
}

// Refresh handles JWT token refresh requests.
// It verifies the existing token and issues a new one with extended expiration.
//
// HTTP Method: POST
// Route: /api/v1/auth/refresh
//
// Request: Expects refresh token in authorization header or cookie
// Response:
//   - 200 OK with new token pair on success
//   - 401 Unauthorized if refresh token is invalid
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement token refresh logic here
	// 1. Extract refresh token from request
	// 2. Validate refresh token
	// 3. Generate new access token with updated expiration
	// 4. Return new token pair
}

// ResetPassword initiates the password reset process.
// It sends a reset link to the user's email address.
//
// HTTP Method: POST
// Route: /api/v1/auth/reset-password
//
// Request: JSON with user's email
// Response:
//   - 200 OK if reset email was sent successfully
//   - 404 Not Found if email doesn't match any user
//   - 500 Internal Server Error if email sending fails
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement password reset logic here
	// 1. Validate email address
	// 2. Generate reset token
	// 3. Store token in database with expiration
	// 4. Send reset email with token link
}

// ChangePassword handles password change requests.
// It validates the reset token and updates the user's password.
//
// HTTP Method: POST
// Route: /api/v1/auth/change-password
//
// Request: JSON with reset token and new password
// Response:
//   - 200 OK if password was changed successfully
//   - 400 Bad Request if password doesn't meet requirements
//   - 401 Unauthorized if reset token is invalid or expired
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement password change logic here
	// 1. Validate reset token
	// 2. Validate new password meets requirements
	// 3. Hash the new password
	// 4. Update password in database
	// 5. Invalidate used token
	// 6. Notify user of successful password change
}

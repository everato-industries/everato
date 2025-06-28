package v1

import (
	"context"
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/services/user"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// AuthHandler handles authentication-related routes and logic
//
// Route prefix:
//   - `/api/vX/auth/`
//
// Routes:
//   - /register
//   - /login
//   - /verify-email
//   - /refresh
//   - /reset-password
//   - /change-password
type AuthHandler struct {
	Repo     *repository.Queries
	Conn     *pgx.Conn
	BasePath string
}

// Asserting the implementation of the handler interface
var _ handlers.Handler = (*AuthHandler)(nil)

// Returns a pointer to a newly created the AuthHandler instance
func NewAuthHandler() *AuthHandler {
	logger := pkg.NewLogger()
	defer logger.Close()

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

	repo := repository.New(conn)

	return &AuthHandler{
		Repo:     repo,
		Conn:     conn,
		BasePath: "/auth",
	}
}

func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	auth_router := router.PathPrefix("/auth").Subrouter()

	auth_router.HandleFunc("/register", h.Register).Methods(http.MethodPost)              // Register a new user
	auth_router.HandleFunc("/login", h.Login).Methods(http.MethodPost)                    // Login an existing user (returning the jwt token)
	auth_router.HandleFunc("/verify-email", h.VerifyEmail).Methods(http.MethodGet)        // Verify the user's email address
	auth_router.HandleFunc("/refresh", h.Refresh).Methods(http.MethodPost)                // Refresh the JWT token
	auth_router.HandleFunc("/reset-password", h.ResetPassword).Methods(http.MethodPost)   // Reset the user's password
	auth_router.HandleFunc("/change-password", h.ChangePassword).Methods(http.MethodPost) // Change the user's password
}

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
	user.CreateUser(wr, h.Repo, h.Conn)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Implement login logic here
}

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

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Implement token refresh logic here
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	// Implement password reset logic here
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Implement password change logic here
}

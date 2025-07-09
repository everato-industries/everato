package api

import (
	"context"
	"net/http"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/services/admin"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// AdminHandler handles administrative operations for the Everato platform.
//
// It manages the creation of other admin accounts which will help
// other admins to login and manage the events and stuffs in the dashboard
//
// Route group:
//   - /admin/*
//
// Routes:
//   - POST /admin/login - to login into a super-admin account
//   - POST /admin/create - to create sub-admin accounts
//   - POST /admin/:id - to update the details and roles of a sub-admin
//   - POST /admin/send-verification/:id - to send the verification email to the sub-admin's email
//   - DELETE /admin/:id - to delete an admin account
//   - GET /admin/all - to get all the admin accounts
//   - GET /admin/:id - to get a single admin details (with id)
//   - GET /admin/u/:username - to get a single admin details (with user name)
//   - GET /admin/:query - search any admin that mathes the query with email or username
type AdminHandler struct {
	Repo     *repository.Queries // Database repository for admin operations
	Conn     *pgx.Conn           // Database connection
	Cfg      *config.Config      // Application configuration
	BasePath string              // Base URL path for admin endpoints
}

// Manually assert the implementation
var _ handlers.Handler = (*AdminHandler)(nil)

// Initialize a new instance of admin handler
func NewAdminHandler(cfg *config.Config) *AdminHandler {
	logger := pkg.NewLogger()
	defer logger.Close()

	// Establish connection to the PostgreSQL database using connection string from environment
	conn, err := pgx.Connect(
		context.Background(),
		utils.GetEnv("DB_URL", "postgres://piush:root_access@localhost:5432/everato?ssl_mode=disable"),
	)
	if err != nil {
		logger.StdoutLogger.Error("Error connecting to the postgres db", "err", err.Error())
		return &AdminHandler{
			Repo:     nil,
			Conn:     nil,
			BasePath: "",
			Cfg:      cfg,
		}
	}

	repo := repository.New(conn)

	return &AdminHandler{
		Repo:     repo,
		Conn:     conn,
		Cfg:      cfg,
		BasePath: "/admin",
	}
}

// RegisterRoutes registers the admin handler routes with the provided router.
//
// It sets up the necessary endpoints for managing admin accounts,
// including login, creation, updates, deletion, and retrieval of admin accounts.
func (h *AdminHandler) RegisterRoutes(r *mux.Router) {
	router := r.PathPrefix(h.BasePath).Subrouter()

	router.HandleFunc("/login", h.Login).Methods(http.MethodPost)                                  // Login to an admin account
	router.HandleFunc("/create", h.CreateAdmin).Methods(http.MethodPost)                           // Create a new admin account
	router.HandleFunc("/{id}", h.UpdateAdmin).Methods(http.MethodPut)                              // Update an admin by ID
	router.HandleFunc("/send-verification/{id}", h.SendVerificationEmail).Methods(http.MethodPost) // Send verification email to an admin by ID
	router.HandleFunc("/{id}", h.DeleteAdmin).Methods(http.MethodDelete)                           // Delete an admin by ID
	router.HandleFunc("/all", h.GetAllAdmins).Methods(http.MethodGet)                              // Get all admins
	router.HandleFunc("/{id}", h.GetAdminByID).Methods(http.MethodGet)                             // Get admin by ID
	router.HandleFunc("/u/{username}", h.GetAdminByID).Methods(http.MethodGet)                     // Get admin by username
	router.HandleFunc("/{query}", h.SearchAdminByQeury).Methods(http.MethodGet)                    // Search admin by query
}

// Close cleans up the database connection when the handler is no longer needed.
func (h *AdminHandler) Close() {
	if h.Conn != nil {
		err := h.Conn.Close(context.Background())
		if err != nil {
			pkg.NewLogger().StdoutLogger.Error("Error closing the database connection", "err", err.Error())
		}
	}
}

// This handles the login service
func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	// If either of the repo or the conn is nil that means that there is some error
	if h.Repo == nil || h.Conn == nil {
		pkg.NewLogger().StdoutLogger.Error("Database repository or connection is not initialized")
		return
	}

	admin.Login(wr, h.Repo, h.Conn, h.Cfg)
}

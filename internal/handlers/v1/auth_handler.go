package v1

import (
	"context"
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	_ "github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/services/user"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	Repo *repository.Queries
}

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
		Repo: repo,
	}
}

func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	auth_router := router.PathPrefix("/auth").Subrouter()

	auth_router.HandleFunc("/register", h.Register).Methods(http.MethodPost)
	auth_router.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	auth_router.HandleFunc("/verify-email", h.VerifyEmail).Methods(http.MethodPost)
	auth_router.HandleFunc("/refresh", h.Refresh).Methods(http.MethodPost)
	auth_router.HandleFunc("/reset-password", h.ResetPassword).Methods(http.MethodPost)
	auth_router.HandleFunc("/change-password", h.ChangePassword).Methods(http.MethodPost)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	if h.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(
			utils.M{
				"message": "BAD_GATEWAY, No database connection, Oops!",
			},
		)
		return
	}

	user.CreateUser(wr, h.Repo)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Implement login logic here
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Implement email verification logic here
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

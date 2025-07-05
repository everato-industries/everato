package views

import (
	"context"
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/middlewares"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	BasePath string
	Repo     *repository.Queries
	Conn     *pgx.Conn
}

var _ handlers.Handler = (*AuthHandler)(nil) // Manually assert the interface implementation

func NewAuthHandler(basepath string) *AuthHandler {
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
		BasePath: basepath,
	}
}

func (a *AuthHandler) RegisterRoutes(router *mux.Router) {
	authRouter := router.PathPrefix(a.BasePath).Subrouter()

	guard := middlewares.NewAuthMiddleware(a.Repo, a.Conn, true)
	authRouter.Use(guard.Guard) // Guard the route

	authRouter.HandleFunc("/login", a.LoginPageHandler)
}

func (a *AuthHandler) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)
	if a.Repo == nil {
		wr.Status(http.StatusBadGateway).Json(
			utils.M{
				"message": "BAD_GATEWAY, No database connection, Oops!",
			},
		)
		return
	}

	wr.Html("templates/pages/login_page.html", utils.M{})
}

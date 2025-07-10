package admin

import (
	"net/http"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/jackc/pgx/v5"
)

func Login(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	// This function will handle the login logic for admin accounts.
	logger := pkg.NewLogger()
	defer logger.Close()

	// Get the DTO
	login_dto := &AdminLoginDTO{}
	err := wr.ParseBody(login_dto)
	if err != nil {
		logger.StdoutLogger.Error("Error parsing login request body", "err", err.Error())
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"message": "Invalid request body",
				"error":   err.Error(),
			},
		)
		return
	}

	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Logged in",
		},
	)
}

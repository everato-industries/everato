package admin

import (
	"net/http"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/jackc/pgx/v5"
)

func UpdateAdmin(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	// This function will handle the update of an existing admin account

	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Admin account updated",
		},
	)
}

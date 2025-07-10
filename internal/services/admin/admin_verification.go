package admin

import (
	"net/http"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/jackc/pgx/v5"
)

func SendVerificationEmail(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	// This function will handle sending a verification email to the admin
	// after they have created their account or updated their email address.

	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Verification email sent",
		},
	)
}

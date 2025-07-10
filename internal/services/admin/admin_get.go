package admin

import (
	"net/http"

	"github.com/dtg-lucifer/everato/config"
	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/jackc/pgx/v5"
)

func GetAllAdmins(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	// This function will handle fetching all admin accounts
	// It will return a list of all admins in the system
	wr.Status(http.StatusOK).Json(
		utils.M{
			"admins": "",
		},
	)
}

func GetAdminByID(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	// This function will handle fetching a specific admin account by ID
	// It will return the details of the admin with the given ID
	wr.Status(http.StatusOK).Json(
		utils.M{
			"admin": "",
		},
	)
}

func GetAdminByUserName(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	// This function will handle fetching a specific admin account by username
	// It will return the details of the admin with the given username
	wr.Status(http.StatusOK).Json(
		utils.M{
			"admin": "",
		},
	)
}

func SearchAdminByQuery(wr *utils.HttpWriter, repo *repository.Queries, conn *pgx.Conn, cfg *config.Config) {
	// This function will handle searching for admin accounts based on a query
	// It will return a list of admins that match the search criteria
	wr.Status(http.StatusOK).Json(
		utils.M{
			"admins": "",
		},
	)
}

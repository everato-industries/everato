package auth

import (
	"net/http"

	"github.com/dtg-lucifer/everato/api/internal/db/repository"
	"github.com/dtg-lucifer/everato/api/internal/utils"
)

func CreateUser(wr *utils.HttpWriter, repo *repository.Queries) {
	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Register route",
		},
	)
}

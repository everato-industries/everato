package admin

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/utils"
)

func CreateSubAdmin(w http.ResponseWriter, r *http.Request) {
	// This function will handle the creation of sub-admin accounts.
	wr := utils.NewHttpWriter(w, r)

	wr.Status(http.StatusOK).Json(
		utils.M{
			"message": "Sub-admin created",
		},
	)
}

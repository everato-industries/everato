package user

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/db/repository"
	"github.com/dtg-lucifer/everato/internal/utils"
)

func CreateUser(wr *utils.HttpWriter, repo *repository.Queries) {
	body := make(map[string]any)
	err := wr.ParseBody(&body)
	if err != nil {
		wr.Status(http.StatusBadRequest).Json(
			utils.M{
				"error": err.Error(),
			},
		)
		return
	}

	// TODO: Implement actual user creation logic here
	// For example:
	// email, ok := body["email"].(string)
	// if !ok {
	//     wr.Status(http.StatusBadRequest).Json(utils.M{"error": "Email is required"})
	//     return
	// }

	// For now, just returning a success response
	wr.Status(http.StatusCreated).Json(
		utils.M{
			"success": true,
			"message": "User registration endpoint reached successfully",
			"data":    body,
		},
	)
}

package general

import (
	"net/http"

	"github.com/dtg-lucifer/everato/server/internal/utils"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	wr.Status(http.StatusOK).Json(utils.M{
		"status": "success",
		"data":   "Server is running perfectly fine",
	})
}

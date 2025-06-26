package v1

import (
	"net/http"

	"github.com/dtg-lucifer/everato/api/internal/utils"
	"github.com/gorilla/mux"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type HealthCheckHandler struct {
	BasePath string
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{
		BasePath: "/health",
	}
}

func (h *HealthCheckHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(h.BasePath, GET_HealthCheck).Methods(http.MethodGet)
}

func GET_HealthCheck(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)
	requestId := w.Header().Get("X-Request-ID")

	if requestId == "" {
		wr.Status(http.StatusInternalServerError).Json(utils.M{
			"status": "error",
			"data":   "Request ID not found, aborting...",
		})
		return
	}

	wr.Status(http.StatusOK).Json(utils.M{
		"status":    "success",
		"data":      "Server is running perfectly fine",
		"requestId": requestId,
	})
}

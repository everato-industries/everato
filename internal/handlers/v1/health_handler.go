package v1

import (
	"net/http"

	_ "github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/utils"
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
	router.HandleFunc(h.BasePath, HealthCheck).Methods(http.MethodGet)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
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

package api

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/gorilla/mux"
)

// Package v1 provides the implementation of the health check handler
//
// Route prefix:
//   - /health
//
// Routes:
//   - /health
type HealthCheckHandler struct {
	BasePath string
}

// Asserting the implementation of the handler interface
var _ handlers.Handler = (*HealthCheckHandler)(nil)

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

	wr.Status(http.StatusOK).Json(utils.M{
		"status": "success",
		"data":   "Server is running perfectly fine",
	})
}

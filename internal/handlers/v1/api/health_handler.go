// Package api provides handlers for the REST API endpoints of the Everato application.
// This package contains all the HTTP handlers for API routes under /api/v1/
package api

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/gorilla/mux"
)

// HealthCheckHandler manages the health check endpoint for the application.
// It provides a simple status endpoint that can be used by monitoring tools,
// load balancers, and container orchestration systems to verify service health.
//
// Route prefix:
//   - `/api/vX/health`
//
// Routes:
//   - / - Returns a success response if the server is running (GET)
type HealthCheckHandler struct {
	BasePath string // Base URL path for health check endpoint
}

// Asserting the implementation of the handler interface
var _ handlers.Handler = (*HealthCheckHandler)(nil)

// NewHealthCheckHandler creates and initializes a new HealthCheckHandler instance.
// It configures the base path for the health check endpoint.
//
// Returns:
//   - A fully initialized HealthCheckHandler instance ready to register routes
func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{
		BasePath: "/health",
	}
}

// RegisterRoutes registers the health check route with the provided router.
// It maps the GET method on the health endpoint to the HealthCheck handler function.
//
// Parameters:
//   - router: The main router to attach the health check route to
func (h *HealthCheckHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(h.BasePath, HealthCheck).Methods(http.MethodGet)
}

// HealthCheck handles health check requests and returns a success response if the server is running.
// This function provides a simple status endpoint that external services can use to verify
// that the application is operational.
//
// HTTP Method: GET
// Route: /api/v1/health
//
// Response:
//   - 200 OK with JSON indicating successful health status
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	wr := utils.NewHttpWriter(w, r)

	wr.Status(http.StatusOK).Json(utils.M{
		"status": "success",
		"data":   "Server is running perfectly fine",
	})
}

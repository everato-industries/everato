// Package api provides handlers for the REST API endpoints of the Everato application.
// This package contains all the HTTP handlers for API routes under /api/v1/
package api

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/gorilla/mux"
)

// NotFoundHandler provides a custom 404 Not Found response handler for the application.
// It returns a standardized JSON response when users attempt to access non-existent routes.
//
// This handler ensures consistent error responses throughout the API by following
// the same response format as other API endpoints.
type NotFoundHandler struct{}

// Asserting the implementation of the handler interface
var _ handlers.Handler = (*NotFoundHandler)(nil)

// NewNotFoundHandler creates and initializes a new NotFoundHandler instance.
//
// Returns:
//   - A fully initialized NotFoundHandler instance ready to register with the router
func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

// RegisterRoutes configures the router's NotFoundHandler to use our custom handler.
// This method is called during application startup to set up the 404 handler.
//
// Instead of registering a standard route, this method sets the router's
// NotFoundHandler property to intercept all requests to undefined routes.
//
// Parameters:
//   - router: The main router to attach the 404 handler to
func (n *NotFoundHandler) RegisterRoutes(router *mux.Router) {
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := utils.NewHttpWriter(w, r)

		// Return a standardized 404 response in JSON format
		// This maintains consistency with other API responses
		wr.Status(http.StatusNotFound).Json(
			utils.M{
				"message": "Can't find the route you are looking for :)",
			},
		)
	})
}

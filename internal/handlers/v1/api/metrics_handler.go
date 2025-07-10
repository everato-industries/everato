// Package api provides handlers for the REST API endpoints of the Everato application.
// This package contains all the HTTP handlers for API routes under /api/v1/
package api

import (
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler manages Prometheus metrics endpoints for application monitoring.
// It exposes metrics that can be scraped by Prometheus for real-time observability.
//
// Route prefix:
//   - `/api/vX/metrics/`
//
// Routes:
//   - / - Exposes all registered Prometheus metrics (GET)
type MetricsHandler struct {
	BasePath string // Base URL path for metrics endpoints
}

// NewMetricsHandler creates and initializes a new MetricsHandler instance.
// It configures the base path for metrics endpoints.
//
// Returns:
//   - A fully initialized MetricsHandler instance ready to register routes
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		BasePath: "/metrics",
	}
}

// -----------------------------------------------
var _ handlers.Handler = (*MetricsHandler)(nil) // Assert the interface implementation to catch errors
// -----------------------------------------------

// RegisterRoutes registers all metrics-related routes with the provided router.
// It creates a subrouter with the base path and maps the Prometheus HTTP handler
// for metrics collection and exposure.
//
// This exposes standard Go runtime metrics plus any custom metrics registered
// throughout the application via the Prometheus client library.
//
// Parameters:
//   - router: The main router to attach metrics routes to
func (h *MetricsHandler) RegisterRoutes(router *mux.Router) {
	// Create a subrouter for metrics endpoints
	metrics := router.PathPrefix(h.BasePath).Subrouter()

	// Register the Prometheus HTTP handler for metrics collection
	// This handler automatically gathers all registered metrics on request
	metrics.Handle("/", promhttp.Handler()) // Export all metrics
}

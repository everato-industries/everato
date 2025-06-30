package v1

import (
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics API
type MetricsHandler struct {
	BasePath string
}

// NewMetricsHandler creates a new MetricsHandler instance.
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		BasePath: "/metrics",
	}
}

// -----------------------------------------------
var _ handlers.Handler = (*MetricsHandler)(nil) // Assert the interface implementation to catch errors
// -----------------------------------------------

// RegisterRoutes registers the metrics-related routes with the given router.
func (h *MetricsHandler) RegisterRoutes(router *mux.Router) {
	// Create a subrouter
	metrics := router.PathPrefix(h.BasePath).Subrouter()

	// Routes
	metrics.Handle("/", promhttp.Handler()) // Export all metrics
}

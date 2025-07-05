package views

import (
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/gorilla/mux"
)

type EventHandler struct {
	BasePath string
}

var _ handlers.Handler = (*EventHandler)(nil) // Manually asserting the implementation

func NewEventHandler(basePath string) *EventHandler {
	return &EventHandler{
		BasePath: basePath,
	}
}

func (h *EventHandler) RegisterRoutes(router *mux.Router) {
	// Register the event routes
	// eventRouter := router.PathPrefix(h.BasePath).Subrouter()

	// Define the routes for the event handler
	// eventRouter.HandleFunc("/", /* TODO */).Methods("GET")
	// eventRouter.HandleFunc("/search/{query}", /* TODO */).Methods("POST")
	// eventRouter.HandleFunc("/{slug}", /* TODO */).Methods("GET")
	// eventRouter.HandleFunc("/{id}", /* TODO */).Methods("GET")
	// eventRouter.HandleFunc("/new", /* TODO */).Methods("POST")
	// eventRouter.HandleFunc("/{id}", /* TODO */).Methods("PUT")
	// eventRouter.HandleFunc("/{id}", /* TODO */).Methods("DELETE")
}

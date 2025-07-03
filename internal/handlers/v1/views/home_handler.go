package views

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/dtg-lucifer/everato/templates/pages/home"
	"github.com/gorilla/mux"
)

// ViewsHandler is a struct that holds the templ instance for rendering views.
type ViewsHandler struct {
	BasePath string
}

func NewViewsHandler(basePath string) *ViewsHandler {
	return &ViewsHandler{
		BasePath: basePath,
	}
}

func (h *ViewsHandler) RegisterRoutes(router *mux.Router) {
	home_router := router.PathPrefix(h.BasePath).Subrouter()

	// Register the home routes
	home_router.HandleFunc("/", h.HomeRoute)
}

func (h *ViewsHandler) HomeRoute(w http.ResponseWriter, r *http.Request) {
	// Render the home page using the templ package
	templ.Handler(home.Home()).ServeHTTP(w, r)
}

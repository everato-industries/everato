package views

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/middlewares"
	"github.com/dtg-lucifer/everato/pages"
	"github.com/gorilla/mux"
)

// ViewsHandler is a struct that holds the templ instance for rendering views.
type ViewsHandler struct {
	BasePath string
}

var _ handlers.Handler = (*ViewsHandler)(nil) // Manually asserting the implementation

func NewViewsHandler(basePath string) *ViewsHandler {
	return &ViewsHandler{
		BasePath: basePath,
	}
}

func (h *ViewsHandler) RegisterRoutes(router *mux.Router) {
	home_router := router.PathPrefix(h.BasePath).Subrouter()

	guard := middlewares.NewAuthMiddleware(nil, nil, true) // No repo or connection needed for views
	home_router.Use(guard.Guard)

	// Register the home routes
	home_router.HandleFunc("/", h.HomeRoute)
}

func (h *ViewsHandler) HomeRoute(w http.ResponseWriter, r *http.Request) {
	// Render the home page using the templ package

	// wr := utils.NewHttpWriter(w, r)                      // DEMO FOR EXECUTING HTML TEMPLATES
	// wr.Html("templates/pages/home_page.html", utils.M{}) // DEMO FOR EXECUTING HTML TEMPLATES

	templ.Handler(pages.Home()).ServeHTTP(w, r)
}

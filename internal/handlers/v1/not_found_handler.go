package v1

import (
	"net/http"

	"github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/utils"
	"github.com/gorilla/mux"
)

type NotFoundHandler struct{}

var _ handlers.Handler = (*NotFoundHandler)(nil)

func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (n *NotFoundHandler) RegisterRoutes(router *mux.Router) {
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := utils.NewHttpWriter(w, r)

		wr.Status(http.StatusNotFound).Json(
			utils.M{
				"message": "Can't find the route you are looking for :)",
			},
		)
	})
}

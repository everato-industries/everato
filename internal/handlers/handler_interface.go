package handlers

import "github.com/gorilla/mux"

type Handler interface {
	// RegisterRoutes registers the routes for the handler
	// it should be called in the main function to register all the routes
	// with the router
	RegisterRoutes(router *mux.Router)
}

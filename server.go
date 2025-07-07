package main

import (
	"net/http"
	"strconv"

	"github.com/dtg-lucifer/everato/config"
	_ "github.com/dtg-lucifer/everato/internal/handlers"
	"github.com/dtg-lucifer/everato/internal/handlers/v1/api"
	"github.com/dtg-lucifer/everato/internal/handlers/v1/views"
	"github.com/dtg-lucifer/everato/internal/middlewares"
	"github.com/dtg-lucifer/everato/pkg"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Cfg    *config.Config
}

// Initialize the server
func NewServer(cfg *config.Config) *Server {
	router := mux.NewRouter()
	logger := pkg.NewLogger()
	defer logger.Close() // Ensure the logger is closed when the server is done

	server := &Server{
		Router: router,
		Cfg:    cfg,
	}

	server.initializeStaticFS()    // Initialize the static file system to serve files
	server.initializeViews()       // Initialize the handlers to handle the html views
	server.initializeMiddlewares() // Initialize the middlewares for the server
	server.initializeRoutes()      // Initialize the routes for the server

	return server
}

func (s *Server) initializeMiddlewares() {
	s.Router.Use(middlewares.RequestIDMiddleware)
	s.Router.Use(middlewares.CorsMiddleware)
	s.Router.Use(middlewares.LoggerMiddleware)

	// Add a 30 second timeout to all of the routes
	s.Router.Use(middlewares.TimeoutMiddleware(s.Cfg.RequestTimeout))
}

func (s *Server) initializeRoutes() {
	// Setting up the API prefix
	apivx := s.Router.PathPrefix(s.Cfg.ApiPrefix).Subrouter()

	// Route Group:
	// 	- General
	api.NewHealthCheckHandler().RegisterRoutes(apivx)
	api.NewMetricsHandler().RegisterRoutes(apivx)

	// Route Group:
	// 	- Authentication
	api.NewAuthHandler().RegisterRoutes(apivx)

	// Route Group:
	// 	- Events
	api.NewEventHandler().RegisterRoutes(apivx)

	// @TODO: User routes
	// @TODO: Ticket routes
	// @TODO: Notification routes

	// Notfound handler
	api.NewNotFoundHandler().RegisterRoutes(apivx)
}

func (s *Server) initializeStaticFS() {
	// Serve static files from the public directory
	s.Router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", PublicFS()))
}

func (s *Server) initializeViews() {
	// Initialize views if needed
	// This can be used to render HTML templates or other view engines
	//
	//
	view_router := s.Router.PathPrefix("/").Subrouter()

	views.NewViewsHandler("/").RegisterRoutes(view_router)
	views.NewAuthHandler("/auth").RegisterRoutes(view_router)
	views.NewEventHandler("/events").RegisterRoutes(view_router)
}

func (s *Server) Start() error {
	logger := pkg.NewLogger()
	defer logger.Close() // Ensure the logger is closed when the server is done

	addr := ":" + strconv.Itoa(s.Cfg.Server.Port)
	logger.Info("Server started running on", "port", addr)
	return http.ListenAndServe(addr, s.Router)
}

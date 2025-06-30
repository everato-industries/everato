package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dtg-lucifer/everato/config"
	_ "github.com/dtg-lucifer/everato/internal/handlers"
	v1 "github.com/dtg-lucifer/everato/internal/handlers/v1"
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

	server.initializeMiddlewares()
	server.initializeRoutes()
	return server
}

func (s *Server) initializeMiddlewares() {
	s.Router.Use(middlewares.RequestIDMiddleware)
	s.Router.Use(middlewares.CorsMiddleware)
	s.Router.Use(middlewares.LoggerMiddleware)

	// Add a 30 second timeout to all of the routes
	s.Router.Use(middlewares.TimeoutMiddleware(time.Second * 10))
}

func (s *Server) initializeRoutes() {
	// Setting up the API prefix
	apiv1 := s.Router.PathPrefix("/api/v1").Subrouter()

	// Route Group:
	// 	- General
	v1.NewHealthCheckHandler().RegisterRoutes(apiv1)

	// Route Group:
	// 	- Authentication
	v1.NewAuthHandler().RegisterRoutes(apiv1)

	// Route Group:
	// 	- Events
	v1.NewEventHandler().RegisterRoutes(apiv1)

	// @TODO: User routes
	// @TODO: Ticket routes
	// @TODO: Notification routes

	// Notfound handler
	v1.NewNotFoundHandler().RegisterRoutes(apiv1)
}

func (s *Server) Start() error {
	logger := pkg.NewLogger()
	defer logger.Close() // Ensure the logger is closed when the server is done

	addr := ":" + strconv.Itoa(s.Cfg.Server.Port)
	logger.Info("Server started running on", "port", addr)
	return http.ListenAndServe(addr, s.Router)
}

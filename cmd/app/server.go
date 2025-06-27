package main

import (
	"net/http"
	"strconv"

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
}

func (s *Server) initializeRoutes() {
	// Setting up the API prefix
	apiv1 := s.Router.PathPrefix("/api/v1").Subrouter()

	v1.NewHealthCheckHandler().RegisterRoutes(apiv1)
	v1.NewAuthHandler().RegisterRoutes(apiv1)

	// TODO: Authentication routes
	// TODO: User routes
	// TODO: Event routes
	// TODO: Ticket routes
	// TODO: Notification routes
}

func (s *Server) Start() error {
	logger := pkg.NewLogger().StdoutLogger
	addr := ":" + strconv.Itoa(s.Cfg.Server.Port)
	logger.Info("Server started running on", "port", addr)
	return http.ListenAndServe(addr, s.Router)
}

package main

import (
	"net/http"
	"strconv"

	"github.com/dtg-lucifer/everato/server/config"
	"github.com/dtg-lucifer/everato/server/internal/middlewares"
	"github.com/dtg-lucifer/everato/server/internal/routes/general"
	"github.com/dtg-lucifer/everato/server/pkg/logger"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Cfg    *config.Config
	Logger *logger.Logger
}

// Initialize the server
func NewServer(cfg *config.Config) *Server {
	router := mux.NewRouter()
	logger := logger.NewLogger()

	server := &Server{
		Router: router,
		Cfg:    cfg,
		Logger: logger,
	}

	server.initializeMiddlewares()
	server.initializeRoutes()
	return server
}

func (s *Server) initializeMiddlewares() {
	s.Router.Use(middlewares.CorsMiddleware)
	s.Router.Use(middlewares.LoggerMiddleware)
}

func (s *Server) initializeRoutes() {
	// Setting up the API prefix
	apiv1 := s.Router.PathPrefix("/api/v1").Subrouter()

	apiv1.HandleFunc("/health", general.HealthCheckHandler).Methods("GET")

	// TODO: Authentication routes
	// TODO: User routes
	// TODO: Event routes
	// TODO: Ticket routes
	// TODO: Notification routes
}

func (s *Server) Start() error {
	addr := ":" + strconv.Itoa(s.Cfg.Server.Port)
	s.Logger.StdoutLogger.Info("Server started listening", "addr", addr)
	return http.ListenAndServe(addr, s.Router)
}

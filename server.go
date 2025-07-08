/*
Server - Core HTTP Server Implementation for Everato

This file implements the HTTP server for the Everato platform, handling route setup,
middleware configuration, and server startup. It follows a clean architecture approach
with separate API and view routes.
*/
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

// Server represents the HTTP server instance for the Everato application.
// It encapsulates the router and application configuration.
type Server struct {
	Router *mux.Router    // HTTP router that handles all incoming requests
	Cfg    *config.Config // Application configuration
}

// NewServer initializes and returns a new Server instance with the provided configuration.
// It sets up all necessary components like middlewares, routes, and static file handlers.
//
// Parameters:
//   - cfg: Application configuration containing server settings, database details, etc.
//
// Returns:
//   - A fully initialized Server instance ready to be started
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

// initializeMiddlewares configures all global middlewares for the server.
// These middlewares will be applied to all routes.
//
// Middleware chain:
// 1. RequestID - Adds a unique identifier to each request
// 2. CORS - Handles Cross-Origin Resource Sharing
// 3. Logger - Logs all HTTP requests
// 4. Timeout - Enforces request timeout limits from configuration
func (s *Server) initializeMiddlewares() {
	s.Router.Use(middlewares.RequestIDMiddleware)
	s.Router.Use(middlewares.CorsMiddleware)
	s.Router.Use(middlewares.LoggerMiddleware)

	// Add a timeout middleware based on the configuration
	s.Router.Use(middlewares.TimeoutMiddleware(s.Cfg.RequestTimeout))
}

// initializeRoutes sets up all API routes for the application.
// Routes are grouped by functionality and registered with the appropriate handlers.
//
// Route groups:
// - General: Health check and metrics endpoints
// - Authentication: User registration, login, verification
// - Events: Event creation, management, and search
//
// All API routes are prefixed with the configured API prefix (e.g., /api/v1)
func (s *Server) initializeRoutes() {
	// Setting up the API prefix
	apivx := s.Router.PathPrefix(s.Cfg.ApiPrefix).Subrouter()

	// Route Group: General
	// Health check and monitoring endpoints
	api.NewHealthCheckHandler().RegisterRoutes(apivx)
	api.NewMetricsHandler().RegisterRoutes(apivx)

	// Route Group: Authentication
	// User registration, login, verification, etc.
	api.NewAuthHandler(s.Cfg).RegisterRoutes(apivx)

	// Route Group: Events
	// Event creation, management, search, etc.
	api.NewEventHandler(s.Cfg).RegisterRoutes(apivx)

	// @TODO: User routes - User profile, management, etc.
	// @TODO: Ticket routes - Ticket creation, validation, etc.
	// @TODO: Notification routes - Email/push notifications, etc.

	// Register the NotFound handler for unmatched routes
	api.NewNotFoundHandler().RegisterRoutes(apivx)
}

// initializeStaticFS configures static file serving for the application.
// In development mode, it serves files directly from the filesystem.
// In production mode, it serves files from the embedded filesystem.
//
// Static files include CSS, JavaScript, images, and other assets.
func (s *Server) initializeStaticFS() {
	// Serve static files from the public directory
	// The PublicFS() function is defined differently in dev/prod builds
	s.Router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", PublicFS()))
}

// initializeViews sets up the server-side rendering routes for the application.
// These routes serve HTML pages rendered using the templ templating engine.
//
// View routes:
// - Root (/): Home page and general site pages
// - Auth (/auth): Authentication-related pages (login, register, etc.)
// - Events (/events): Event browsing and detail pages
func (s *Server) initializeViews() {
	// Create a subrouter for all view routes starting at the root path
	view_router := s.Router.PathPrefix("/").Subrouter()

	// Register view handlers for different sections of the application
	views.NewViewsHandler("/").RegisterRoutes(view_router)       // Home and general pages
	views.NewAuthHandler("/auth").RegisterRoutes(view_router)    // Authentication pages
	views.NewEventHandler("/events").RegisterRoutes(view_router) // Event pages
}

// Start begins listening for HTTP requests on the configured port.
// This method is blocking and will only return if there's an error starting the server.
//
// Returns:
//   - An error if the server fails to start or encounters a fatal error while running
func (s *Server) Start() error {
	logger := pkg.NewLogger()
	defer logger.Close() // Ensure the logger is closed when the server is done

	// Construct the server address from the configured port
	addr := ":" + strconv.Itoa(s.Cfg.Server.Port)
	logger.Info("Server started running on", "port", addr)

	// Start the HTTP server with the configured router
	// This call is blocking and will only return on error
	return http.ListenAndServe(addr, s.Router)
}

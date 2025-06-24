package main

import (
	"net/http"
	"strconv"

	"github.com/dtg-lucifer/everato/server/config"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Cfg    *config.Config
}

func NewServer(cfg *config.Config) *Server {
	router := mux.NewRouter()
	server := &Server{
		Router: router,
		Cfg:    cfg,
	}

	server.initializeMiddlewares()
	server.initializeRoutes()
	return server
}

func (s *Server) initializeMiddlewares() {
	s.Router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Log the request
			println("Request:", r.Method, r.URL.Path)

			// Call the next handler
			h.ServeHTTP(w, r)
		})
	})
}

func (s *Server) initializeRoutes() {
	s.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Everato!"))
	}).Methods("GET")
}

func (s *Server) Start() error {
	addr := ":" + strconv.Itoa(s.Cfg.Server.Port)
	return http.ListenAndServe(addr, s.Router)
}

// internal/server/server.go
package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xmonader/team-timetracker/internal/server/handlers"
	"gorm.io/gorm"
)

// Server encapsulates the server dependencies
type Server struct {
	Router          *mux.Router
	DB              *gorm.DB
	TrackingHandler *handlers.TrackingHandler
	EntriesHandler  *handlers.EntriesHandler
	HealthHandler   *handlers.HealthHandler
}

// NewServer initializes the server with routes and handlers
func NewServer(db *gorm.DB) *Server {
	router := mux.NewRouter()

	trackingHandler := &handlers.TrackingHandler{DB: db}
	entriesHandler := &handlers.EntriesHandler{DB: db}
	healthHandler := &handlers.HealthHandler{DB: db}

	server := &Server{
		Router:          router,
		DB:              db,
		TrackingHandler: trackingHandler,
		EntriesHandler:  entriesHandler,
		HealthHandler:   healthHandler,
	}

	server.setupRoutes()

	return server
}

// setupRoutes sets up the server routes
func (s *Server) setupRoutes() {
	// Tracking routes
	s.Router.HandleFunc("/api/start", s.TrackingHandler.StartTracking).Methods("POST")
	s.Router.HandleFunc("/api/stop", s.TrackingHandler.StopTracking).Methods("POST")

	// Entries routes
	s.Router.HandleFunc("/api/entries", s.EntriesHandler.GetEntries).Methods("GET")

	// Health and Liveness routes
	s.Router.HandleFunc("/live", s.HealthHandler.LivenessHandler).Methods("GET")
	s.Router.HandleFunc("/health", s.HealthHandler.HealthCheckHandler).Methods("GET")
}

// Run starts the HTTP server
func (s *Server) Run(addr string) error {
	log.Printf("Server is running on %s", addr)
	return http.ListenAndServe(addr, s.Router)
}

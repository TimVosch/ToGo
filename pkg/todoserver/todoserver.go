package todoserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server contains relevant todo Server objects
type Server struct {
	httpServer *http.Server
	router     *mux.Router
}

// NewServer creates a new server
func NewServer() *Server {
	// Set up
	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	// Build Server struct
	s := &Server{
		httpServer,
		router,
	}

	setRoutes(s)

	return s
}

// StartServer creates and starts a TodoServer
func (s *Server) StartServer() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown the server
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}

}

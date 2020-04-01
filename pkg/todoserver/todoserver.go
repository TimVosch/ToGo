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
	s := &Server{
		nil,
		mux.NewRouter(),
	}

	s.httpServer = &http.Server{
		Addr:    ":4000",
		Handler: s.router,
	}

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

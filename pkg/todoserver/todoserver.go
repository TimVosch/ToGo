package todoserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// TodoServer contains relevant todo Server objects
type TodoServer struct {
	httpServer *http.Server
	router     *mux.Router
	db         TodoDAL
}

// NewServer creates a new server
func NewServer() *TodoServer {
	// Set up
	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	db := NewMemoryDB()

	// Build TodoServer struct
	s := &TodoServer{
		httpServer,
		router,
		db,
	}

	setRoutes(s)

	return s
}

// StartServer creates and starts a TodoServer
func (s *TodoServer) StartServer() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown the server
func (s *TodoServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}

}

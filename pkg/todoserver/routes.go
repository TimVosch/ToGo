package todoserver

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleHealthCheck() http.HandlerFunc {
	// Create response
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]bool{"healthy": true}
		json.NewEncoder(w).Encode(response)
	}
}

func (s *Server) handleGetTodos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos := s.db.GetTodosForUser(0)
		json.NewEncoder(w).Encode(todos)
	}
}

func setRoutes(s *Server) {
	r := s.router

	r.HandleFunc("/health", s.handleHealthCheck()).Methods("GET")
	r.HandleFunc("/todos/{id:[0-9]+}", s.handleGetTodos()).Methods("GET")
}

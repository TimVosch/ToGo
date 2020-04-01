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

func (s *Server) handleCreateTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo TodoEntry
		json.NewDecoder(r.Body).Decode(&todo)

		// Insert
		if err := s.db.InsertTodo(todo); err != nil {
			// Error!
		}
	}
}

func setRoutes(s *Server) {
	r := s.router

	r.HandleFunc("/health", s.handleHealthCheck()).Methods("GET")
	r.HandleFunc("/todos", s.handleGetTodos()).Methods("GET")
	r.HandleFunc("/todos", s.handleCreateTodo()).Methods("POST")
}

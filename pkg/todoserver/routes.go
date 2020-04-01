package todoserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func (s *Server) handleGetTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.ParseInt(params["id"], 10, 0)
		todo := s.db.GetTodoByID(int(id))
		json.NewEncoder(w).Encode(todo)
	}
}

func (s *Server) handleCreateTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo TodoEntry
		json.NewDecoder(r.Body).Decode(&todo)

		// Insert
		created, err := s.db.InsertTodo(todo)
		if err != nil {
			// Error!
		}

		json.NewEncoder(w).Encode(created)
	}
}

func (s *Server) handleDeleteTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.ParseInt(params["id"], 10, 0)
		err := s.db.DeleteTodo(int(id))

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func setRoutes(s *Server) {
	r := s.router

	r.HandleFunc("/health", s.handleHealthCheck()).Methods("GET")
	r.HandleFunc("/todos/{id:[0-9]+}", s.handleGetTodo()).Methods("GET")
	r.HandleFunc("/todos/{id:[0-9]+}", s.handleDeleteTodo()).Methods("DELETE")
	r.HandleFunc("/todos", s.handleGetTodos()).Methods("GET")
	r.HandleFunc("/todos", s.handleCreateTodo()).Methods("POST")
}

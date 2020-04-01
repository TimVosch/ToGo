package todoserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/timvosch/togo/pkg/api"

	"github.com/gorilla/mux"
)

func (s *TodoServer) handleHealthCheck() http.HandlerFunc {
	// Create response
	return func(w http.ResponseWriter, r *http.Request) {
		body := map[string]interface{}{"healthy": true}
		api.SendResponse(w, http.StatusOK, body, "Everything is O.K.")
	}
}

func (s *TodoServer) handleGetTodos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos := s.db.GetTodosForUser(0)
		api.SendResponse(w, http.StatusOK, todos, "Returned all Todos for given user")
	}
}

func (s *TodoServer) handleGetTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.ParseInt(params["id"], 10, 0)

		todo := s.db.GetTodoByID(int(id))

		api.SendResponse(w, http.StatusOK, todo, "Returned Todo with given id")
	}
}

func (s *TodoServer) handleCreateTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo TodoEntry
		json.NewDecoder(r.Body).Decode(&todo)

		// Insert
		created, err := s.db.InsertTodo(todo)
		if err != nil {
			api.SendResponse(w, http.StatusInternalServerError, nil, "An error occured while creating Todo")
			return
		}

		api.SendResponse(w, http.StatusOK, created, "Created a new Todo")
	}
}

func (s *TodoServer) handleDeleteTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.ParseInt(params["id"], 10, 0)

		err := s.db.DeleteTodo(int(id))
		if err != nil {
			api.SendResponse(w, http.StatusNotFound, nil, "An error occured while deleting Todo")
			return
		}

		body := map[string]string{
			"message": "Todo has been removed",
		}
		api.SendResponse(w, http.StatusOK, body, "Deleted Todo with given ID")
	}
}

func setRoutes(s *TodoServer) {
	r := s.router

	r.HandleFunc("/health", s.handleHealthCheck()).Methods("GET")
	r.HandleFunc("/todos/{id:[0-9]+}", s.handleGetTodo()).Methods("GET")
	r.HandleFunc("/todos/{id:[0-9]+}", s.handleDeleteTodo()).Methods("DELETE")
	r.HandleFunc("/todos", s.handleGetTodos()).Methods("GET")
	r.HandleFunc("/todos", s.handleCreateTodo()).Methods("POST")
}

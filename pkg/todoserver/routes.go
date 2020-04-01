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
		res := api.NewResponse(w)
		res.Body = map[string]interface{}{"healthy": true}
		res.Send()
	}
}

func (s *TodoServer) handleGetTodos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := api.NewResponse(w)
		res.Body = s.db.GetTodosForUser(0)
		res.Send()
	}
}

func (s *TodoServer) handleGetTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.ParseInt(params["id"], 10, 0)
		res := api.NewResponse(w)
		res.Body = s.db.GetTodoByID(int(id))
		res.Send()
	}
}

func (s *TodoServer) handleCreateTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo TodoEntry
		json.NewDecoder(r.Body).Decode(&todo)

		res := api.NewResponse(w)

		// Insert
		created, err := s.db.InsertTodo(todo)
		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Send()
			return
		}

		res.Body = created
		res.Send()
	}
}

func (s *TodoServer) handleDeleteTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := api.NewResponse(w)
		params := mux.Vars(r)
		id, _ := strconv.ParseInt(params["id"], 10, 0)
		err := s.db.DeleteTodo(int(id))

		if err != nil {
			res.Send()
			res.Status = http.StatusNotFound
			return
		}

		res.Body = map[string]string{
			"message": "Todo has been removed",
		}
		res.Send()
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

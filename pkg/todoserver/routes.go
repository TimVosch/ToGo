package todoserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/timvosch/togo/pkg/api"

	"github.com/gorilla/mux"
)

func (s *TodoServer) handleHealthCheck() api.HandlerFunc {
	// Create response
	return func(ctx *api.CTX, next func()) {
		body := map[string]interface{}{"healthy": true}
		api.SendResponse(ctx.W, http.StatusOK, body, "Everything is O.K.")
	}
}

func (s *TodoServer) handleGetTodos() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		todos := s.db.GetTodosForUser(0)
		api.SendResponse(ctx.W, http.StatusOK, todos, "Returned all Todos for given user")
	}
}

func (s *TodoServer) handleGetTodo() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		params := mux.Vars(ctx.R)
		id, _ := strconv.ParseInt(params["id"], 10, 0)

		todo := s.db.GetTodoByID(int(id))

		api.SendResponse(ctx.W, http.StatusOK, todo, "Returned Todo with given id")
	}
}

func (s *TodoServer) handleCreateTodo() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		var todo TodoEntry
		json.NewDecoder(ctx.R.Body).Decode(&todo)

		// Insert
		created, err := s.db.InsertTodo(todo)
		if err != nil {
			api.SendResponse(ctx.W, http.StatusInternalServerError, nil, "An error occured while creating Todo")
			return
		}

		api.SendResponse(ctx.W, http.StatusOK, created, "Created a new Todo")
	}
}

func (s *TodoServer) handleDeleteTodo() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		params := mux.Vars(ctx.R)
		id, _ := strconv.ParseInt(params["id"], 10, 0)

		err := s.db.DeleteTodo(int(id))
		if err != nil {
			api.SendResponse(ctx.W, http.StatusNotFound, nil, "An error occured while deleting Todo")
			return
		}

		body := map[string]string{
			"message": "Todo has been removed",
		}
		api.SendResponse(ctx.W, http.StatusOK, body, "Deleted Todo with given ID")
	}
}

func setRoutes(s *TodoServer) {
	r := s.router

	r.HandleFunc(
		"/health",
		api.Handler(
			s.handleHealthCheck(),
		),
	).Methods("GET")
	r.HandleFunc(
		"/todos/{id:[0-9]+}",
		api.Handler(
			s.handleGetTodo(),
		),
	).Methods("GET")
	r.HandleFunc(
		"/todos/{id:[0-9]+}",
		api.Handler(
			s.handleDeleteTodo(),
		),
	).Methods("DELETE")
	r.HandleFunc(
		"/todos",
		api.Handler(
			s.handleGetTodos(),
		),
	).Methods("GET")
	r.HandleFunc(
		"/todos",
		api.Handler(
			s.handleCreateTodo(),
		),
	).Methods("POST")
}

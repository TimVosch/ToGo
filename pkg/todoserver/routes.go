package todoserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/timvosch/togo/pkg/api"
	"github.com/timvosch/togo/pkg/common/middleware"

	"github.com/gorilla/mux"
)

func (ts *TodoServer) handleHealthCheck() api.HandlerFunc {
	// Create response
	return func(ctx *api.CTX, next func()) {
		body := map[string]interface{}{"healthy": true}
		ctx.SendResponse(http.StatusOK, body, "Everything is O.K.")
	}
}

func (ts *TodoServer) handleGetTodos() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		todos := ts.db.GetTodosForUser(0)
		ctx.SendResponse(http.StatusOK, todos, "Returned all Todos for given user")
	}
}

func (ts *TodoServer) handleGetTodo() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		params := mux.Vars(ctx.R)
		id, _ := strconv.ParseInt(params["id"], 10, 0)

		todo := ts.db.GetTodoByID(int(id))

		ctx.SendResponse(http.StatusOK, todo, "Returned Todo with given id")
	}
}

func (ts *TodoServer) handleCreateTodo() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		var todo TodoEntry
		json.NewDecoder(ctx.R.Body).Decode(&todo)

		// Set owner
		id, ok := ctx.User["sub"].(float64)
		if !ok {
			ctx.SendResponse(http.StatusUnauthorized, nil, "Authenticated is not a user")
			return
		}
		todo.OwnerID = int(id)

		// Insert
		created, err := ts.db.InsertTodo(todo)
		if err != nil {
			ctx.SendResponse(http.StatusInternalServerError, nil, "An error occured while creating Todo")
			return
		}

		ctx.SendResponse(http.StatusOK, created, "Created a new Todo")
	}
}

func (ts *TodoServer) handleDeleteTodo() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		params := mux.Vars(ctx.R)
		id, _ := strconv.ParseInt(params["id"], 10, 0)

		err := ts.db.DeleteTodo(int(id))
		if err != nil {
			ctx.SendResponse(http.StatusNotFound, nil, "An error occured while deleting Todo")
			return
		}

		body := map[string]string{
			"message": "Todo has been removed",
		}
		ctx.SendResponse(http.StatusOK, body, "Deleted Todo with given ID")
	}
}

func setRoutes(ts *TodoServer) {
	r := ts.router
	auth := middleware.MakeAuth(ts.jwt)

	//
	r.HandleFunc(
		"/health",
		api.Handler(
			ts.handleHealthCheck(),
		),
	).Methods("GET")

	//
	r.HandleFunc(
		"/todos/{id:[0-9]+}",
		api.Handler(
			auth(),
			ts.handleGetTodo(),
		),
	).Methods("GET")

	//
	r.HandleFunc(
		"/todos/{id:[0-9]+}",
		api.Handler(
			auth(),
			ts.handleDeleteTodo(),
		),
	).Methods("DELETE")

	//
	r.HandleFunc(
		"/todos",
		api.Handler(
			auth(),
			ts.handleGetTodos(),
		),
	).Methods("GET")

	//
	r.HandleFunc(
		"/todos",
		api.Handler(
			auth(),
			ts.handleCreateTodo(),
		),
	).Methods("POST")
}

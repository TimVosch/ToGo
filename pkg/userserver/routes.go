package userserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/timvosch/togo/pkg/api"
	"github.com/timvosch/togo/pkg/jwt"
)

func (us *UserServer) handleHealthCheck() api.HandlerFunc {
	return func(ctx *api.CTX) {
		body := map[string]interface{}{
			"healthy": true,
		}
		api.SendResponse(ctx.W, http.StatusOK, body, "Everything is O.K.")
	}
}

func (us *UserServer) handleRegisterUser() api.HandlerFunc {
	return func(ctx *api.CTX) {
		var user User
		json.NewDecoder(ctx.R.Body).Decode(&user)

		// Hash password
		user.SetPassword(user.Password)

		created, err := us.repo.InsertUser(user)
		if err != nil {
			api.SendResponse(ctx.W, http.StatusInternalServerError, err, "An error occured while registering user")
		}

		api.SendResponse(ctx.W, http.StatusCreated, created, "User registered")
	}
}

func (us *UserServer) handleGetUserByID() api.HandlerFunc {
	return func(ctx *api.CTX) {
		params := mux.Vars(ctx.R)
		id, _ := strconv.ParseInt(params["id"], 10, 0)

		user := us.repo.GetUserByID(int(id))
		if user == nil {
			api.SendResponse(ctx.W, http.StatusNotFound, nil, "User not found")
			return
		}

		api.SendResponse(ctx.W, http.StatusOK, user, "Fetched user by ID")
	}
}

func (us *UserServer) handleLogin() api.HandlerFunc {
	return func(ctx *api.CTX) {
		token := us.jwt.CreateToken()
		body := map[string]string{
			"token": us.jwt.Sign(token),
		}
		api.SendResponse(ctx.W, http.StatusOK, body, "Succesfully logged in")
	}
}

func makeAuthMiddleware(jwt *jwt.JWT) func(api.HandlerFunc) api.HandlerFunc {
	// Middleware
	return func(next api.HandlerFunc) api.HandlerFunc {
		// Handler
		return func(ctx *api.CTX) {
			r := ctx.R
			w := ctx.W
			header := r.Header.Get("Authorization")
			parts := strings.Split(header, " ")
			if len(parts) != 2 {
				api.SendResponse(w, http.StatusUnauthorized, nil, "Must be authenticated")
				return
			}
			if parts[0] != "Bearer" {
				api.SendResponse(w, http.StatusUnauthorized, nil, "Authorization method not supported")
				return
			}

			_, err := jwt.Verify(parts[1])
			if err != nil {
				api.SendResponse(w, http.StatusUnauthorized, nil, "Provided JWT is invalid")
				return
			}

			// Continue
			next(ctx)
		}
	}
}

func setRoutes(us *UserServer) {
	r := us.router
	auth := makeAuthMiddleware(us.jwt)

	r.HandleFunc(
		"/health",
		api.Handler(auth(us.handleHealthCheck())),
	).Methods("GET")
	r.HandleFunc(
		"/auth",
		api.Handler(us.handleLogin()),
	).Methods("POST")
	r.HandleFunc(
		"/users",
		api.Handler(us.handleRegisterUser()),
	).Methods("POST")
	r.HandleFunc(
		"/users/{id:[0-9]+}",
		api.Handler(auth(us.handleGetUserByID())),
	).Methods("GET")
}

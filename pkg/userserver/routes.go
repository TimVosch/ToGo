package userserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/timvosch/togo/pkg/api"
	"github.com/timvosch/togo/pkg/jwt"
)

func (us *UserServer) handleHealthCheck() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		body := map[string]interface{}{
			"healthy": true,
		}
		api.SendResponse(ctx.W, http.StatusOK, body, "Everything is O.K.")
	}
}

func (us *UserServer) handleRegisterUser() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
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

func (us *UserServer) handleGetUserSelf() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		// Get uid from authenticated user
		idf, ok := ctx.User["sub"].(float64)
		id := int(idf)

		if ok == false {
			api.SendResponse(ctx.W, http.StatusForbidden, ctx.User, "Not authenticated as user")
			return
		}

		user := us.repo.GetUserByID(id)
		if user == nil {
			api.SendResponse(ctx.W, http.StatusNotFound, nil, "User not found")
			return
		}

		api.SendResponse(ctx.W, http.StatusOK, user, "Fetched user by ID")
	}
}

func (us *UserServer) handleLogin() api.HandlerFunc {
	// Request DTO
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(ctx *api.CTX, next func()) {
		var req Request
		if err := json.NewDecoder(ctx.R.Body).Decode(&req); err != nil {
			api.SendResponse(ctx.W, http.StatusBadRequest, nil, "Bad request")
			return
		}

		token, err := us.loginUser(req.Email, req.Password)
		if err != nil {
			api.SendResponse(ctx.W, http.StatusForbidden, nil, "Incorrect login information")
			return
		}

		body := map[string]string{
			"token": token,
		}
		api.SendResponse(ctx.W, http.StatusOK, body, "Succesfully logged in")
	}
}

func makeAuthMiddleware(jwt *jwt.JWT) func() api.HandlerFunc {
	// Handler
	return func() api.HandlerFunc {
		return func(ctx *api.CTX, next func()) {
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

			token, err := jwt.Verify(parts[1])
			if err != nil {
				api.SendResponse(w, http.StatusUnauthorized, nil, "Provided JWT is invalid")
				return
			}

			ctx.User = token.Body
			next()
		}
	}
}

func setRoutes(us *UserServer) {
	r := us.router
	auth := makeAuthMiddleware(us.jwt)

	//
	r.HandleFunc(
		"/health",
		api.Handler(
			auth(),
			us.handleHealthCheck(),
		),
	).Methods("GET")

	//
	r.HandleFunc(
		"/auth",
		api.Handler(us.handleLogin()),
	).Methods("POST")

	//
	r.HandleFunc(
		"/users",
		api.Handler(us.handleRegisterUser()),
	).Methods("POST")

	//
	r.HandleFunc(
		"/users/self",
		api.Handler(
			auth(),
			us.handleGetUserSelf(),
		),
	).Methods("GET")
}

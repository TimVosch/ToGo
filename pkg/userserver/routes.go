package userserver

import (
	"encoding/json"
	"net/http"

	"github.com/timvosch/togo/pkg/api"
	"github.com/timvosch/togo/pkg/common/middleware"
)

func (us *UserServer) handleHealthCheck() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		body := map[string]interface{}{
			"healthy": true,
		}
		ctx.SendResponse(http.StatusOK, body, "Everything is O.K.")
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
			ctx.SendResponse(http.StatusBadRequest, nil, err.Error())
			return
		}

		ctx.SendResponse(http.StatusCreated, created, "User registered")
	}
}

func (us *UserServer) handleGetUserSelf() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		user := us.repo.GetUserByID(*ctx.User.ID)
		if user == nil {
			ctx.SendResponse(http.StatusNotFound, nil, "User not found")
			return
		}

		ctx.SendResponse(http.StatusOK, user, "Fetched user by ID")
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
			ctx.SendResponse(http.StatusBadRequest, nil, "Bad request")
			return
		}

		token, err := us.loginUser(req.Email, req.Password)
		if err != nil {
			ctx.SendResponse(http.StatusForbidden, nil, "Incorrect login information")
			return
		}

		body := map[string]string{
			"token": token,
		}
		ctx.SendResponse(http.StatusOK, body, "Succesfully logged in")
	}
}

func (us *UserServer) handleGetJWKS() api.HandlerFunc {
	return func(ctx *api.CTX, next func()) {
		json.NewEncoder(ctx.W).Encode(us.jwks)
	}
}

func setRoutes(us *UserServer) {
	r := us.router
	auth := middleware.MakeAuth(us.verifier)

	r.HandleFunc(
		"/.well-known/jwks.json",
		api.Handler(
			us.handleGetJWKS(),
		),
	)

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

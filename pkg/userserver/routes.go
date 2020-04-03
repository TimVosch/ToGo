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

func (us *UserServer) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := map[string]interface{}{
			"healthy": true,
		}
		api.SendResponse(w, http.StatusOK, body, "Everything is O.K.")
	}
}

func (us *UserServer) handleRegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		// Hash password
		user.SetPassword(user.Password)

		created, err := us.repo.InsertUser(user)
		if err != nil {
			api.SendResponse(w, http.StatusInternalServerError, err, "An error occured while registering user")
		}

		api.SendResponse(w, http.StatusCreated, created, "User registered")
	}
}

func (us *UserServer) handleGetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, _ := strconv.ParseInt(params["id"], 10, 0)

		user := us.repo.GetUserByID(int(id))
		if user == nil {
			api.SendResponse(w, http.StatusNotFound, nil, "User not found")
			return
		}

		api.SendResponse(w, http.StatusOK, user, "Fetched user by ID")
	}
}

func (us *UserServer) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := us.jwt.CreateToken()
		body := map[string]string{
			"token": us.jwt.Sign(token),
		}
		api.SendResponse(w, http.StatusOK, body, "Succesfully logged in")
	}
}

func makeAuthMiddleware(jwt *jwt.JWT) func(http.HandlerFunc) http.HandlerFunc {
	// Middleware
	return func(h http.HandlerFunc) http.HandlerFunc {
		// Handler
		return func(w http.ResponseWriter, r *http.Request) {
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
			h(w, r)
		}
	}
}

func setRoutes(us *UserServer) {
	r := us.router
	auth := makeAuthMiddleware(us.jwt)

	r.HandleFunc(
		"/health",
		auth(us.handleHealthCheck()),
	).Methods("GET")
	r.HandleFunc("/auth", us.handleLogin()).Methods("POST")
	r.HandleFunc("/users", us.handleRegisterUser()).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", auth(us.handleGetUserByID())).Methods("GET")
}

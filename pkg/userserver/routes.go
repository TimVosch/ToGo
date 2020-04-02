package userserver

import (
	"encoding/json"
	"net/http"

	"github.com/timvosch/togo/pkg/api"
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

		api.SendResponse(w, http.StatusOK, created, "User registered")
	}
}

func setRoutes(us *UserServer) {
	r := us.router

	r.HandleFunc("/health", us.handleHealthCheck()).Methods("GET")
	r.HandleFunc("/users", us.handleRegisterUser()).Methods("POST")
}

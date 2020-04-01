package userserver

import (
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

func setRoutes(us *UserServer) {
	r := us.router

	r.HandleFunc("/health", us.handleHealthCheck()).Methods("GET")
}

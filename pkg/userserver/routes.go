package userserver

import (
	"encoding/json"
	"net/http"
)

func (us *UserServer) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{
			"healthy": true,
		})
	}
}

func setRoutes(us *UserServer) {
	r := us.router

	r.HandleFunc("/health", us.handleHealthCheck()).Methods("GET")
}

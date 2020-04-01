package userserver

import (
	"net/http"

	"github.com/timvosch/togo/pkg/api"
)

func (us *UserServer) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := api.NewResponse(w)
		res.Body = map[string]interface{}{
			"healthy": true,
		}
		res.Send()
	}
}

func setRoutes(us *UserServer) {
	r := us.router

	r.HandleFunc("/health", us.handleHealthCheck()).Methods("GET")
}

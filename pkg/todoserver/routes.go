package todoserver

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleHealthCheck() http.HandlerFunc {
	// Create response
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]bool{"healthy": true}
		json.NewEncoder(w).Encode(response)
	}
}

func setRoutes(s *Server) {
	r := s.router

	r.HandleFunc("/health", s.handleHealthCheck()).Methods("GET")
}

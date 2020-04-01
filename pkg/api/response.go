package api

import (
	"encoding/json"
	"net/http"
)

// Response is an API Response
type Response struct {
	Status         int
	Body           interface{}
	responseWriter http.ResponseWriter
	useWrapper     bool
}

// NewResponse creates a new API Response
func NewResponse(w http.ResponseWriter) *Response {
	w.Header().Set("Content-Type", "application/json")
	return &Response{
		responseWriter: w,
		useWrapper:     true,
		Status:         http.StatusOK,
	}
}

// Send the built response
func (r *Response) Send() {
	r.writeBody()
}

// writeBody send the JSON response
func (r *Response) writeBody() {
	w := r.responseWriter
	var body interface{}

	if r.useWrapper == true {
		body = map[string]interface{}{
			"status": r.Status,
			"data":   r.Body,
		}
	} else {
		body = r.Body
	}
	json.NewEncoder(w).Encode(body)
}

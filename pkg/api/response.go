package api

import (
	"encoding/json"
	"net/http"
)

// Response is an API Response
type Response struct {
	Status         int
	Message        string
	Body           interface{}
	responseWriter http.ResponseWriter
	useWrapper     bool
}

// SendResponse creates and immediately sends a response
func SendResponse(w http.ResponseWriter, status int, body interface{}, message string) {
	r := NewResponse(w)
	r.Status = status
	r.Body = body
	r.Message = message
	r.Send()
}

// NewResponse creates a new API Response
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{
		responseWriter: w,
		useWrapper:     true,
		Status:         http.StatusOK,
	}
}

// Send the built response
func (r *Response) Send() {
	r.responseWriter.Header().Set("Content-Type", "application/json")
	r.responseWriter.WriteHeader(r.Status)
	r.writeBody()
}

// writeBody send the JSON response
func (r *Response) writeBody() {
	w := r.responseWriter
	var body interface{}

	if r.useWrapper == true {
		body = map[string]interface{}{
			"status":  r.Status,
			"message": r.Message,
			"data":    r.Body,
		}
	} else {
		body = r.Body
	}
	json.NewEncoder(w).Encode(body)
}

package api

import "net/http"

// CTX represents the request context
type CTX struct {
	W    http.ResponseWriter
	R    *http.Request
	user interface{}
}

// HandlerFunc represents and API request handler
type HandlerFunc = func(ctx *CTX)

// Handler wraps handlers and creates a context
func Handler(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &CTX{
			w,
			r,
			nil,
		}
		next(ctx)
	}
}

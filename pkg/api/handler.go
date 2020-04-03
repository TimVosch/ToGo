package api

import (
	"net/http"
)

// CTX represents the request context
type CTX struct {
	W    http.ResponseWriter
	R    *http.Request
	User interface{}
}

// HandlerFunc represents and API request handler
type HandlerFunc = func(ctx *CTX, next func())

// Handler wraps handlers and creates a context
func Handler(handlers ...HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create the context for this request
		ctx := &CTX{
			w,
			r,
			nil,
		}

		/*
			Create the `next` closure. Everytime next is called
			the counter increases and the next handler from the
			`handlers` array will be used. If i > len(handlers)
			nothing happens
		*/
		i := 0
		var next func()
		next = func() {
			i++
			if i < len(handlers) {
				handlers[i](ctx, next)
			}
		}

		// Call the first handler
		handlers[0](ctx, next)
	}
}

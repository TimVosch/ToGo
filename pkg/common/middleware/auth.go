package middleware

import (
	"net/http"
	"strings"

	"github.com/timvosch/togo/pkg/api"
	"github.com/timvosch/togo/pkg/jwt"
)

type wrappedHandlerFunc = func() api.HandlerFunc

// MakeAuth returns an authentication middleware
func MakeAuth(v *jwt.Verifier) wrappedHandlerFunc {
	return func() api.HandlerFunc {
		return func(ctx *api.CTX, next func()) {
			header := ctx.R.Header.Get("Authorization")
			parts := strings.Split(header, " ")
			if len(parts) != 2 {
				ctx.SendResponse(http.StatusUnauthorized, nil, "Must be authenticated")
				return
			}
			if parts[0] != "Bearer" {
				ctx.SendResponse(http.StatusUnauthorized, nil, "Authorization method not supported")
				return
			}

			token, err := v.Verify(parts[1])
			if err != nil {
				ctx.SendResponse(http.StatusUnauthorized, nil, "Provided JWT is invalid")
				return
			}

			ctx.User = token.Body
			next()
		}
	}
}

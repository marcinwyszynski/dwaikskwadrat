package pkg

import (
	"context"
	"crypto/subtle"
	"errors"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

type contextKey string

const contextKeyToken contextKey = "token"

// ErrUnauthorized is returned when the request is not authorized.
var ErrUnauthorized = errors.New("unauthorized")

// AddBearerTokenFromHTTP gets a bearer token from the HTTP request and puts it
// in the returned context.
func AddBearerTokenFromHTTP(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(
		ctx,
		contextKeyToken,
		strings.TrimPrefix(
			r.Header.Get("Authorization"),
			"Bearer ",
		),
	)
}

// Authorize ensures that only authorized requests get a response.
func Authorize(token string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r interface{}) (interface{}, error) {
			requestToken, ok := ctx.Value(contextKeyToken).(string)
			if !ok {
				requestToken = ""
			}

			if subtle.ConstantTimeCompare([]byte(token), []byte(requestToken)) == 1 {
				return next(ctx, r)
			}

			return nil, ErrUnauthorized
		}
	}
}

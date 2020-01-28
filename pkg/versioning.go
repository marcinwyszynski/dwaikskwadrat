package pkg

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

const versionKeyToken contextKey = "version"

// AddVersionFromHTTP gets API version info from the HTTP request and puts it
// in the returned context.
func AddVersionFromHTTP(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(
		ctx,
		versionKeyToken,
		r.Header.Get("Accept-version"),
	)
}

// Versioning directs the request to various endpoints based on version
// information in the context.
func Versioning(fallback endpoint.Endpoint, versions map[string]endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		if next, ok := versions[ctx.Value(versionKeyToken).(string)]; ok {
			return next(ctx, r)
		}

		return fallback(ctx, r)
	}
}

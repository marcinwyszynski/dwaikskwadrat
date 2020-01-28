package pkg

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// LoggingMiddleware logs request duration and the presence of bearer token.
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r interface{}) (interface{}, error) {
			logger = log.With(logger, "has_token", ctx.Value(contextKeyToken) != "")

			defer func(begin time.Time) {
				level.Info(logger).Log("duration_micro", time.Since(begin).Microseconds())
			}(time.Now())

			return next(ctx, r)
		}
	}
}

package pkg

import (
	"context"
	"runtime"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
)

// MetricsMiddleware reports the number of goroutines and request duration.
func MetricsMiddleware(goroutines metrics.Gauge, requests metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r interface{}) (interface{}, error) {
			goroutines.Set(float64(runtime.NumGoroutine()))

			defer func(begin time.Time) {
				requests.Observe(float64(time.Since(begin).Microseconds()))
			}(time.Now())

			return next(ctx, r)
		}
	}
}

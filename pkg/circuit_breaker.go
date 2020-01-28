package pkg

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
)

// CircuitBreaker is a circuit breaker with a fallback value.
func CircuitBreaker(commandName string, fallback interface{}) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r interface{}) (interface{}, error) {
			var resp interface{}

			return resp, hystrix.Do(commandName, func() (err error) {
				resp, err = next(ctx, r)
				return err
			}, func(_ error) error {
				resp = fallback
				return nil
			})
		}
	}
}

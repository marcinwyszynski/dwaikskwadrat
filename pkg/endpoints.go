package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	zipkinkit "github.com/go-kit/kit/tracing/zipkin"
	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/openzipkin/zipkin-go"
)

// MakeIntegerClientEndpoint creates a client endpoint.
func MakeIntegerClientEndpoint(tracer *zipkin.Tracer, host, path string) endpoint.Endpoint {
	return httpkit.NewClient(
		http.MethodPost,
		&url.URL{Scheme: "http", Host: host, Path: path},
		httpkit.EncodeJSONRequest,
		DecodeIntegerResponse,
		httpkit.ClientBefore(func(ctx context.Context, r *http.Request) context.Context {
			r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctx.Value(contextKeyToken)))
			return ctx
		}),
		zipkinkit.HTTPClientTrace(tracer, zipkinkit.Tags(map[string]string{"role": "client"})),
	).Endpoint()
}

// MakeDoublerServerEndpoint creates a server endpoint wrapping a Doubler.
func MakeDoublerServerEndpoint(doubler Doubler) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		var ret IntegerResponse
		var err error

		ret.Body.Output, err = doubler.Double(req.(IntegerRequest).Input)

		return ret, err
	}
}

// MakeDoublerServerEndpointV2 creates a server endpoint wrapping a Doubler.
func MakeDoublerServerEndpointV2(doubler Doubler) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		var ret IntegerResponse
		var err error

		ret.Body.Output, err = doubler.Double(req.(IntegerRequest).Input)
		ret.Body.Output++

		return ret, err
	}
}

// MakeSquarerServerEndpoint creates a server endpoint wrapping a Squarer.
func MakeSquarerServerEndpoint(squarer Squarer) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		var ret IntegerResponse
		var err error

		ret.Body.Output, err = squarer.Square(req.(IntegerRequest).Input)

		return ret, err
	}
}

// MakeDoubleSquarerServerEndpoint creates a server endpoint wrapping double and
// square client endpoints.
func MakeDoubleSquarerServerEndpoint(double, square endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ret, err := square(ctx, req)
		if err != nil {
			return ret, err
		}

		out, ok := ret.(IntegerResponse)
		if !ok {
			return nil, errors.New("unexpected double response format")
		}

		return double(ctx, IntegerRequest{Input: out.Body.Output})
	}
}

package pkg

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

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

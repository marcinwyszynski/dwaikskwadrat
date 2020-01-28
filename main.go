package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"

	"github.com/marcinwyszynski/dwaikskwadrat/pkg"
)

const (
	binding = "localhost:8080"
	token   = "bacon"
)

func main() {
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.Timestamp(time.Now))

	mathService := new(pkg.MathService)
	router := httprouter.New()

	middlewareChain := endpoint.Chain(
		pkg.LoggingMiddleware(logger),
		pkg.Authorize(token),
	)

	defaultDoubler := pkg.MakeDoublerServerEndpoint(mathService)

	router.Handler(http.MethodPost, "/double", httpkit.NewServer(
		middlewareChain(pkg.Versioning(defaultDoubler, map[string]endpoint.Endpoint{
			"v1": defaultDoubler,
			"v2": pkg.MakeDoublerServerEndpointV2(mathService),
		})),
		pkg.DecodeIntegerRequest,
		httpkit.EncodeJSONResponse,
		httpkit.ServerBefore(
			pkg.AddBearerTokenFromHTTP,
			pkg.AddVersionFromHTTP,
		),
	))

	router.Handler(http.MethodPost, "/square", httpkit.NewServer(
		middlewareChain(pkg.MakeSquarerServerEndpoint(mathService)),
		pkg.DecodeIntegerRequest,
		httpkit.EncodeJSONResponse,
		httpkit.ServerBefore(pkg.AddBearerTokenFromHTTP),
	))

	level.Info(logger).Log("msg", "starting server", "binding", binding)
	level.Error(logger).Log("problem", http.ListenAndServe(binding, router))
}

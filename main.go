package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	influxkit "github.com/go-kit/kit/metrics/influx"
	zipkinkit "github.com/go-kit/kit/tracing/zipkin"
	httpkit "github.com/go-kit/kit/transport/http"
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"

	"github.com/marcinwyszynski/dwaikskwadrat/pkg"
)

const (
	binding = "localhost:8080"
	token   = "bacon"
)

func main() {
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.Timestamp(time.Now))

	// Zipkin
	reporter := zipkinhttp.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()
	zipkinEndpoint, err := zipkin.NewEndpoint("dwaikskwadrat", binding)
	if err != nil {
		level.Error(logger).Log("error", "failed to create a Zipkin endpoint", "cause", err)
		os.Exit(1)
	}
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zipkinEndpoint))
	if err != nil {
		level.Error(logger).Log("error", "failed to create a Zipkin tracer", "cause", err)
		os.Exit(1)
	}

	zipkinServerOption := zipkinkit.HTTPServerTrace(
		tracer,
		zipkinkit.Logger(logger),
		zipkinkit.Tags(map[string]string{"role": "server"}),
	)

	influxClient, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{Addr: "http://localhost:8086"})
	if err != nil {
		level.Error(logger).Log("error", "failed to create InfluxDB client", "cause", err)
		os.Exit(1)
	}
	defer influxClient.Close()

	influxFactory := influxkit.New(
		map[string]string{"service": "dwaikskwadrat"},
		influxdb.BatchPointsConfig{Database: "metrics"},
		logger,
	)

	go func() {
		influxFactory.WriteLoop(context.Background(), time.Tick(10*time.Second), influxClient)
	}()

	mathService := &pkg.MathService{
		Operations: influxFactory.NewCounter("operations"),
	}
	router := httprouter.New()

	var defaultResponse pkg.IntegerResponse
	defaultResponse.Body.Output = 42

	middlewareChain := endpoint.Chain(
		pkg.MetricsMiddleware(influxFactory.NewGauge("goroutines"), influxFactory.NewHistogram("requests")),
		pkg.LoggingMiddleware(logger),
		zipkinkit.TraceEndpoint(tracer, "auth"),
		pkg.Authorize(token),
		zipkinkit.TraceEndpoint(tracer, "hystrix"),
		pkg.CircuitBreaker("hystrix", defaultResponse),
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
		zipkinServerOption,
	))

	router.Handler(http.MethodPost, "/square", httpkit.NewServer(
		middlewareChain(pkg.MakeSquarerServerEndpoint(mathService)),
		pkg.DecodeIntegerRequest,
		httpkit.EncodeJSONResponse,
		httpkit.ServerBefore(pkg.AddBearerTokenFromHTTP),
		zipkinServerOption,
	))

	router.Handler(http.MethodPost, "/doublesquare", httpkit.NewServer(
		middlewareChain(pkg.MakeDoubleSquarerServerEndpoint(
			pkg.MakeIntegerClientEndpoint(tracer, binding, "/double"),
			pkg.MakeIntegerClientEndpoint(tracer, binding, "/square"),
		)),
		pkg.DecodeIntegerRequest,
		httpkit.EncodeJSONResponse,
		httpkit.ServerBefore(pkg.AddBearerTokenFromHTTP),
		zipkinServerOption,
	))

	level.Info(logger).Log("msg", "starting server", "binding", binding)
	level.Error(logger).Log("problem", http.ListenAndServe(binding, router))
}

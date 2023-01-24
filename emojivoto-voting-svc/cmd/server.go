package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/buoyantio/emojivoto/emojivoto-voting-svc/api"
	"github.com/buoyantio/emojivoto/emojivoto-voting-svc/voting"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	grpcPort                   = os.Getenv("GRPC_PORT")
	promPort                   = os.Getenv("PROM_PORT")
	ocagentHost                = os.Getenv("OC_AGENT_HOST")
	failureRateVar             = os.Getenv("FAILURE_RATE")
	failureRateFloat           = float64(0.0)
	artificialDelayVar         = os.Getenv("ARTIFICIAL_DELAY")
	artificialDelayDuration, _ = time.ParseDuration("0ms")
)

func main() {

	if grpcPort == "" {
		log.Fatalf("GRPC_PORT (currently [%s]) environment variable must me set to run the server.", grpcPort)
	}

	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("votoemoji-voting"),
	)

	client := otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint(ocagentHost), otlptracegrpc.WithInsecure())
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		panic(fmt.Sprintf("creating OTLP trace exporter: %v", err))
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(r),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)

	poll := voting.NewPoll()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		panic(err)
	}

	errs := make(chan error, 1)

	if promPort != "" {
		// Start prometheus server
		go func() {
			log.Printf("Starting prom metrics on PROM_PORT=[%s]", promPort)
			http.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(fmt.Sprintf(":%s", promPort), nil)
			errs <- err
		}()
	}

	// Start grpc server
	go func() {
		grpc_prometheus.EnableHandlingTimeHistogram()
		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
			grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		)

		setFailureRateOrDefault(failureRateVar, &failureRateFloat)

		setArtificialDelayOrDefault(artificialDelayVar, &artificialDelayDuration)

		api.NewGrpServer(grpcServer, poll, float32(failureRateFloat), artificialDelayDuration)
		grpc_prometheus.Register(grpcServer)
		log.Printf("Starting grpc server on GRPC_PORT=[%s]", grpcPort)
		log.Printf("Using failureRate [%f] and artificialDelayDuration [%v]", failureRateFloat, artificialDelayDuration)
		err := grpcServer.Serve(lis)
		errs <- err
	}()

	// Catch shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT)
		s := <-sig
		errs <- fmt.Errorf("caught signal %v", s)
	}()

	log.Fatal(<-errs)
}

func setFailureRateOrDefault(failureRateVar string, failureRateFloat *float64) {
	if failureRateVar != "" {
		var err error
		*failureRateFloat, err = strconv.ParseFloat(failureRateVar, 32)
		if err != nil {
			log.Printf("Invalid value for FAILURE_RATE %v. Using %f instead", failureRateVar, *failureRateFloat)
		}

		if *failureRateFloat > 1.0 {
			log.Printf("FAILURE_RATE is set to %f. It should be a value between 0.0 and 1.0", *failureRateFloat)
		}
	}
}

func setArtificialDelayOrDefault(artificialDelayVar string, artificialDelayDuration *time.Duration) {
	if artificialDelayVar != "" {
		var err error
		*artificialDelayDuration, err = time.ParseDuration(artificialDelayVar)
		if err != nil {
			log.Printf("ParseDuration failed for %v using %v instead", artificialDelayVar, *artificialDelayDuration)
		}
	}
}

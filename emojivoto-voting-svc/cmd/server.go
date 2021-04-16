package main

import (
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

	"contrib.go.opencensus.io/exporter/ocagent"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
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

	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress(ocagentHost),
		ocagent.WithServiceName("voting"))
	if err != nil {
		log.Fatalf("Failed to create ocagent-exporter: %v", err)
	}
	trace.RegisterExporter(oce)

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
			grpc.StatsHandler(&ocgrpc.ServerHandler{}),
			grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
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

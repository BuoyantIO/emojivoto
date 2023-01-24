package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/buoyantio/emojivoto/emojivoto-web/gen/proto"
	"github.com/buoyantio/emojivoto/emojivoto-web/web"
	"google.golang.org/grpc"

	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	webPort              = os.Getenv("WEB_PORT")
	emojisvcHost         = os.Getenv("EMOJISVC_HOST")
	votingsvcHost        = os.Getenv("VOTINGSVC_HOST")
	indexBundle          = os.Getenv("INDEX_BUNDLE")
	webpackDevServerHost = os.Getenv("WEBPACK_DEV_SERVER")
	ocagentHost          = os.Getenv("OC_AGENT_HOST")
)

func main() {

	if webPort == "" || emojisvcHost == "" || votingsvcHost == "" {
		log.Fatalf("WEB_PORT (currently [%s]) EMOJISVC_HOST (currently [%s]) and VOTINGSVC_HOST (currently [%s]) INDEX_BUNDLE (currently [%s]) environment variables must me set.", webPort, emojisvcHost, votingsvcHost, indexBundle)
	}

	// Identify service
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("votomoji-web"),
	)

	// Create exporter client
	client := otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint(ocagentHost), otlptracegrpc.WithInsecure())
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		panic(fmt.Sprintf("creating OTLP trace exporter: %v", err))
	}

	// Create Tracer Provider
	// TracerProvider will provide instrumentations with an impl of Tracer.
	// Tracer funnels data to export pipelines (span processors)
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

	votingSvcConn := openGrpcClientConnection(votingsvcHost)
	votingClient := pb.NewVotingServiceClient(votingSvcConn)
	defer votingSvcConn.Close()

	emojiSvcConn := openGrpcClientConnection(emojisvcHost)
	emojiSvcClient := pb.NewEmojiServiceClient(emojiSvcConn)
	defer emojiSvcConn.Close()

	web.StartServer(webPort, webpackDevServerHost, indexBundle, emojiSvcClient, votingClient)
}

func openGrpcClientConnection(host string) *grpc.ClientConn {
	log.Printf("Connecting to [%s]", host)
	conn, err := grpc.Dial(
		host,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))

	if err != nil {
		panic(err)
	}
	return conn
}

package main

import (
	"log"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	pb "github.com/buoyantio/emojivoto/emojivoto-web/gen/proto"
	"github.com/buoyantio/emojivoto/emojivoto-web/web"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
	tracer.Start(
		tracer.WithEnv("staging"),
		tracer.WithService("webserver"),
		tracer.WithServiceVersion("env"),
	)
	if webPort == "" || emojisvcHost == "" || votingsvcHost == "" {
		log.Fatalf("WEB_PORT (currently [%s]) EMOJISVC_HOST (currently [%s]) and VOTINGSVC_HOST (currently [%s]) INDEX_BUNDLE (currently [%s]) environment variables must me set.", webPort, emojisvcHost, votingsvcHost, indexBundle)
	}

	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress(ocagentHost),
		ocagent.WithServiceName("web"))
	if err != nil {
		log.Fatalf("Failed to create ocagent-exporter: %v", err)
	}
	trace.RegisterExporter(oce)

	votingSvcConn := openGrpcClientConnection(votingsvcHost)
	votingClient := pb.NewVotingServiceClient(votingSvcConn)
	defer votingSvcConn.Close()

	emojiSvcConn := openGrpcClientConnection(emojisvcHost)
	emojiSvcClient := pb.NewEmojiServiceClient(emojiSvcConn)
	defer emojiSvcConn.Close()

	web.StartServer(webPort, webpackDevServerHost, indexBundle, emojiSvcClient, votingClient)
	defer tracer.Stop()
}

func openGrpcClientConnection(host string) *grpc.ClientConn {
	log.Printf("Connecting to [%s]", host)
	conn, err := grpc.Dial(
		host,
		grpc.WithInsecure(),
		grpc.WithStatsHandler(new(ocgrpc.ClientHandler)))

	if err != nil {
		panic(err)
	}
	return conn
}

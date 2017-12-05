package main

import (
	"log"
	"os"

	pb "github.com/runconduit/conduit-examples/emojivoto/emojivoto-web/gen/proto"
	"github.com/runconduit/conduit-examples/emojivoto/emojivoto-web/web"
	"google.golang.org/grpc"
)

var (
	webPort       = os.Getenv("WEB_PORT")
	emojisvcHost  = os.Getenv("EMOJISVC_HOST")
	votingsvcHost = os.Getenv("VOTINGSVC_HOST")
	indexBundle   = os.Getenv("INDEX_BUNDLE")
)

func main() {

	if webPort == "" || emojisvcHost == "" || votingsvcHost == "" {
		log.Fatalf("WEB_PORT (currently [%s]) EMOJISVC_HOST (currently [%s]) and VOTINGSVC_HOST (currently [%s]) INDEX_BUNDLE (currently [%s]) environment variables must me set.", webPort, emojisvcHost, votingsvcHost, indexBundle)
	}

	votingSvcConn := openGrpcClientConnection(votingsvcHost)
	votingClient := pb.NewVotingServiceClient(votingSvcConn)
	defer votingSvcConn.Close()

	emojiSvcConn := openGrpcClientConnection(emojisvcHost)
	emojiSvcClient := pb.NewEmojiServiceClient(emojiSvcConn)
	defer emojiSvcConn.Close()

	web.StartServer(webPort, indexBundle, emojiSvcClient, votingClient)
}

func openGrpcClientConnection(host string) *grpc.ClientConn {
	log.Printf("Connecting to [%s]", host)
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}

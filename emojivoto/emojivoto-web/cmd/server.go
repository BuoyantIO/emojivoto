package main

import (
	"os"
	"log"
	"google.golang.org/grpc"
	pb "github.com/buoyantio/conduit-examples/emojivoto/emojivoto-web/gen/proto"
	"github.com/buoyantio/conduit-examples/emojivoto/emojivoto-web/web"
)

var (
	webPort       = os.Getenv("WEB_PORT")
	emojisvcHost  = os.Getenv("EMOJISVC_HOST")
	votingsvcHost = os.Getenv("VOTINGSVC_HOST")
)

func main() {

	if webPort == "" || emojisvcHost == "" || votingsvcHost == "" {
		log.Fatalf("WEB_PORT (currently [%s]) EMOJISVC_HOST (currently [%s]) and VOTINGSVC_HOST (currently [%s]) environment variables must me set.", webPort, emojisvcHost, votingsvcHost)
	}

	votingSvcConn := openGrpcClientConnection(votingsvcHost)
	votingClient := pb.NewVotingServiceClient(votingSvcConn)
	defer votingSvcConn.Close()

	emojiSvcConn := openGrpcClientConnection(emojisvcHost)
	emojiSvcClient := pb.NewEmojiServiceClient(emojiSvcConn)
	defer emojiSvcConn.Close()

	web.StartServer(webPort, emojiSvcClient, votingClient)
}

func openGrpcClientConnection(host string) *grpc.ClientConn {
	log.Printf("Connecting to [%s]", host)
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}

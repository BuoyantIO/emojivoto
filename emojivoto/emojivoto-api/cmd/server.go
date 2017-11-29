package main

import (
	"os"
	"log"
	"google.golang.org/grpc"
	pb "github.com/buoyantio/conduit-examples/emojivoto/emojivoto-api/gen/proto"
	"github.com/buoyantio/conduit-examples/emojivoto/emojivoto-api/api"
)

var (
	apiPort       = os.Getenv("API_PORT")
	emojisvcHost = os.Getenv("EMOJISVC_HOST")
)

func main() {

	if apiPort == "" || emojisvcHost == "" {
		log.Fatalf("WEB_PORT (currently [%s]) and EMOJISVC_HOST (currently [%s]) environment variables must me set.", apiPort, emojisvcHost)
	}

	conn := openGrpcClientConnection(emojisvcHost)

	emojivotoClient := pb.NewEmojiServiceClient(conn)
	api.StartServer(apiPort, emojivotoClient)
	defer conn.Close()
}

func openGrpcClientConnection(host string) *grpc.ClientConn {
	log.Printf("Connecting to [%s]", host)
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return conn
}

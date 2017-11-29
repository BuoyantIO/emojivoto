package main

import (
	"os"
	"log"
	"github.com/buoyantio/conduit-examples/emojivoto/emojivoto-voting-svc/voting"
	"fmt"
	"net"
	"google.golang.org/grpc"
	"github.com/buoyantio/conduit-examples/emojivoto/emojivoto-voting-svc/api"
)

var (
	grpcPort = os.Getenv("GRPC_PORT")
)

func main() {

	if grpcPort == "" {
		log.Fatalf("GRPC_PORT (currently [%s]) environment variable must me set to run the server.", grpcPort)
	}

	poll := voting.NewPoll()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	api.NewGrpServer(grpcServer, poll)
	log.Printf("Starting grpc server on GRPC_PORT=[%s]", grpcPort)
	grpcServer.Serve(lis)
}

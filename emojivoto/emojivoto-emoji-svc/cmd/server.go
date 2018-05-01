package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/runconduit/conduit-examples/emojivoto/emojivoto-emoji-svc/api"
	"github.com/runconduit/conduit-examples/emojivoto/emojivoto-emoji-svc/emoji"
	"google.golang.org/grpc"
)

var (
	grpcPort  = os.Getenv("GRPC_PORT")
	HTTP1Addr = os.Getenv("HTTP1_ADDR")
)

func main() {

	if grpcPort == "" {
		log.Fatalf("GRPC_PORT (currently [%s]) environment variable must me set to run the server.", grpcPort)
	}
	if HTTP1Addr == "" {
		log.Fatalf("HTTP1_ADDR (currently [%s]) environment variable must me set to run the server.", HTTP1Addr)
	}

	allEmoji := emoji.NewAllEmoji()

	go func() {
		log.Printf("Starting HTTP1 server on HTTP1_ADDR=[%s]", HTTP1Addr)
		api.NewHTTP1Server(HTTP1Addr, allEmoji)
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	api.NewGrpServer(grpcServer, allEmoji)
	log.Printf("Starting grpc server on GRPC_PORT=[%s]", grpcPort)
	grpcServer.Serve(lis)
}

package api

import (
	"context"
	"fmt"

	pb "github.com/buoyantio/conduit-examples/emojivoto/emojivoto-voting-svc/gen/proto"
	"github.com/buoyantio/conduit-examples/emojivoto/emojivoto-voting-svc/voting"
	"google.golang.org/grpc"
)

type PollServiceServer struct {
	poll voting.Poll
}

func (pS *PollServiceServer) Vote(context context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	sc := req.Shortcode
	if sc == ":shit:" || sc == ":poop:" || sc == ":hankey:" {
		return nil, fmt.Errorf("ERROR")
	}
	err := pS.poll.Vote(sc)
	return &pb.VoteResponse{}, err
}

func (pS *PollServiceServer) Results(context.Context, *pb.ResultsRequest) (*pb.ResultsResponse, error) {
	results, e := pS.poll.Results()
	if e != nil {
		return nil, e
	}

	votingResults := make([]*pb.VotingResult, 0)
	for _, e := range results {
		result := pb.VotingResult{
			Shortcode: e.Shortcode,
			Votes:     int32(e.NumVotes),
		}
		votingResults = append(votingResults, &result)
	}

	response := &pb.ResultsResponse{
		Results: votingResults,
	}
	return response, nil
}

func NewGrpServer(grpcServer *grpc.Server, poll voting.Poll) {
	server := &PollServiceServer{
		poll,
	}

	pb.RegisterVotingServiceServer(grpcServer, server)
}

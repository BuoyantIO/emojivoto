package api

import (
	"context"
	"google.golang.org/grpc"
	pb "github.com/buoyantio/conduit-examples/emojivoto/emojivoto-voting-svc/gen/proto"
	"github.com/buoyantio/conduit-examples/emojivoto/emojivoto-voting-svc/voting"
)

type PollServiceServer struct {
	poll voting.Poll
}

func (pS *PollServiceServer) Vote(context context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	err := pS.poll.Vote(req.Shortcode)
	return nil, err
}

func (pS *PollServiceServer) Results(context.Context, *pb.ResultsRequest) (*pb.ResultsResponse, error) {
	results, e := pS.poll.Results()
	if e != nil {
		return &pb.ResultsResponse{}, e
	}

	votingResults := make([]*pb.VotingResult, 0)
	for _, e := range results {
		result := pb.VotingResult{
			Shortcode: e.Shortcode,
			Votes: int32(e.NumVotes),
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

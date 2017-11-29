package api

import (
	"testing"
	"context"
	pb "github.com/buoyantio/boron/demos/emojivoto/emojivoto-voting-svc/gen/proto"
	"github.com/buoyantio/boron/demos/emojivoto/emojivoto-voting-svc/voting"
)

func TestVote(t *testing.T) {
	t.Run("Computes vote", func(t *testing.T) {
		ctx := context.Background()
		poll := voting.NewPoll()
		emojivotoService := PollServiceServer{
			poll: poll,
		}

		shortcodeVotedFor := ":wave:"

		request := pb.VoteRequest{
			Shortcode: shortcodeVotedFor,
		}
		_, err := emojivotoService.Vote(ctx, &request)

		if err != nil {
			t.Fatal(err)
		}

		if r, _ := poll.Results(); len(r) == 0 || r[0].Shortcode != shortcodeVotedFor {
			t.Fatalf("Voted for [%s] but results were [%v]", shortcodeVotedFor, r)
		}
	})
}

func TestLeaderboard(t *testing.T) {
	t.Run("Returns expected leaderboard", func(t *testing.T) {
		ctx := context.Background()
		poll := voting.NewPoll()
		emojivotoService := PollServiceServer{
			poll: poll,
		}

		votedForTwice := ":wave:"
		votedForOnce := ":poop:"

		emojivotoService.Vote(ctx, &pb.VoteRequest{Shortcode: votedForTwice})
		emojivotoService.Vote(ctx, &pb.VoteRequest{Shortcode: votedForTwice})
		emojivotoService.Vote(ctx, &pb.VoteRequest{Shortcode: votedForOnce})

		response, err := emojivotoService.Results(context.Background(), &pb.ResultsRequest{})

		if err != nil {
			t.Fatal(err)
		}

		if len(response.Results) != 2 {
			t.Fatalf("Expected results to contain two emoji, found: [%v]", response.Results)
		}

		if response.Results[0].Shortcode != votedForTwice|| response.Results[0].Votes != 2 {
			t.Fatalf("Expected results to be [%s,%s], found: [%v]", votedForTwice, 2, response.Results)
		}

		if response.Results[1].Shortcode != votedForOnce || response.Results[1].Votes != 1 {
			t.Fatalf("Expected results to be [%s,%s], found: [%v]", votedForOnce, 1, response.Results)
		}
	})
}

//TODO: test for errors

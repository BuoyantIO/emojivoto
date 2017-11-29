package web

import (
	"net/http"
	"testing"

	"net/http/httptest"
	"encoding/json"
	"context"
	"google.golang.org/grpc"
	pb "github.com/buoyantio/boron/demos/emojivoto/emojivoto-web/gen/proto"
	"net/url"
)

type MockEmojiServiceClient struct {
	emojiList []*pb.Emoji
}

func (c *MockEmojiServiceClient) ListAll(ctx context.Context, in *pb.ListAllEmojiRequest, opts ...grpc.CallOption) (*pb.ListAllEmojiResponse, error) {
	response := pb.ListAllEmojiResponse{
		List: c.emojiList,
	}

	return &response, nil
}

func (c *MockEmojiServiceClient) FindByShortcode(ctx context.Context, req *pb.FindByShortcodeRequest, opts ...grpc.CallOption) (*pb.FindByShortcodeResponse, error) {
	foundEmoji := c.findByShortcode(req.Shortcode)

	return &pb.FindByShortcodeResponse{
		Emoji: foundEmoji,
	}, nil
}

func (c *MockEmojiServiceClient) findByShortcode(shortcode string) *pb.Emoji {
	var foundEmoji *pb.Emoji
	for _, e := range c.emojiList {
		if e.Shortcode == shortcode {
			foundEmoji = &pb.Emoji{
				Shortcode: e.Shortcode,
				Unicode:   e.Unicode,
			}
		}
	}
	return foundEmoji
}

type MockVotingServiceClient struct {
	lastChoiceShortcode string
	resultToReturn      []*pb.VotingResult
}

func (c *MockVotingServiceClient) Vote(ctx context.Context, req *pb.VoteRequest, opts ...grpc.CallOption) (*pb.VoteResponse, error) {
	c.lastChoiceShortcode = req.Shortcode
	return &pb.VoteResponse{}, nil
}

func (c *MockVotingServiceClient) Results(ctx context.Context, in *pb.ResultsRequest, opts ...grpc.CallOption) (*pb.ResultsResponse, error) {
	return &pb.ResultsResponse{
		Results: c.resultToReturn,
	}, nil
}

func TestListEmojiHandler(t *testing.T) {

	t.Run("returns correct list", func(t *testing.T) {
		expectedList := []*pb.Emoji{{Shortcode: "a", Unicode: "\\a"}, {Shortcode: "b", Unicode: "\\b"}}
		emojiSvcClient := &MockEmojiServiceClient{emojiList: expectedList}
		webApp := &WebApp{
			emojiServiceClient: emojiSvcClient,
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(webApp.listEmojiHandler)

		req, err := http.NewRequest("GET", "/list", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var responseList []map[string]string

		if err := json.Unmarshal(rr.Body.Bytes(), &responseList); err != nil {
			t.Fatalf("parsing response returned error [%v] response:\n%s", err, string(rr.Body.Bytes()))
		}

		if len(responseList) != len(expectedList) {
			t.Fatalf("Expected response to contain [%d] emoji, got [%d]", len(expectedList), len(responseList))
		}

		responseMap := make(map[string]map[string]string)
		for _, obj := range responseList {
			responseMap[obj["unicode"]] = obj
		}

		for _, e := range expectedList {
			responseEmoji := responseMap[e.Unicode]
			if responseEmoji["shortcode"] != e.Shortcode || responseEmoji["unicode"] != e.Unicode {
				t.Fatalf("Response didnt contain [%v]", e)
			}
		}
	})
}

func TestVoteHandler(t *testing.T) {
	t.Run("registers the vote if everything is valid", func(t *testing.T) {
		emojiIWantToVoteFor := &pb.Emoji{Shortcode: "voted", Unicode: "\\voted"}
		expectedList := []*pb.Emoji{emojiIWantToVoteFor}
		emojiSvcClient := &MockEmojiServiceClient{emojiList: expectedList}
		votingServiceClient := &MockVotingServiceClient{}
		webApp := &WebApp{
			emojiServiceClient:  emojiSvcClient,
			votingServiceClient: votingServiceClient,
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(webApp.voteEmojiHandler)

		req, err := http.NewRequest("POST", "/voting", nil)
		if err != nil {
			t.Fatal(err)
		}

		form := url.Values{}
		form.Add("choice", emojiIWantToVoteFor.Shortcode)
		req.PostForm = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		if votingServiceClient.lastChoiceShortcode != emojiIWantToVoteFor.Shortcode {
			t.Fatalf("Expected vote to be for [%v], but got [%v]", emojiIWantToVoteFor, votingServiceClient.lastChoiceShortcode)
		}
	})

	t.Run("rejects request if doesnt contain choice parameter", func(t *testing.T) {
		webApp := &WebApp{
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(webApp.voteEmojiHandler)

		req, err := http.NewRequest("POST", "/voting", nil)
		if err != nil {
			t.Fatal(err)
		}

		form := url.Values{}

		req.PostForm = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code when missing voter: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("rejects request if emoji doesnt exist", func(t *testing.T) {
		emojiIWantToVoteFor := &pb.Emoji{Shortcode: "voted", Unicode: "\\voted"}
		expectedList := []*pb.Emoji{emojiIWantToVoteFor}
		emojiSvcClient := &MockEmojiServiceClient{emojiList: expectedList}
		votingServiceClient := &MockVotingServiceClient{}
		webApp := &WebApp{
			emojiServiceClient:  emojiSvcClient,
			votingServiceClient: votingServiceClient,
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(webApp.voteEmojiHandler)

		req, err := http.NewRequest("POST", "/api/vote", nil)
		if err != nil {
			t.Fatal(err)
		}

		form := url.Values{}
		form.Add("choice", "thiz doesnt exist")

		req.PostForm = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Fatalf("handler returned wrong status code when missing voter: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
}

func TestLeaderboard(t *testing.T) {

	t.Run("registers the vote if everything is valid", func(t *testing.T) {
		expectedList := []*pb.Emoji{
			{Shortcode: "1voted", Unicode: "\\1voted"},
			{Shortcode: "2voted", Unicode: "\\2voted"},
			{Shortcode: "3voted", Unicode: "\\3voted"},
			{Shortcode: "4voted", Unicode: "\\4voted"},
		}

		expectedResults := []*pb.VotingResult{
			{
				Votes:     10,
				Shortcode: expectedList[0].Shortcode,
			},
			{
				Votes:     5,
				Shortcode: expectedList[1].Shortcode,
			},
			{
				Votes:     1,
				Shortcode: expectedList[2].Shortcode,
			},
		}
		emojiSvcClient := &MockEmojiServiceClient{emojiList: expectedList}
		votingServiceClient := &MockVotingServiceClient{resultToReturn: expectedResults}
		webApp := &WebApp{
			emojiServiceClient:  emojiSvcClient,
			votingServiceClient: votingServiceClient,
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(webApp.leaderboardHandler)

		req, err := http.NewRequest("GET", "/leaderboard", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var responseList []map[string]string

		if err := json.Unmarshal(rr.Body.Bytes(), &responseList); err != nil {
			t.Fatalf("parsing response returned error %v", err)
		}

		if len(responseList) != len(expectedResults) {
			t.Fatalf("Expected response to contain [%d] items but it has [%d]", len(expectedResults), len(responseList))
		}

		for i := 0; i < len(responseList); i++ {
			actualShortcode := responseList[i]["shortcode"]
			expectedShortcode := expectedResults[i].Shortcode
			if actualShortcode != expectedShortcode {
				t.Fatalf("Expected shortcode for item [%d] to be [%s] but it was [%s]", i, expectedShortcode, actualShortcode)
			}

			actualUnicode := responseList[i]["unicode"]
			expectedUnicode := emojiSvcClient.findByShortcode(expectedShortcode).Unicode
			if actualUnicode != expectedUnicode {
				t.Fatalf("Expected unicode for item [%d] to be [%s] but it was [%s]", i, expectedUnicode, actualUnicode)
			}
		}
	})
}
//TODO: test for errors

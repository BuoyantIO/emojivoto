package api

import (
	"net/http"
	"testing"

	"net/http/httptest"
	"encoding/json"
	"context"
	"google.golang.org/grpc"
	pb "github.com/buoyantio/boron/demos/emojivoto/emojivoto-api/gen/proto"
)

type MockEmojiServiceClient struct {
	EmojiList []*pb.Emoji
}

func (c *MockEmojiServiceClient) ListAll(ctx context.Context, req *pb.ListAllEmojiRequest, opts ...grpc.CallOption) (*pb.ListAllEmojiResponse, error) {
	response := pb.ListAllEmojiResponse{
		List: c.EmojiList,
	}

	return &response, nil
}

func (c *MockEmojiServiceClient) FindByShortcode(ctx context.Context, req *pb.FindByShortcodeRequest, opts ...grpc.CallOption) (*pb.FindByShortcodeResponse, error) {
	var foundEmoji *pb.Emoji

	for _, e := range c.EmojiList {
		if e.Shortcode == req.Shortcode {
			foundEmoji = &pb.Emoji{
				Shortcode: e.Shortcode,
				Unicode:   e.Unicode,
			}
		}
	}

	return &pb.FindByShortcodeResponse{
		Emoji: foundEmoji,
	}, nil
}

func TestListEmojiHandler(t *testing.T) {
	expectedList := []*pb.Emoji{{Shortcode: "a", Unicode: "\\a"}, {Shortcode: "b", Unicode: "\\b"}}
	emojiSvcClient := &MockEmojiServiceClient{EmojiList: expectedList}
	apiApp := &ApiApp{
		emojiClient: emojiSvcClient,
	}

	t.Run("returns correct list", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiApp.listEmojiHandler)

		req, err := http.NewRequest("GET", "/emoji", nil)
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

func TestFindEMojiByShortcodeHandler(t *testing.T) {
	expectedList := []*pb.Emoji{{Shortcode: "a", Unicode: "\\a"}, {Shortcode: "b", Unicode: "\\b"}}
	emojiSvcClient := &MockEmojiServiceClient{EmojiList: expectedList}
	apiApp := &ApiApp{
		emojiClient: emojiSvcClient,
	}

	t.Run("returns correct emoji when found", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiApp.findEmojiHandler)

		req, err := http.NewRequest("GET", "/emoji/a", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var responseEmoji map[string]string

		if err := json.Unmarshal(rr.Body.Bytes(), &responseEmoji); err != nil {
			t.Fatalf("parsing response returned error [%v] response:\n%s", err, string(rr.Body.Bytes()))
		}

		if responseEmoji["shortcode"] != "a" || responseEmoji["unicode"] != "\\a" {
			t.Fatalf("Response didnt contain expected emoji [%v]", responseEmoji)

		}
	})

	t.Run("returns correct status when emoji not found", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiApp.findEmojiHandler)

		req, err := http.NewRequest("GET", "/emoji/does-not-exist", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound{
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})
}

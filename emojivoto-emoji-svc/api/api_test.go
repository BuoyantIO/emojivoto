package api

import (
	"context"
	"testing"

	"github.com/buoyantio/emojivoto/emojivoto-emoji-svc/emoji"
	pb "github.com/buoyantio/emojivoto/emojivoto-emoji-svc/gen/proto"
)

func TestListAll(t *testing.T) {
	t.Run("return all existing emoji", func(t *testing.T) {
		ctx := context.Background()
		allEmoji := emoji.NewAllEmoji()
		emojiService := EmojiServiceServer{
			allEmoji: allEmoji,
		}

		response, err := emojiService.ListAll(ctx, &pb.ListAllEmojiRequest{})

		if err != nil {
			t.Fatal(err)
		}

		if response == nil || len(response.List) == 0 {
			t.Fatal("Returned emoji list was empty")
		}

		responseMap := make(map[string]*pb.Emoji)
		for _, obj := range response.List {
			responseMap[obj.Unicode] = obj
		}

		for _, e := range allEmoji.List() {
			responseEmoji := responseMap[e.Unicode]
			if responseEmoji.Shortcode != e.Shortcode || responseEmoji.Unicode != e.Unicode {
				t.Fatalf("Response didnt contain [%v]", e)
			}
		}
	})
}

func TestFindByShortcode(t *testing.T) {
	t.Run("return emoji by shortcode, if exists", func(t *testing.T) {
		allEmoji := emoji.NewAllEmoji()
		emojivotoService := EmojiServiceServer{
			allEmoji: allEmoji,
		}

		emojiSearchedFor := allEmoji.List()[3]

		response, err := emojivotoService.FindByShortcode(context.Background(), &pb.FindByShortcodeRequest{
			Shortcode: emojiSearchedFor.Shortcode,
		})

		if err != nil {
			t.Fatal(err)
		}

		if response.Emoji == nil {
			t.Fatal("Didnt return an emoji")
		}

		if response.Emoji.Shortcode != emojiSearchedFor.Shortcode || response.Emoji.Unicode != emojiSearchedFor.Unicode {
			t.Fatalf("Response didnt contain [%v]", response.Emoji)
		}
	})

	t.Run("return nil if no emoji with such shortcode", func(t *testing.T) {
		allEmoji := emoji.NewAllEmoji()
		emojivotoService := EmojiServiceServer{
			allEmoji: allEmoji,
		}

		response, err := emojivotoService.FindByShortcode(context.Background(), &pb.FindByShortcodeRequest{
			Shortcode: "doesnt-really-exist",
		})

		if err != nil {
			t.Fatal(err)
		}

		if response.Emoji != nil {
			t.Fatalf("Expected to return nil for emoji, returned [%s]", response.Emoji)
		}
	})

}

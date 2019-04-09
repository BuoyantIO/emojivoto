package emoji

import (
	"testing"
)

func TestListAllEmoji(t *testing.T) {
	allEmoji := NewAllEmoji()

	t.Run("has all emoji from the generated code map", func(t *testing.T) {
		allEmojiSet := make(map[string]bool)
		for _, emoji := range allEmoji.List() {
			allEmojiSet[emoji.Unicode] = true
		}

		for _, alias := range top100Emoji {
			if !allEmojiSet[emojiCodeMap[alias]] {

				t.Fatalf("All Emoji doesnt contain [%s]", emojiCodeMap[alias])
			}
		}
	})

	t.Run("has all emoji from the generated code map", func(t *testing.T) {
		allEmoji := NewAllEmoji()

		alreadySeenEmojiCodes := make(map[string]bool, 0)

		for _, emoji := range allEmoji.List() {
			if alreadySeenEmojiCodes[emoji.Unicode] {
				t.Fatalf("Duplicated emoji [%v]", emoji)
			} else {
				alreadySeenEmojiCodes[emoji.Shortcode] = true
			}
		}
	})
}

func TestFindByShortcode(t *testing.T) {
	allEmoji := NewAllEmoji()

	t.Run("can find all emoji", func(t *testing.T) {
		for _, emoji := range allEmoji.List() {
			found := allEmoji.WithShortcode(emoji.Shortcode)
			if found != emoji {
				t.Fatalf("Couldn't find [%s] by shortcode", emoji)
			}

		}
	})

	t.Run("returns false when can't find emoji", func(t *testing.T) {
		for _, shortcode := range []string{"these", "arent", "emoji", "shortcodes"} {
			found := allEmoji.WithShortcode(shortcode)
			if found != nil {
				t.Fatalf("Returned unexpected [%v]for shortcode [%s]", found, shortcode)
			}

		}
	})
}

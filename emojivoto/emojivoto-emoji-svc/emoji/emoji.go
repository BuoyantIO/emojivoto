package emoji

//go:generate generateEmojiCodeMap -pkg emojivoto

type Emoji struct {
	Unicode   string `json:"unicode"`
	Shortcode string `json:"shortcode"`
}

type AllEmoji interface {
	WithShortcode(shortcode string) *Emoji
	List() []*Emoji
}

type inMemoryAllEmoji struct {
	emojiList []*Emoji
}

func (allEmoji *inMemoryAllEmoji) List() []*Emoji {
	return allEmoji.emojiList
}

func (allEmoji *inMemoryAllEmoji) WithShortcode(shortcode string) *Emoji {
	for _, emoji := range allEmoji.List() {
		if emoji.Shortcode == shortcode {
			return emoji
		}
	}
	return nil
}

func NewAllEmoji() AllEmoji {
	emojiSet := make(map[string]*Emoji, 0)

	for name, unicode := range emojiCodeMap {
		emojiSet[unicode] = &Emoji{
			Unicode:   unicode,
			Shortcode: name,
		}
	}

	emojiList := make([]*Emoji, 0)
	for _, e := range emojiSet {
		emojiList = append(emojiList, e)
	}

	return &inMemoryAllEmoji{
		emojiList,
	}
}

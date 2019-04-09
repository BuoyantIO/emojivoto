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

var top100Emoji = []string{
	":joy:",
	":sunglasses:",
	":doughnut:",
	":stuck_out_tongue_winking_eye:",
	":money_mouth_face:",
	":flushed:",
	":mask:",
	":nerd_face:",
	":ghost:",
	":skull_and_crossbones:",
	":heart_eyes_cat:",
	":hear_no_evil:",
	":see_no_evil:",
	":speak_no_evil:",
	":boy:",
	":girl:",
	":man:",
	":woman:",
	":older_man:",
	":policeman:",
	":guardsman:",
	":construction_worker_man:",
	":prince:",
	":princess:",
	":man_in_tuxedo:",
	":bride_with_veil:",
	":mrs_claus:",
	":santa:",
	":turkey:",
	":rabbit:",
	":no_good_woman:",
	":ok_woman:",
	":raising_hand_woman:",
	":bowing_man:",
	":man_facepalming:",
	":woman_shrugging:",
	":massage_woman:",
	":walking_man:",
	":running_man:",
	":dancer:",
	":man_dancing:",
	":dancing_women:",
	":rainbow:",
	":skier:",
	":golfing_man:",
	":surfing_man:",
	":basketball_man:",
	":biking_man:",
	":point_up_2:",
	":vulcan_salute:",
	":metal:",
	":call_me_hand:",
	":thumbsup:",
	":wave:",
	":clap:",
	":raised_hands:",
	":pray:",
	":dog:",
	":cat2:",
	":pig:",
	":hatching_chick:",
	":snail:",
	":bacon:",
	":pizza:",
	":taco:",
	":burrito:",
	":ramen:",
	":champagne:",
	":tropical_drink:",
	":beer:",
	":tumbler_glass:",
	":world_map:",
	":beach_umbrella:",
	":mountain_snow:",
	":camping:",
	":steam_locomotive:",
	":flight_departure:",
	":rocket:",
	":star2:",
	":sun_behind_small_cloud:",
	":cloud_with_rain:",
	":fire:",
	":jack_o_lantern:",
	":balloon:",
	":tada:",
	":trophy:",
	":iphone:",
	":pager:",
	":fax:",
	":bulb:",
	":money_with_wings:",
	":crystal_ball:",
	":underage:",
	":interrobang:",
	":100:",
	":checkered_flag:",
	":crossed_swords:",
	":floppy_disk:",
	":poop:",
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
	emojiList := make([]*Emoji, 0)

	for _, name := range top100Emoji {
		e := &Emoji{
			Unicode:   emojiCodeMap[name],
			Shortcode: name,
		}
		emojiList = append(emojiList, e)
	}

	return &inMemoryAllEmoji{
		emojiList,
	}
}

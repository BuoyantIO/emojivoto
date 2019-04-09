package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	pb "github.com/buoyantio/emojivoto/emojivoto-web/gen/proto"
	"google.golang.org/grpc"
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

func (c *MockVotingServiceClient) vote(shortcode string) (*pb.VoteResponse, error) {
	c.lastChoiceShortcode = shortcode
	return &pb.VoteResponse{}, nil
}

func (c *MockVotingServiceClient) VoteDoughnut(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return nil, fmt.Errorf("ERROR")
}

func (c *MockVotingServiceClient) VotePoop(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":poop:")
}

func (c *MockVotingServiceClient) VoteJoy(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":joy:")
}

func (c *MockVotingServiceClient) VoteSunglasses(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":sunglasses:")
}

func (c *MockVotingServiceClient) VoteRelaxed(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":relaxed:")
}

func (c *MockVotingServiceClient) VoteStuckOutTongueWinkingEye(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":stuck_out_tongue_winking_eye:")
}

func (c *MockVotingServiceClient) VoteMoneyMouthFace(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":money_mouth_face:")
}

func (c *MockVotingServiceClient) VoteFlushed(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":flushed:")
}

func (c *MockVotingServiceClient) VoteMask(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":mask:")
}

func (c *MockVotingServiceClient) VoteNerdFace(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":nerd_face:")
}

func (c *MockVotingServiceClient) VoteGhost(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":ghost:")
}

func (c *MockVotingServiceClient) VoteSkullAndCrossbones(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":skull_and_crossbones:")
}

func (c *MockVotingServiceClient) VoteHeartEyesCat(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":heart_eyes_cat:")
}

func (c *MockVotingServiceClient) VoteHearNoEvil(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":hear_no_evil:")
}

func (c *MockVotingServiceClient) VoteSeeNoEvil(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":see_no_evil:")
}

func (c *MockVotingServiceClient) VoteSpeakNoEvil(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":speak_no_evil:")
}

func (c *MockVotingServiceClient) VoteBoy(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":boy:")
}

func (c *MockVotingServiceClient) VoteGirl(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":girl:")
}

func (c *MockVotingServiceClient) VoteMan(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":man:")
}

func (c *MockVotingServiceClient) VoteWoman(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":woman:")
}

func (c *MockVotingServiceClient) VoteOlderMan(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":older_man:")
}

func (c *MockVotingServiceClient) VotePoliceman(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":policeman:")
}

func (c *MockVotingServiceClient) VoteGuardsman(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":guardsman:")
}

func (c *MockVotingServiceClient) VoteConstructionWorkerMan(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":construction_worker_man:")
}

func (c *MockVotingServiceClient) VotePrince(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":prince:")
}

func (c *MockVotingServiceClient) VotePrincess(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":princess:")
}

func (c *MockVotingServiceClient) VoteManInTuxedo(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":man_in_tuxedo:")
}

func (c *MockVotingServiceClient) VoteBrideWithVeil(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":bride_with_veil:")
}

func (c *MockVotingServiceClient) VoteMrsClaus(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":mrs_claus:")
}

func (c *MockVotingServiceClient) VoteSanta(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":santa:")
}

func (c *MockVotingServiceClient) VoteTurkey(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":turkey:")
}

func (c *MockVotingServiceClient) VoteRabbit(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":rabbit:")
}

func (c *MockVotingServiceClient) VoteNoGoodWoman(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":no_good_woman:")
}

func (c *MockVotingServiceClient) VoteOkWoman(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":ok_woman:")
}

func (c *MockVotingServiceClient) VoteRaisingHandWoman(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":raising_hand_woman:")
}

func (c *MockVotingServiceClient) VoteBowingMan(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":bowing_man:")
}

func (c *MockVotingServiceClient) VoteManFacepalming(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":man_facepalming:")
}

func (c *MockVotingServiceClient) VoteWomanShrugging(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":woman_shrugging:")
}

func (c *MockVotingServiceClient) VoteMassageWoman(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":massage_woman:")
}

func (c *MockVotingServiceClient) VoteWalkingMan(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":walking_man:")
}

func (c *MockVotingServiceClient) VoteRunningMan(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":running_man:")
}

func (c *MockVotingServiceClient) VoteDancer(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":dancer:")
}

func (c *MockVotingServiceClient) VoteManDancing(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":man_dancing:")
}

func (c *MockVotingServiceClient) VoteDancingWomen(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":dancing_women:")
}

func (c *MockVotingServiceClient) VoteRainbow(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":rainbow:")
}

func (c *MockVotingServiceClient) VoteSkier(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":skier:")
}

func (c *MockVotingServiceClient) VoteGolfingMan(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":golfing_man:")
}

func (c *MockVotingServiceClient) VoteSurfingMan(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":surfing_man:")
}

func (c *MockVotingServiceClient) VoteBasketballMan(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":basketball_man:")
}

func (c *MockVotingServiceClient) VoteBikingMan(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":biking_man:")
}

func (c *MockVotingServiceClient) VotePointUp2(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":point_up_2:")
}

func (c *MockVotingServiceClient) VoteVulcanSalute(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":vulcan_salute:")
}

func (c *MockVotingServiceClient) VoteMetal(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":metal:")
}

func (c *MockVotingServiceClient) VoteCallMeHand(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":call_me_hand:")
}

func (c *MockVotingServiceClient) VoteThumbsup(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":thumbsup:")
}

func (c *MockVotingServiceClient) VoteWave(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":wave:")
}

func (c *MockVotingServiceClient) VoteClap(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":clap:")
}

func (c *MockVotingServiceClient) VoteRaisedHands(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":raised_hands:")
}

func (c *MockVotingServiceClient) VotePray(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":pray:")
}

func (c *MockVotingServiceClient) VoteDog(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":dog:")
}

func (c *MockVotingServiceClient) VoteCat2(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":cat2:")
}

func (c *MockVotingServiceClient) VotePig(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":pig:")
}

func (c *MockVotingServiceClient) VoteHatchingChick(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":hatching_chick:")
}

func (c *MockVotingServiceClient) VoteSnail(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":snail:")
}

func (c *MockVotingServiceClient) VoteBacon(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":bacon:")
}

func (c *MockVotingServiceClient) VotePizza(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":pizza:")
}

func (c *MockVotingServiceClient) VoteTaco(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":taco:")
}

func (c *MockVotingServiceClient) VoteBurrito(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":burrito:")
}

func (c *MockVotingServiceClient) VoteRamen(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":ramen:")
}

func (c *MockVotingServiceClient) VoteChampagne(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":champagne:")
}

func (c *MockVotingServiceClient) VoteTropicalDrink(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":tropical_drink:")
}

func (c *MockVotingServiceClient) VoteBeer(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":beer:")
}

func (c *MockVotingServiceClient) VoteTumblerGlass(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":tumbler_glass:")
}

func (c *MockVotingServiceClient) VoteWorldMap(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":world_map:")
}

func (c *MockVotingServiceClient) VoteBeachUmbrella(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":beach_umbrella:")
}

func (c *MockVotingServiceClient) VoteMountainSnow(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":mountain_snow:")
}

func (c *MockVotingServiceClient) VoteCamping(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":camping:")
}

func (c *MockVotingServiceClient) VoteSteamLocomotive(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":steam_locomotive:")
}

func (c *MockVotingServiceClient) VoteFlightDeparture(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":flight_departure:")
}

func (c *MockVotingServiceClient) VoteRocket(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":rocket:")
}

func (c *MockVotingServiceClient) VoteStar2(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":star2:")
}

func (c *MockVotingServiceClient) VoteSunBehindSmallCloud(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":sun_behind_small_cloud:")
}

func (c *MockVotingServiceClient) VoteCloudWithRain(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":cloud_with_rain:")
}

func (c *MockVotingServiceClient) VoteFire(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":fire:")
}

func (c *MockVotingServiceClient) VoteJackOLantern(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":jack_o_lantern:")
}

func (c *MockVotingServiceClient) VoteBalloon(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":balloon:")
}

func (c *MockVotingServiceClient) VoteTada(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":tada:")
}

func (c *MockVotingServiceClient) VoteTrophy(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":trophy:")
}

func (c *MockVotingServiceClient) VoteIphone(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":iphone:")
}

func (c *MockVotingServiceClient) VotePager(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":pager:")
}

func (c *MockVotingServiceClient) VoteFax(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":fax:")
}

func (c *MockVotingServiceClient) VoteBulb(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":bulb:")
}

func (c *MockVotingServiceClient) VoteMoneyWithWings(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":money_with_wings:")
}

func (c *MockVotingServiceClient) VoteCrystalBall(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":crystal_ball:")
}

func (c *MockVotingServiceClient) VoteUnderage(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":underage:")
}

func (c *MockVotingServiceClient) VoteInterrobang(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":interrobang:")
}

func (c *MockVotingServiceClient) Vote100(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":100:")
}

func (c *MockVotingServiceClient) VoteCheckeredFlag(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":checkered_flag:")
}

func (c *MockVotingServiceClient) VoteCrossedSwords(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":crossed_swords:")
}

func (c *MockVotingServiceClient) VoteFloppyDisk(_ context.Context, _ *pb.VoteRequest, _ ...grpc.CallOption) (*pb.VoteResponse, error) {
	return c.vote(":floppy_disk:")
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
		emojiIWantToVoteFor := &pb.Emoji{Shortcode: ":100:", Unicode: "\U0001f4af"}
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
		webApp := &WebApp{}

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
		emojiIWantToVoteFor := &pb.Emoji{Shortcode: ":100:", Unicode: "\U0001f4af"}
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
			{Shortcode: ":100:", Unicode: "\U0001f4af"},
			{Shortcode: ":checkered_flag:", Unicode: "\U0001f3c1"},
			{Shortcode: ":crossed_swords:", Unicode: "\u2694\ufe0f"},
			{Shortcode: ":floppy_disk:", Unicode: "\U0001f4be"},
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

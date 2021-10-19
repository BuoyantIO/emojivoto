package api

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	pb "github.com/buoyantio/emojivoto/emojivoto-voting-svc/gen/proto"
	"github.com/buoyantio/emojivoto/emojivoto-voting-svc/voting"
	"google.golang.org/grpc"
)

var (
	FloatZero = float32(0.0)
)

type PollServiceServer struct {
	poll                    voting.Poll
	failureRate             float32
	artificialDelayDuration time.Duration
	pb.UnimplementedVotingServiceServer
}

func (pS *PollServiceServer) vote(shortcode string) (*pb.VoteResponse, error) {

	time.Sleep(pS.artificialDelayDuration)

	err := pS.poll.Vote(shortcode)
	return &pb.VoteResponse{}, err
}

func (pS *PollServiceServer) VoteDoughnut(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {

	if pS.failureRate > FloatZero {
		probability := rand.Float32()

		if probability < pS.failureRate {
			log.Printf("probability [%f] is less than failureRate [%f]", probability, pS.failureRate)
			log.Printf("logging an error for doughnut")
			return nil, fmt.Errorf("ERROR")
		}
	}
	log.Printf("voting for doughnut")
	return pS.vote(":doughnut:")
}

func (pS *PollServiceServer) VotePoop(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":poop:")
}

func (pS *PollServiceServer) VoteJoy(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":joy:")
}

func (pS *PollServiceServer) VoteSunglasses(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":sunglasses:")
}

func (pS *PollServiceServer) VoteRelaxed(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":relaxed:")
}

func (pS *PollServiceServer) VoteStuckOutTongueWinkingEye(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":stuck_out_tongue_winking_eye:")
}

func (pS *PollServiceServer) VoteMoneyMouthFace(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":money_mouth_face:")
}

func (pS *PollServiceServer) VoteFlushed(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":flushed:")
}

func (pS *PollServiceServer) VoteMask(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":mask:")
}

func (pS *PollServiceServer) VoteNerdFace(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":nerd_face:")
}

func (pS *PollServiceServer) VoteGhost(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":ghost:")
}

func (pS *PollServiceServer) VoteSkullAndCrossbones(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":skull_and_crossbones:")
}

func (pS *PollServiceServer) VoteHeartEyesCat(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":heart_eyes_cat:")
}

func (pS *PollServiceServer) VoteHearNoEvil(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":hear_no_evil:")
}

func (pS *PollServiceServer) VoteSeeNoEvil(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":see_no_evil:")
}

func (pS *PollServiceServer) VoteSpeakNoEvil(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":speak_no_evil:")
}

func (pS *PollServiceServer) VoteBoy(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":boy:")
}

func (pS *PollServiceServer) VoteGirl(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":girl:")
}

func (pS *PollServiceServer) VoteMan(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":man:")
}

func (pS *PollServiceServer) VoteWoman(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":woman:")
}

func (pS *PollServiceServer) VoteOlderMan(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":older_man:")
}

func (pS *PollServiceServer) VotePoliceman(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":policeman:")
}

func (pS *PollServiceServer) VoteGuardsman(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":guardsman:")
}

func (pS *PollServiceServer) VoteConstructionWorkerMan(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":construction_worker_man:")
}

func (pS *PollServiceServer) VotePrince(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":prince:")
}

func (pS *PollServiceServer) VotePrincess(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":princess:")
}

func (pS *PollServiceServer) VoteManInTuxedo(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":man_in_tuxedo:")
}

func (pS *PollServiceServer) VoteBrideWithVeil(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":bride_with_veil:")
}

func (pS *PollServiceServer) VoteMrsClaus(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":mrs_claus:")
}

func (pS *PollServiceServer) VoteSanta(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":santa:")
}

func (pS *PollServiceServer) VoteTurkey(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":turkey:")
}

func (pS *PollServiceServer) VoteRabbit(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":rabbit:")
}

func (pS *PollServiceServer) VoteNoGoodWoman(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":no_good_woman:")
}

func (pS *PollServiceServer) VoteOkWoman(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":ok_woman:")
}

func (pS *PollServiceServer) VoteRaisingHandWoman(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":raising_hand_woman:")
}

func (pS *PollServiceServer) VoteBowingMan(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":bowing_man:")
}

func (pS *PollServiceServer) VoteManFacepalming(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":man_facepalming:")
}

func (pS *PollServiceServer) VoteWomanShrugging(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":woman_shrugging:")
}

func (pS *PollServiceServer) VoteMassageWoman(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":massage_woman:")
}

func (pS *PollServiceServer) VoteWalkingMan(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":walking_man:")
}

func (pS *PollServiceServer) VoteRunningMan(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":running_man:")
}

func (pS *PollServiceServer) VoteDancer(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":dancer:")
}

func (pS *PollServiceServer) VoteManDancing(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":man_dancing:")
}

func (pS *PollServiceServer) VoteDancingWomen(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":dancing_women:")
}

func (pS *PollServiceServer) VoteRainbow(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":rainbow:")
}

func (pS *PollServiceServer) VoteSkier(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":skier:")
}

func (pS *PollServiceServer) VoteGolfingMan(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":golfing_man:")
}

func (pS *PollServiceServer) VoteSurfingMan(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":surfing_man:")
}

func (pS *PollServiceServer) VoteBasketballMan(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":basketball_man:")
}

func (pS *PollServiceServer) VoteBikingMan(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":biking_man:")
}

func (pS *PollServiceServer) VotePointUp2(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":point_up_2:")
}

func (pS *PollServiceServer) VoteVulcanSalute(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":vulcan_salute:")
}

func (pS *PollServiceServer) VoteMetal(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":metal:")
}

func (pS *PollServiceServer) VoteCallMeHand(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":call_me_hand:")
}

func (pS *PollServiceServer) VoteThumbsup(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":thumbsup:")
}

func (pS *PollServiceServer) VoteWave(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":wave:")
}

func (pS *PollServiceServer) VoteClap(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":clap:")
}

func (pS *PollServiceServer) VoteRaisedHands(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":raised_hands:")
}

func (pS *PollServiceServer) VotePray(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":pray:")
}

func (pS *PollServiceServer) VoteDog(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":dog:")
}

func (pS *PollServiceServer) VoteCat2(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":cat2:")
}

func (pS *PollServiceServer) VotePig(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":pig:")
}

func (pS *PollServiceServer) VoteHatchingChick(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":hatching_chick:")
}

func (pS *PollServiceServer) VoteSnail(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":snail:")
}

func (pS *PollServiceServer) VoteBacon(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":bacon:")
}

func (pS *PollServiceServer) VotePizza(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":pizza:")
}

func (pS *PollServiceServer) VoteTaco(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":taco:")
}

func (pS *PollServiceServer) VoteBurrito(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":burrito:")
}

func (pS *PollServiceServer) VoteRamen(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":ramen:")
}

func (pS *PollServiceServer) VoteChampagne(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":champagne:")
}

func (pS *PollServiceServer) VoteTropicalDrink(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":tropical_drink:")
}

func (pS *PollServiceServer) VoteBeer(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":beer:")
}

func (pS *PollServiceServer) VoteTumblerGlass(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":tumbler_glass:")
}

func (pS *PollServiceServer) VoteWorldMap(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":world_map:")
}

func (pS *PollServiceServer) VoteBeachUmbrella(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":beach_umbrella:")
}

func (pS *PollServiceServer) VoteMountainSnow(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":mountain_snow:")
}

func (pS *PollServiceServer) VoteCamping(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":camping:")
}

func (pS *PollServiceServer) VoteSteamLocomotive(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":steam_locomotive:")
}

func (pS *PollServiceServer) VoteFlightDeparture(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":flight_departure:")
}

func (pS *PollServiceServer) VoteRocket(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":rocket:")
}

func (pS *PollServiceServer) VoteStar2(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":star2:")
}

func (pS *PollServiceServer) VoteSunBehindSmallCloud(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":sun_behind_small_cloud:")
}

func (pS *PollServiceServer) VoteCloudWithRain(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":cloud_with_rain:")
}

func (pS *PollServiceServer) VoteFire(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":fire:")
}

func (pS *PollServiceServer) VoteJackOLantern(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":jack_o_lantern:")
}

func (pS *PollServiceServer) VoteBalloon(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":balloon:")
}

func (pS *PollServiceServer) VoteTada(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":tada:")
}

func (pS *PollServiceServer) VoteTrophy(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":trophy:")
}

func (pS *PollServiceServer) VoteIphone(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":iphone:")
}

func (pS *PollServiceServer) VotePager(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":pager:")
}

func (pS *PollServiceServer) VoteFax(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":fax:")
}

func (pS *PollServiceServer) VoteBulb(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":bulb:")
}

func (pS *PollServiceServer) VoteMoneyWithWings(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":money_with_wings:")
}

func (pS *PollServiceServer) VoteCrystalBall(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":crystal_ball:")
}

func (pS *PollServiceServer) VoteUnderage(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":underage:")
}

func (pS *PollServiceServer) VoteInterrobang(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":interrobang:")
}

func (pS *PollServiceServer) Vote100(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":100:")
}

func (pS *PollServiceServer) VoteCheckeredFlag(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":checkered_flag:")
}

func (pS *PollServiceServer) VoteCrossedSwords(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":crossed_swords:")
}

func (pS *PollServiceServer) VoteFloppyDisk(_ context.Context, _ *pb.VoteRequest) (*pb.VoteResponse, error) {
	return pS.vote(":floppy_disk:")
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

func NewGrpServer(grpcServer *grpc.Server, poll voting.Poll, failureRate float32, artificialDelayDuration time.Duration) {
	server := &PollServiceServer{
		poll,
		failureRate,
		artificialDelayDuration,
		pb.UnimplementedVotingServiceServer{},
	}

	pb.RegisterVotingServiceServer(grpcServer, server)
}

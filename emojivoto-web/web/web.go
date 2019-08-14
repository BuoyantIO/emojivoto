package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"go.opencensus.io/plugin/ochttp"

	pb "github.com/buoyantio/emojivoto/emojivoto-web/gen/proto"
)

type WebApp struct {
	emojiServiceClient  pb.EmojiServiceClient
	votingServiceClient pb.VotingServiceClient
	indexBundle         string
	webpackDevServer    string
}

func (app *WebApp) listEmojiHandler(w http.ResponseWriter, r *http.Request) {
	serviceResponse, err := app.emojiServiceClient.ListAll(r.Context(), &pb.ListAllEmojiRequest{})
	if err != nil {
		writeError(err, w, r, http.StatusInternalServerError)
		return
	}

	list := make([]map[string]string, 0)
	for _, e := range serviceResponse.List {
		list = append(list, map[string]string{
			"shortcode": e.Shortcode,
			"unicode":   e.Unicode,
		})
	}

	err = writeJsonBody(w, http.StatusOK, list)

	if err != nil {
		writeError(err, w, r, http.StatusInternalServerError)
	}
}

func (app *WebApp) leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	results, err := app.votingServiceClient.Results(r.Context(), &pb.ResultsRequest{})

	if err != nil {
		writeError(err, w, r, http.StatusInternalServerError)
		return
	}

	representations := make([]map[string]string, 0)
	for _, result := range results.Results {
		findByShortcodeRequest := &pb.FindByShortcodeRequest{
			Shortcode: result.Shortcode,
		}

		findByShortcodeResponse, err := app.emojiServiceClient.FindByShortcode(r.Context(), findByShortcodeRequest)

		if err != nil {
			writeError(err, w, r, http.StatusInternalServerError)
			return
		}

		emoji := findByShortcodeResponse.Emoji
		representation := make(map[string]string)
		representation["votes"] = strconv.Itoa(int(result.Votes))
		representation["unicode"] = emoji.Unicode
		representation["shortcode"] = emoji.Shortcode

		representations = append(representations, representation)
	}

	err = writeJsonBody(w, http.StatusOK, representations)

	if err != nil {
		writeError(err, w, r, http.StatusInternalServerError)
	}
}

func (app *WebApp) voteEmojiHandler(w http.ResponseWriter, r *http.Request) {
	emojiShortcode := r.FormValue("choice")
	if emojiShortcode == "" {
		error := errors.New(fmt.Sprintf("Emoji choice [%s] is mandatory", emojiShortcode))
		writeError(error, w, r, http.StatusBadRequest)
		return
	}

	request := &pb.FindByShortcodeRequest{
		Shortcode: emojiShortcode,
	}
	response, err := app.emojiServiceClient.FindByShortcode(r.Context(), request)
	if err != nil {
		writeError(err, w, r, http.StatusInternalServerError)
		return
	}

	if response.Emoji == nil {
		err = errors.New(fmt.Sprintf("Choosen emoji shortcode [%s] doesnt exist", emojiShortcode))
		writeError(err, w, r, http.StatusBadRequest)
		return
	}

	voteRequest := &pb.VoteRequest{}
	switch emojiShortcode {
	case ":poop:":
		_, err = app.votingServiceClient.VotePoop(r.Context(), voteRequest)
	case ":joy:":
		_, err = app.votingServiceClient.VoteJoy(r.Context(), voteRequest)
	case ":sunglasses:":
		_, err = app.votingServiceClient.VoteSunglasses(r.Context(), voteRequest)
	case ":relaxed:":
		_, err = app.votingServiceClient.VoteRelaxed(r.Context(), voteRequest)
	case ":stuck_out_tongue_winking_eye:":
		_, err = app.votingServiceClient.VoteStuckOutTongueWinkingEye(r.Context(), voteRequest)
	case ":money_mouth_face:":
		_, err = app.votingServiceClient.VoteMoneyMouthFace(r.Context(), voteRequest)
	case ":flushed:":
		_, err = app.votingServiceClient.VoteFlushed(r.Context(), voteRequest)
	case ":mask:":
		_, err = app.votingServiceClient.VoteMask(r.Context(), voteRequest)
	case ":nerd_face:":
		_, err = app.votingServiceClient.VoteNerdFace(r.Context(), voteRequest)
	case ":ghost:":
		_, err = app.votingServiceClient.VoteGhost(r.Context(), voteRequest)
	case ":skull_and_crossbones:":
		_, err = app.votingServiceClient.VoteSkullAndCrossbones(r.Context(), voteRequest)
	case ":heart_eyes_cat:":
		_, err = app.votingServiceClient.VoteHeartEyesCat(r.Context(), voteRequest)
	case ":hear_no_evil:":
		_, err = app.votingServiceClient.VoteHearNoEvil(r.Context(), voteRequest)
	case ":see_no_evil:":
		_, err = app.votingServiceClient.VoteSeeNoEvil(r.Context(), voteRequest)
	case ":speak_no_evil:":
		_, err = app.votingServiceClient.VoteSpeakNoEvil(r.Context(), voteRequest)
	case ":boy:":
		_, err = app.votingServiceClient.VoteBoy(r.Context(), voteRequest)
	case ":girl:":
		_, err = app.votingServiceClient.VoteGirl(r.Context(), voteRequest)
	case ":man:":
		_, err = app.votingServiceClient.VoteMan(r.Context(), voteRequest)
	case ":woman:":
		_, err = app.votingServiceClient.VoteWoman(r.Context(), voteRequest)
	case ":older_man:":
		_, err = app.votingServiceClient.VoteOlderMan(r.Context(), voteRequest)
	case ":policeman:":
		_, err = app.votingServiceClient.VotePoliceman(r.Context(), voteRequest)
	case ":guardsman:":
		_, err = app.votingServiceClient.VoteGuardsman(r.Context(), voteRequest)
	case ":construction_worker_man:":
		_, err = app.votingServiceClient.VoteConstructionWorkerMan(r.Context(), voteRequest)
	case ":prince:":
		_, err = app.votingServiceClient.VotePrince(r.Context(), voteRequest)
	case ":princess:":
		_, err = app.votingServiceClient.VotePrincess(r.Context(), voteRequest)
	case ":man_in_tuxedo:":
		_, err = app.votingServiceClient.VoteManInTuxedo(r.Context(), voteRequest)
	case ":bride_with_veil:":
		_, err = app.votingServiceClient.VoteBrideWithVeil(r.Context(), voteRequest)
	case ":mrs_claus:":
		_, err = app.votingServiceClient.VoteMrsClaus(r.Context(), voteRequest)
	case ":santa:":
		_, err = app.votingServiceClient.VoteSanta(r.Context(), voteRequest)
	case ":turkey:":
		_, err = app.votingServiceClient.VoteTurkey(r.Context(), voteRequest)
	case ":rabbit:":
		_, err = app.votingServiceClient.VoteRabbit(r.Context(), voteRequest)
	case ":no_good_woman:":
		_, err = app.votingServiceClient.VoteNoGoodWoman(r.Context(), voteRequest)
	case ":ok_woman:":
		_, err = app.votingServiceClient.VoteOkWoman(r.Context(), voteRequest)
	case ":raising_hand_woman:":
		_, err = app.votingServiceClient.VoteRaisingHandWoman(r.Context(), voteRequest)
	case ":bowing_man:":
		_, err = app.votingServiceClient.VoteBowingMan(r.Context(), voteRequest)
	case ":man_facepalming:":
		_, err = app.votingServiceClient.VoteManFacepalming(r.Context(), voteRequest)
	case ":woman_shrugging:":
		_, err = app.votingServiceClient.VoteWomanShrugging(r.Context(), voteRequest)
	case ":massage_woman:":
		_, err = app.votingServiceClient.VoteMassageWoman(r.Context(), voteRequest)
	case ":walking_man:":
		_, err = app.votingServiceClient.VoteWalkingMan(r.Context(), voteRequest)
	case ":running_man:":
		_, err = app.votingServiceClient.VoteRunningMan(r.Context(), voteRequest)
	case ":dancer:":
		_, err = app.votingServiceClient.VoteDancer(r.Context(), voteRequest)
	case ":man_dancing:":
		_, err = app.votingServiceClient.VoteManDancing(r.Context(), voteRequest)
	case ":dancing_women:":
		_, err = app.votingServiceClient.VoteDancingWomen(r.Context(), voteRequest)
	case ":rainbow:":
		_, err = app.votingServiceClient.VoteRainbow(r.Context(), voteRequest)
	case ":skier:":
		_, err = app.votingServiceClient.VoteSkier(r.Context(), voteRequest)
	case ":golfing_man:":
		_, err = app.votingServiceClient.VoteGolfingMan(r.Context(), voteRequest)
	case ":surfing_man:":
		_, err = app.votingServiceClient.VoteSurfingMan(r.Context(), voteRequest)
	case ":basketball_man:":
		_, err = app.votingServiceClient.VoteBasketballMan(r.Context(), voteRequest)
	case ":biking_man:":
		_, err = app.votingServiceClient.VoteBikingMan(r.Context(), voteRequest)
	case ":point_up_2:":
		_, err = app.votingServiceClient.VotePointUp2(r.Context(), voteRequest)
	case ":vulcan_salute:":
		_, err = app.votingServiceClient.VoteVulcanSalute(r.Context(), voteRequest)
	case ":metal:":
		_, err = app.votingServiceClient.VoteMetal(r.Context(), voteRequest)
	case ":call_me_hand:":
		_, err = app.votingServiceClient.VoteCallMeHand(r.Context(), voteRequest)
	case ":thumbsup:":
		_, err = app.votingServiceClient.VoteThumbsup(r.Context(), voteRequest)
	case ":wave:":
		_, err = app.votingServiceClient.VoteWave(r.Context(), voteRequest)
	case ":clap:":
		_, err = app.votingServiceClient.VoteClap(r.Context(), voteRequest)
	case ":raised_hands:":
		_, err = app.votingServiceClient.VoteRaisedHands(r.Context(), voteRequest)
	case ":pray:":
		_, err = app.votingServiceClient.VotePray(r.Context(), voteRequest)
	case ":dog:":
		_, err = app.votingServiceClient.VoteDog(r.Context(), voteRequest)
	case ":cat2:":
		_, err = app.votingServiceClient.VoteCat2(r.Context(), voteRequest)
	case ":pig:":
		_, err = app.votingServiceClient.VotePig(r.Context(), voteRequest)
	case ":hatching_chick:":
		_, err = app.votingServiceClient.VoteHatchingChick(r.Context(), voteRequest)
	case ":snail:":
		_, err = app.votingServiceClient.VoteSnail(r.Context(), voteRequest)
	case ":bacon:":
		_, err = app.votingServiceClient.VoteBacon(r.Context(), voteRequest)
	case ":pizza:":
		_, err = app.votingServiceClient.VotePizza(r.Context(), voteRequest)
	case ":taco:":
		_, err = app.votingServiceClient.VoteTaco(r.Context(), voteRequest)
	case ":burrito:":
		_, err = app.votingServiceClient.VoteBurrito(r.Context(), voteRequest)
	case ":ramen:":
		_, err = app.votingServiceClient.VoteRamen(r.Context(), voteRequest)
	case ":doughnut:":
		_, err = app.votingServiceClient.VoteDoughnut(r.Context(), voteRequest)
	case ":champagne:":
		_, err = app.votingServiceClient.VoteChampagne(r.Context(), voteRequest)
	case ":tropical_drink:":
		_, err = app.votingServiceClient.VoteTropicalDrink(r.Context(), voteRequest)
	case ":beer:":
		_, err = app.votingServiceClient.VoteBeer(r.Context(), voteRequest)
	case ":tumbler_glass:":
		_, err = app.votingServiceClient.VoteTumblerGlass(r.Context(), voteRequest)
	case ":world_map:":
		_, err = app.votingServiceClient.VoteWorldMap(r.Context(), voteRequest)
	case ":beach_umbrella:":
		_, err = app.votingServiceClient.VoteBeachUmbrella(r.Context(), voteRequest)
	case ":mountain_snow:":
		_, err = app.votingServiceClient.VoteMountainSnow(r.Context(), voteRequest)
	case ":camping:":
		_, err = app.votingServiceClient.VoteCamping(r.Context(), voteRequest)
	case ":steam_locomotive:":
		_, err = app.votingServiceClient.VoteSteamLocomotive(r.Context(), voteRequest)
	case ":flight_departure:":
		_, err = app.votingServiceClient.VoteFlightDeparture(r.Context(), voteRequest)
	case ":rocket:":
		_, err = app.votingServiceClient.VoteRocket(r.Context(), voteRequest)
	case ":star2:":
		_, err = app.votingServiceClient.VoteStar2(r.Context(), voteRequest)
	case ":sun_behind_small_cloud:":
		_, err = app.votingServiceClient.VoteSunBehindSmallCloud(r.Context(), voteRequest)
	case ":cloud_with_rain:":
		_, err = app.votingServiceClient.VoteCloudWithRain(r.Context(), voteRequest)
	case ":fire:":
		_, err = app.votingServiceClient.VoteFire(r.Context(), voteRequest)
	case ":jack_o_lantern:":
		_, err = app.votingServiceClient.VoteJackOLantern(r.Context(), voteRequest)
	case ":balloon:":
		_, err = app.votingServiceClient.VoteBalloon(r.Context(), voteRequest)
	case ":tada:":
		_, err = app.votingServiceClient.VoteTada(r.Context(), voteRequest)
	case ":trophy:":
		_, err = app.votingServiceClient.VoteTrophy(r.Context(), voteRequest)
	case ":iphone:":
		_, err = app.votingServiceClient.VoteIphone(r.Context(), voteRequest)
	case ":pager:":
		_, err = app.votingServiceClient.VotePager(r.Context(), voteRequest)
	case ":fax:":
		_, err = app.votingServiceClient.VoteFax(r.Context(), voteRequest)
	case ":bulb:":
		_, err = app.votingServiceClient.VoteBulb(r.Context(), voteRequest)
	case ":money_with_wings:":
		_, err = app.votingServiceClient.VoteMoneyWithWings(r.Context(), voteRequest)
	case ":crystal_ball:":
		_, err = app.votingServiceClient.VoteCrystalBall(r.Context(), voteRequest)
	case ":underage:":
		_, err = app.votingServiceClient.VoteUnderage(r.Context(), voteRequest)
	case ":interrobang:":
		_, err = app.votingServiceClient.VoteInterrobang(r.Context(), voteRequest)
	case ":100:":
		_, err = app.votingServiceClient.Vote100(r.Context(), voteRequest)
	case ":checkered_flag:":
		_, err = app.votingServiceClient.VoteCheckeredFlag(r.Context(), voteRequest)
	case ":crossed_swords:":
		_, err = app.votingServiceClient.VoteCrossedSwords(r.Context(), voteRequest)
	case ":floppy_disk:":
		_, err = app.votingServiceClient.VoteFloppyDisk(r.Context(), voteRequest)
	}
	if err != nil {
		writeError(err, w, r, http.StatusInternalServerError)
		return
	}
}

func (app *WebApp) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	indexTemplate := `
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>Emoji Vote</title>
			<link rel="icon" href="/img/favicon.ico">
			<!-- Global site tag (gtag.js) - Google Analytics -->
			<script async src="https://www.googletagmanager.com/gtag/js?id=UA-60040560-4"></script>
			<script>
			  window.dataLayer = window.dataLayer || [];
			  function gtag(){dataLayer.push(arguments);}
			  gtag('js', new Date());
			  gtag('config', 'UA-60040560-4');
			</script>
		</head>
		<body>
			<div id="main" class="main"></div>
		</body>
		{{ if ne . ""}}
			<script type="text/javascript" src="{{ . }}/dist/index_bundle.js" async></script>
		{{else}}
			<script type="text/javascript" src="/js" async></script>
		{{end}}
	</html>`
	t, err := template.New("indexTemplate").Parse(indexTemplate)
	if err != nil {
		panic(err)
	}
	t.Execute(w, app.webpackDevServer)
}

func (app *WebApp) jsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	f, err := ioutil.ReadFile(app.indexBundle)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(f))
}

func (app *WebApp) faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/favicon.ico")
}

func writeJsonBody(w http.ResponseWriter, status int, body interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}

func writeError(err error, w http.ResponseWriter, r *http.Request, status int) {
	log.Printf("Error serving request [%v]: %v", r, err)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	errorMessage := make(map[string]string)
	errorMessage["error"] = fmt.Sprintf("%v", err)
	json.NewEncoder(w).Encode(errorMessage)
}

func handle(path string, h func (w http.ResponseWriter, r *http.Request)) {
	http.Handle(path, &ochttp.Handler {
		Handler: http.HandlerFunc(h),
	})
}

func StartServer(webPort, webpackDevServer, indexBundle string, emojiServiceClient pb.EmojiServiceClient, votingClient pb.VotingServiceClient) {
	webApp := &WebApp{
		emojiServiceClient:  emojiServiceClient,
		votingServiceClient: votingClient,
		indexBundle:         indexBundle,
		webpackDevServer:    webpackDevServer,
	}

	log.Printf("Starting web server on WEB_PORT=[%s]", webPort)
	handle("/", webApp.indexHandler)
	handle("/leaderboard", webApp.indexHandler)
	handle("/js", webApp.jsHandler)
	handle("/img/favicon.ico", webApp.faviconHandler)
	handle("/api/list", webApp.listEmojiHandler)
	handle("/api/vote", webApp.voteEmojiHandler)
	handle("/api/leaderboard", webApp.leaderboardHandler)

	// TODO: make static assets dir configurable
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

	err := http.ListenAndServe(fmt.Sprintf(":%s", webPort), nil)
	if err != nil {
		panic(err)
	}
}

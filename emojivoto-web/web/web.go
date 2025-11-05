package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	pb "github.com/buoyantio/emojivoto/emojivoto-web/gen/proto"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type WebApp struct {
	emojiServiceClient  pb.EmojiServiceClient
	votingServiceClient pb.VotingServiceClient
	indexBundle         string
	webpackDevServer    string
	messageOfTheDay     string
}

func (app *WebApp) listEmojiHandler(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	serviceResponse, err := app.emojiServiceClient.ListAll(ctx, &pb.ListAllEmojiRequest{})
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
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	results, err := app.votingServiceClient.Results(ctx, &pb.ResultsRequest{})

	if err != nil {
		writeError(err, w, r, http.StatusInternalServerError)
		return
	}

	representations := make([]map[string]string, 0)
	for _, result := range results.Results {
		findByShortcodeRequest := &pb.FindByShortcodeRequest{
			Shortcode: result.Shortcode,
		}

		findByShortcodeResponse, err := app.emojiServiceClient.FindByShortcode(ctx, findByShortcodeRequest)

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
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	emojiShortcode := r.FormValue("choice")
	if emojiShortcode == "" {
		error := errors.New(fmt.Sprintf("Emoji choice [%s] is mandatory", emojiShortcode))
		writeError(error, w, r, http.StatusBadRequest)
		return
	}

	request := &pb.FindByShortcodeRequest{
		Shortcode: emojiShortcode,
	}
	response, err := app.emojiServiceClient.FindByShortcode(ctx, request)
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
		_, err = app.votingServiceClient.VotePoop(ctx, voteRequest)
	case ":joy:":
		_, err = app.votingServiceClient.VoteJoy(ctx, voteRequest)
	case ":sunglasses:":
		_, err = app.votingServiceClient.VoteSunglasses(ctx, voteRequest)
	case ":relaxed:":
		_, err = app.votingServiceClient.VoteRelaxed(ctx, voteRequest)
	case ":stuck_out_tongue_winking_eye:":
		_, err = app.votingServiceClient.VoteStuckOutTongueWinkingEye(ctx, voteRequest)
	case ":money_mouth_face:":
		_, err = app.votingServiceClient.VoteMoneyMouthFace(ctx, voteRequest)
	case ":flushed:":
		_, err = app.votingServiceClient.VoteFlushed(ctx, voteRequest)
	case ":mask:":
		_, err = app.votingServiceClient.VoteMask(ctx, voteRequest)
	case ":nerd_face:":
		_, err = app.votingServiceClient.VoteNerdFace(ctx, voteRequest)
	case ":ghost:":
		_, err = app.votingServiceClient.VoteGhost(ctx, voteRequest)
	case ":skull_and_crossbones:":
		_, err = app.votingServiceClient.VoteSkullAndCrossbones(ctx, voteRequest)
	case ":heart_eyes_cat:":
		_, err = app.votingServiceClient.VoteHeartEyesCat(ctx, voteRequest)
	case ":hear_no_evil:":
		_, err = app.votingServiceClient.VoteHearNoEvil(ctx, voteRequest)
	case ":see_no_evil:":
		_, err = app.votingServiceClient.VoteSeeNoEvil(ctx, voteRequest)
	case ":speak_no_evil:":
		_, err = app.votingServiceClient.VoteSpeakNoEvil(ctx, voteRequest)
	case ":boy:":
		_, err = app.votingServiceClient.VoteBoy(ctx, voteRequest)
	case ":girl:":
		_, err = app.votingServiceClient.VoteGirl(ctx, voteRequest)
	case ":man:":
		_, err = app.votingServiceClient.VoteMan(ctx, voteRequest)
	case ":woman:":
		_, err = app.votingServiceClient.VoteWoman(ctx, voteRequest)
	case ":older_man:":
		_, err = app.votingServiceClient.VoteOlderMan(ctx, voteRequest)
	case ":policeman:":
		_, err = app.votingServiceClient.VotePoliceman(ctx, voteRequest)
	case ":guardsman:":
		_, err = app.votingServiceClient.VoteGuardsman(ctx, voteRequest)
	case ":construction_worker_man:":
		_, err = app.votingServiceClient.VoteConstructionWorkerMan(ctx, voteRequest)
	case ":prince:":
		_, err = app.votingServiceClient.VotePrince(ctx, voteRequest)
	case ":princess:":
		_, err = app.votingServiceClient.VotePrincess(ctx, voteRequest)
	case ":man_in_tuxedo:":
		_, err = app.votingServiceClient.VoteManInTuxedo(ctx, voteRequest)
	case ":bride_with_veil:":
		_, err = app.votingServiceClient.VoteBrideWithVeil(ctx, voteRequest)
	case ":mrs_claus:":
		_, err = app.votingServiceClient.VoteMrsClaus(ctx, voteRequest)
	case ":santa:":
		_, err = app.votingServiceClient.VoteSanta(ctx, voteRequest)
	case ":turkey:":
		_, err = app.votingServiceClient.VoteTurkey(ctx, voteRequest)
	case ":rabbit:":
		_, err = app.votingServiceClient.VoteRabbit(ctx, voteRequest)
	case ":no_good_woman:":
		_, err = app.votingServiceClient.VoteNoGoodWoman(ctx, voteRequest)
	case ":ok_woman:":
		_, err = app.votingServiceClient.VoteOkWoman(ctx, voteRequest)
	case ":raising_hand_woman:":
		_, err = app.votingServiceClient.VoteRaisingHandWoman(ctx, voteRequest)
	case ":bowing_man:":
		_, err = app.votingServiceClient.VoteBowingMan(ctx, voteRequest)
	case ":man_facepalming:":
		_, err = app.votingServiceClient.VoteManFacepalming(ctx, voteRequest)
	case ":woman_shrugging:":
		_, err = app.votingServiceClient.VoteWomanShrugging(ctx, voteRequest)
	case ":massage_woman:":
		_, err = app.votingServiceClient.VoteMassageWoman(ctx, voteRequest)
	case ":walking_man:":
		_, err = app.votingServiceClient.VoteWalkingMan(ctx, voteRequest)
	case ":running_man:":
		_, err = app.votingServiceClient.VoteRunningMan(ctx, voteRequest)
	case ":dancer:":
		_, err = app.votingServiceClient.VoteDancer(ctx, voteRequest)
	case ":man_dancing:":
		_, err = app.votingServiceClient.VoteManDancing(ctx, voteRequest)
	case ":dancing_women:":
		_, err = app.votingServiceClient.VoteDancingWomen(ctx, voteRequest)
	case ":rainbow:":
		_, err = app.votingServiceClient.VoteRainbow(ctx, voteRequest)
	case ":skier:":
		_, err = app.votingServiceClient.VoteSkier(ctx, voteRequest)
	case ":golfing_man:":
		_, err = app.votingServiceClient.VoteGolfingMan(ctx, voteRequest)
	case ":surfing_man:":
		_, err = app.votingServiceClient.VoteSurfingMan(ctx, voteRequest)
	case ":basketball_man:":
		_, err = app.votingServiceClient.VoteBasketballMan(ctx, voteRequest)
	case ":biking_man:":
		_, err = app.votingServiceClient.VoteBikingMan(ctx, voteRequest)
	case ":point_up_2:":
		_, err = app.votingServiceClient.VotePointUp2(ctx, voteRequest)
	case ":vulcan_salute:":
		_, err = app.votingServiceClient.VoteVulcanSalute(ctx, voteRequest)
	case ":metal:":
		_, err = app.votingServiceClient.VoteMetal(ctx, voteRequest)
	case ":call_me_hand:":
		_, err = app.votingServiceClient.VoteCallMeHand(ctx, voteRequest)
	case ":thumbsup:":
		_, err = app.votingServiceClient.VoteThumbsup(ctx, voteRequest)
	case ":wave:":
		_, err = app.votingServiceClient.VoteWave(ctx, voteRequest)
	case ":clap:":
		_, err = app.votingServiceClient.VoteClap(ctx, voteRequest)
	case ":raised_hands:":
		_, err = app.votingServiceClient.VoteRaisedHands(ctx, voteRequest)
	case ":pray:":
		_, err = app.votingServiceClient.VotePray(ctx, voteRequest)
	case ":dog:":
		_, err = app.votingServiceClient.VoteDog(ctx, voteRequest)
	case ":cat2:":
		_, err = app.votingServiceClient.VoteCat2(ctx, voteRequest)
	case ":pig:":
		_, err = app.votingServiceClient.VotePig(ctx, voteRequest)
	case ":hatching_chick:":
		_, err = app.votingServiceClient.VoteHatchingChick(ctx, voteRequest)
	case ":snail:":
		_, err = app.votingServiceClient.VoteSnail(ctx, voteRequest)
	case ":bacon:":
		_, err = app.votingServiceClient.VoteBacon(ctx, voteRequest)
	case ":pizza:":
		_, err = app.votingServiceClient.VotePizza(ctx, voteRequest)
	case ":taco:":
		_, err = app.votingServiceClient.VoteTaco(ctx, voteRequest)
	case ":burrito:":
		_, err = app.votingServiceClient.VoteBurrito(ctx, voteRequest)
	case ":ramen:":
		_, err = app.votingServiceClient.VoteRamen(ctx, voteRequest)
	case ":doughnut:":
		_, err = app.votingServiceClient.VoteDoughnut(ctx, voteRequest)
	case ":champagne:":
		_, err = app.votingServiceClient.VoteChampagne(ctx, voteRequest)
	case ":tropical_drink:":
		_, err = app.votingServiceClient.VoteTropicalDrink(ctx, voteRequest)
	case ":beer:":
		_, err = app.votingServiceClient.VoteBeer(ctx, voteRequest)
	case ":tumbler_glass:":
		_, err = app.votingServiceClient.VoteTumblerGlass(ctx, voteRequest)
	case ":world_map:":
		_, err = app.votingServiceClient.VoteWorldMap(ctx, voteRequest)
	case ":beach_umbrella:":
		_, err = app.votingServiceClient.VoteBeachUmbrella(ctx, voteRequest)
	case ":mountain_snow:":
		_, err = app.votingServiceClient.VoteMountainSnow(ctx, voteRequest)
	case ":camping:":
		_, err = app.votingServiceClient.VoteCamping(ctx, voteRequest)
	case ":steam_locomotive:":
		_, err = app.votingServiceClient.VoteSteamLocomotive(ctx, voteRequest)
	case ":flight_departure:":
		_, err = app.votingServiceClient.VoteFlightDeparture(ctx, voteRequest)
	case ":rocket:":
		_, err = app.votingServiceClient.VoteRocket(ctx, voteRequest)
	case ":star2:":
		_, err = app.votingServiceClient.VoteStar2(ctx, voteRequest)
	case ":sun_behind_small_cloud:":
		_, err = app.votingServiceClient.VoteSunBehindSmallCloud(ctx, voteRequest)
	case ":cloud_with_rain:":
		_, err = app.votingServiceClient.VoteCloudWithRain(ctx, voteRequest)
	case ":fire:":
		_, err = app.votingServiceClient.VoteFire(ctx, voteRequest)
	case ":jack_o_lantern:":
		_, err = app.votingServiceClient.VoteJackOLantern(ctx, voteRequest)
	case ":balloon:":
		_, err = app.votingServiceClient.VoteBalloon(ctx, voteRequest)
	case ":tada:":
		_, err = app.votingServiceClient.VoteTada(ctx, voteRequest)
	case ":trophy:":
		_, err = app.votingServiceClient.VoteTrophy(ctx, voteRequest)
	case ":iphone:":
		_, err = app.votingServiceClient.VoteIphone(ctx, voteRequest)
	case ":pager:":
		_, err = app.votingServiceClient.VotePager(ctx, voteRequest)
	case ":fax:":
		_, err = app.votingServiceClient.VoteFax(ctx, voteRequest)
	case ":bulb:":
		_, err = app.votingServiceClient.VoteBulb(ctx, voteRequest)
	case ":money_with_wings:":
		_, err = app.votingServiceClient.VoteMoneyWithWings(ctx, voteRequest)
	case ":crystal_ball:":
		_, err = app.votingServiceClient.VoteCrystalBall(ctx, voteRequest)
	case ":underage:":
		_, err = app.votingServiceClient.VoteUnderage(ctx, voteRequest)
	case ":interrobang:":
		_, err = app.votingServiceClient.VoteInterrobang(ctx, voteRequest)
	case ":100:":
		_, err = app.votingServiceClient.Vote100(ctx, voteRequest)
	case ":checkered_flag:":
		_, err = app.votingServiceClient.VoteCheckeredFlag(ctx, voteRequest)
	case ":crossed_swords:":
		_, err = app.votingServiceClient.VoteCrossedSwords(ctx, voteRequest)
	case ":floppy_disk:":
		_, err = app.votingServiceClient.VoteFloppyDisk(ctx, voteRequest)
	}
	if err != nil {
		writeError(err, w, r, http.StatusInternalServerError)
		return
	}
}

func (app *WebApp) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	indexTemplate := fmt.Sprintf(`
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
			<div id="motd" class="motd">%s</div>
			<div id="main" class="main"></div>
		</body>
		{{ if ne . ""}}
			<script type="text/javascript" src="{{ . }}/dist/index_bundle.js" async></script>
		{{else}}
			<script type="text/javascript" src="/js" async></script>
		{{end}}
	</html>`, app.messageOfTheDay)
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

func handle(path string, h func(w http.ResponseWriter, r *http.Request)) {
	http.Handle(path, otelhttp.WithRouteTag(path, http.HandlerFunc(h)))
}

func StartServer(webPort, webpackDevServer, indexBundle string, emojiServiceClient pb.EmojiServiceClient, votingClient pb.VotingServiceClient) {
	motd := os.Getenv("MESSAGE_OF_THE_DAY")
	webApp := &WebApp{
		emojiServiceClient:  emojiServiceClient,
		votingServiceClient: votingClient,
		indexBundle:         indexBundle,
		webpackDevServer:    webpackDevServer,
		messageOfTheDay:     motd,
	}

	log.Printf("Starting web server on WEB_PORT=[%s] and MESSAGE_OF_THE_DAY=[%s]", webPort, motd)
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

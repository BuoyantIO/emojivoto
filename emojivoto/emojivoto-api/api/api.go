package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	pb "github.com/runconduit/conduit-examples/emojivoto/emojivoto-api/gen/proto"
	"context"
	"strings"
	"errors"
)

type ApiApp struct {
	emojiClient pb.EmojiServiceClient
}

func (app *ApiApp) listEmojiHandler(w http.ResponseWriter, r *http.Request) {
	serviceResponse, err := app.emojiClient.ListAll(context.Background(), &pb.ListAllEmojiRequest{})
	if err != nil {
		writeError(err, w, r, http.StatusInternalServerError)
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

func (app *ApiApp) findEmojiHandler(w http.ResponseWriter, r *http.Request) {
	shortcodeToFind := strings.TrimPrefix(r.URL.Path, "/emoji/")

	serviceResponse, err := app.emojiClient.FindByShortcode(context.Background(), &pb.FindByShortcodeRequest{Shortcode: shortcodeToFind})
	if err != nil {
		writeError(err, w, r, http.StatusInternalServerError)
	}

	emojiReturned := serviceResponse.Emoji
	if emojiReturned == nil {
		err := errors.New("Not found")
		writeError(err, w, r, http.StatusNotFound)
	} else {
		writeJsonBody(w, http.StatusOK, map[string]string{
			"shortcode": emojiReturned.Shortcode,
			"unicode":   emojiReturned.Unicode,
		})
	}
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

func StartServer(webPort string, emojiclientClient pb.EmojiServiceClient) {
	webApp := &ApiApp{
		emojiClient: emojiclientClient,
	}

	log.Printf("Starting api server on API_PORT=[%s]", webPort)
	http.HandleFunc("/emoji", webApp.listEmojiHandler)
	http.HandleFunc("/emoji/", webApp.findEmojiHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%s", webPort), nil)
	if err != nil {
		panic(err)
	}
}

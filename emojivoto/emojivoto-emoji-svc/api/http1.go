package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/runconduit/conduit-examples/emojivoto/emojivoto-emoji-svc/emoji"
)

type EmojiH1Server struct {
	oldEmoji emoji.AllEmoji
	newEmoji emoji.AllEmoji
}

func findByShortcode(allEmoji emoji.AllEmoji, w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	foundEmoji := allEmoji.WithShortcode(params.ByName("shortcode"))
	if foundEmoji != nil {
		selectedEmoji := map[string]string{
			foundEmoji.Shortcode: foundEmoji.Unicode,
		}
		err := json.NewEncoder(w).Encode(selectedEmoji)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	} else {
		http.Error(w, "emoji not found", 500)
	}
}

func (s *EmojiH1Server) FindByShortcode2(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	findByShortcode(s.newEmoji, w, req, params)
}

func (s *EmojiH1Server) FindByShortcode(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	findByShortcode(s.oldEmoji, w, req, params)
}

func NewHTTP1Server(HTTP1Addr string, oldEmoji emoji.AllEmoji, newEmoji emoji.AllEmoji) *EmojiH1Server {
	router := httprouter.New()
	HTTP1Server := http.Server{
		Addr:         HTTP1Addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server := &EmojiH1Server{
		oldEmoji: oldEmoji,
		newEmoji: newEmoji,
	}

	router.GET("/find-by-shortcode2/:shortcode", server.FindByShortcode2)
	router.GET("/find-by-shortcode/:shortcode", server.FindByShortcode)

	err := HTTP1Server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	return server
}

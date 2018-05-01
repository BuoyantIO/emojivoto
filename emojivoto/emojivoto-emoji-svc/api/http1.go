package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/runconduit/conduit-examples/emojivoto/emojivoto-emoji-svc/emoji"
)

type EmojiH1Server struct {
	allEmoji emoji.AllEmoji
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

func (s *EmojiH1Server) FindByShortcode(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	foundEmoji := s.allEmoji.WithShortcode(params.ByName("shortcode"))
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

func NewHTTP1Server(HTTP1Addr string, allEmoji emoji.AllEmoji) *EmojiH1Server {
	router := httprouter.New()
	HTTP1Server := http.Server{
		Addr:         HTTP1Addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server := &EmojiH1Server{
		allEmoji: allEmoji,
	}

	router.GET("/find-by-shortcode/:shortcode", server.FindByShortcode)

	err := HTTP1Server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	return server
}

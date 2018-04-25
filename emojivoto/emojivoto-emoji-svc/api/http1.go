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

func (s *EmojiH1Server) FindByShortcodeNew(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	foundEmoji := s.newEmoji.WithShortcode(params.ByName("shortcode"))
	if foundEmoji != nil {
		selectedEmoji := map[string]string{
			foundEmoji.Shortcode: foundEmoji.Unicode,
		}
		err := json.NewEncoder(w).Encode(selectedEmoji)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
	}
}

func (s *EmojiH1Server) FindByShortcodeOld(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	foundEmoji := s.oldEmoji.WithShortcode(params.ByName("shortcode"))
	if foundEmoji != nil {
		selectedEmoji := map[string]string{
			foundEmoji.Shortcode: foundEmoji.Unicode,
		}
		json.NewEncoder(w).Encode(selectedEmoji)
	} else {
		http.Error(w, "emoji not found", 404)
	}
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

	router.GET("/FindByShortcodeNew/:shortcode", server.FindByShortcodeNew)
	router.GET("/FindByShortcodeOld/:shortcode", server.FindByShortcodeOld)

	err := HTTP1Server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	return server
}

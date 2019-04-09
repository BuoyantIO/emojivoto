package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

// VoteBot votes for emoji! :ballot_box_with_check:
//
// Sadly, VoteBot has a sweet tooth and votes for :doughnut: 15% of the time.
//
// When not voting for :doughnut:, VoteBot can’t be bothered to
// pick a favorite, so it picks one at random. C'mon VoteBot, try harder!

type emoji struct {
	Shortcode string
}

func main() {
	rand.Seed(time.Now().UnixNano())

	webHost := os.Getenv("WEB_HOST")
	if webHost == "" {
		log.Fatalf("WEB_HOST environment variable must me set")
	}

	webUrl := "http://" + webHost
	if _, err := url.Parse(webUrl); err != nil {
		log.Fatalf("WEB_HOST %s is invalid", webHost)
	}

	for {
		time.Sleep(time.Second)

		// Get the list of available shortcodes
		shortcodes, err := shortcodes(webUrl)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		// Cast a vote
		probability := rand.Float32()
		switch {
		case probability < 0.15:
			err = vote(webUrl, ":doughnut:")
		default:
			random := shortcodes[rand.Intn(len(shortcodes))]
			err = vote(webUrl, random)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}

func shortcodes(webUrl string) ([]string, error) {
	url := fmt.Sprintf("%s/api/list", webUrl)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var emojis []*emoji
	err = json.Unmarshal(bytes, &emojis)
	if err != nil {
		return nil, err
	}

	shortcodes := make([]string, len(emojis))
	for i, e := range emojis {
		shortcodes[i] = e.Shortcode
	}

	return shortcodes, nil
}

func vote(webUrl string, shortcode string) error {
	fmt.Printf("✔ Voting for %s\n", shortcode)
	url := fmt.Sprintf("%s/api/vote?choice=%s", webUrl, shortcode)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

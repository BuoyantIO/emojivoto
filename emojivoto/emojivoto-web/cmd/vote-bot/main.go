package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

// VoteBot votes for emoji! :ballot_box_with_check:
//
// Sadly, VoteBot has a sweet tooth and votes for :doughnut: 15% of the time.
// Furthermore, VoteBot is juvenile and votes for :poop: 20% of the time.
//
// When not voting for :doughnut: or :poop:, VoteBot can’t be bothered to
// pick a favorite, so it picks one at random. C'mon VoteBot, try harder!

type emoji struct {
	Shortcode string
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var sleep = flag.Duration("sleep", 100*time.Millisecond, "time to sleep between votes")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: emojivoto-vote-bot [options] <target>\n")
		fmt.Fprintf(os.Stderr, "       where <target> is host:port of web service, and [options] include:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	host := "http://" + flag.Arg(0)
	if _, err := url.Parse(host); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid target: %s\n", flag.Arg(0))
		flag.Usage()
		os.Exit(1)
	}

	for {
		time.Sleep(*sleep)

		// Get the list of available shortcodes
		shortcodes, err := shortcodes(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		// Cast a vote
		probability := rand.Float32()
		switch {
		case probability < 0.15:
			err = vote(host, ":doughnut:")
		case probability < 0.35:
			err = vote(host, ":poop:")
		default:
			random := shortcodes[rand.Intn(len(shortcodes))]
			err = vote(host, random)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		// Get the leaderboard
		err = leaderboard(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
	}
}

func shortcodes(host string) ([]string, error) {
	url := fmt.Sprintf("%s/api/list", host)

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

func vote(host string, shortcode string) error {
	fmt.Printf("✔ Voting for %s\n", shortcode)
	url := fmt.Sprintf("%s/api/vote?choice=%s", host, shortcode)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func leaderboard(host string) error {
	url := fmt.Sprintf("%s/api/leaderboard", host)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

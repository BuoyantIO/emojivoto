package integration_test

import (
	"os"
	"testing"
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/url"
)

var appBaseUrl = os.Getenv("WEB_URL")

func TestVotes(t *testing.T) {
	if appBaseUrl == "" {
		t.Fatalf("Please set env var WEB_URL, currently [%s]", appBaseUrl)
	}

	leaderboardList := getLeaderboard(t)

	if len(leaderboardList) != 0 {
		t.Fatalf("Expected leaderboard to be empty, got [%v]", leaderboardList)
	}

	allEligibleEmoji := getEmojilist(t)
	if len(allEligibleEmoji) < 10 {
		t.Fatal("Expected some emoji to be eligible, got 0")
	}

	someArbitratyIndex := len(allEligibleEmoji) - 1
	emojiToBeFirst := allEligibleEmoji[someArbitratyIndex-1]
	emojiToBeSecond := allEligibleEmoji[someArbitratyIndex-2]
	emojiToBeThird := allEligibleEmoji[someArbitratyIndex-3]

	postVoteFor(t, emojiToBeFirst)
	postVoteFor(t, emojiToBeSecond)
	postVoteFor(t, emojiToBeThird)

	postVoteFor(t, emojiToBeFirst)
	postVoteFor(t, emojiToBeSecond)
	postVoteFor(t, emojiToBeThird)

	postVoteFor(t, emojiToBeFirst)
	postVoteFor(t, emojiToBeSecond)

	postVoteFor(t, emojiToBeFirst)

	leaderboardList = getLeaderboard(t)

	if len(leaderboardList) != 3 {
		t.Fatalf("Expected leaderboard to have 3 emoji, got [%v]", leaderboardList)
	}

	firstItem := leaderboardList[0]
	if firstItem["shortcode"] != emojiToBeFirst["shortcode"] {
		t.Fatalf("Expected 2st emoji in list to be [%v], was [%v]", emojiToBeFirst, firstItem)
	}

	secondItem := leaderboardList[1]
	if secondItem["shortcode"] != emojiToBeSecond["shortcode"] {
		t.Fatalf("Expected 2st emoji in list to be [%v], was [%v]", emojiToBeSecond, secondItem)
	}

	thirdItem := leaderboardList[2]
	if secondItem["shortcode"] != emojiToBeSecond["shortcode"] {
		t.Fatalf("Expected 3rd emoji in list to be [%v], was [%v]", emojiToBeThird, thirdItem)
	}
}

func postVoteFor(t *testing.T, emojiToVoteFor map[string]string) {

	http.Post(fmt.Sprintf("%s/api/vote", appBaseUrl), "", nil)
	resp, err := http.PostForm(fmt.Sprintf("%s/api/vote", appBaseUrl),
		url.Values{"choice": {emojiToVoteFor["shortcode"]}})

	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}
	if status := resp.StatusCode; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func getLeaderboard(t *testing.T) []map[string]string {
	return getJsonMapLisp(fmt.Sprintf("%s/api/leaderboard", appBaseUrl), t)
}

func getEmojilist(t *testing.T) []map[string]string {
	return getJsonMapLisp(fmt.Sprintf("%s/api/list", appBaseUrl), t)
}

func getJsonMapLisp(url string, t *testing.T) []map[string]string {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if status := resp.StatusCode; status != http.StatusOK {
		t.Fatalf("get to [%s] returned wrong status code: got %v want %v",
			url, status, http.StatusOK)
	}
	var list []map[string]string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(body, &list); err != nil {
		t.Fatalf("Error parsing [%v]: %s",err, string(body))
	}
	return list
}

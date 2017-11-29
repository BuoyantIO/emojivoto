package integration_test

import (
	"os"
	"testing"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

var apiBaseUrl = os.Getenv("API_URL")

func TestEmoji(t *testing.T) {
	if apiBaseUrl == "" {
		t.Fatalf("Please set env var API_URL, currently [%s]", apiBaseUrl)
	}

	emojiList := getList(t)

	if len(emojiList) < 1 {
		t.Fatalf("Expected emoji list to be full but it was empty")
	}

	for _, e := range emojiList {
		emoji := getEmoji(t, e["shortcode"])
		if emoji["shortcode"] != e["shortcode"] || emoji["unicode"] != e["unicode"] {
			t.Fatalf("Expected to find [%s], got [%s]", e, emoji)
		}
	}

	targetUrl := fmt.Sprintf("%s/emoji/does-not-exist", apiBaseUrl)
	resp, err := http.Get(targetUrl)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusNotFound{
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func getList(t *testing.T) []map[string]string {
	targetUrl := fmt.Sprintf("%s/emoji", apiBaseUrl)
	resp, err := http.Get(targetUrl)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
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

func getEmoji(t *testing.T, shortcode string) map[string]string {
	targetUrl := fmt.Sprintf("%s/emoji/%s", apiBaseUrl, shortcode)
	resp, err := http.Get(targetUrl)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var emoji map[string]string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(body, &emoji); err != nil {
		t.Fatalf("Error parsing [%v]: %s",err, string(body))
	}
	return emoji
}

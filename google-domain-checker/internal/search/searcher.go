package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Searcher interface {
	Search(request, location string) Result
}

type Result struct {
	OrganicResults []struct {
		URL string `json:"url"`
	} `json:"organic_results"`
}

type googleSearcher struct {
	serpStackAPIKey string
}

func NewGoogleSearcher(APIKey string) Searcher {
	return googleSearcher{serpStackAPIKey: APIKey}
}

func (g googleSearcher) Search(request, location string) Result {
	var result Result

	resp, _ := http.Get(fmt.Sprintf(
		"http://api.serpstack.com/search?access_key=%s&query=%s&location=%s",
		g.serpStackAPIKey, url.QueryEscape(request), url.QueryEscape(location),
	))
	defer resp.Body.Close()

	_ = json.NewDecoder(resp.Body).Decode(&result)

	return result
}

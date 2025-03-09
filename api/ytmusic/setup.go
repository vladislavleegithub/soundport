package ytmusic

import (
	"net/http"
)

const (
	YTMUSIC_BASE_URL = "https://music.youtube.com"
	YTMUSIC_SEARCH   = YTMUSIC_BASE_URL + "/youtubei/v1/search?limit=1"
	YTMUSIC_PLAYLIST = YTMUSIC_BASE_URL + "/youtubei/v1/playlist/create"
	PARAM            = "EgWKAQIIAWoQEAMQBBAJEAoQBRAREBAQFQ%3D%3D"
)

const (
	PRIVATE StatusType = "PRIVATE"
	PUBLIC  StatusType = "PUBLIC"
)

type SearchRequestBody struct {
	Ctx    *Context `json:"context"`
	Query  string   `json:"query"`
	Params string   `json:"params"`
}

type CreatePlaylistRequestBody struct {
	Ctx           *Context   `json:"context"`
	Title         string     `json:"title"`
	PrivacyStatus StatusType `json:"privacyStatus"`
	VideoIds      []string   `json:"videoIds"`
}

type Context struct {
	Client struct {
		Hl            string `json:"hl"`
		Gl            string `json:"gl"`
		ClientName    string `json:"client_name"`
		ClientVersion string `json:"client_version"`
	} `json:"client"`
}

func initHeaders(r *http.Request) {
	r.Header.Add(
		"user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0",
	)
	r.Header.Add("accept", "*/*")
	r.Header.Add("content-type", "application/json")
	r.Header.Add("X-origin", YTMUSIC_BASE_URL)
	r.Header.Add("Origin", YTMUSIC_BASE_URL)
}

func sendGetRequest() (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", YTMUSIC_BASE_URL, nil)
	if err != nil {
		return nil, err
	}

	initHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

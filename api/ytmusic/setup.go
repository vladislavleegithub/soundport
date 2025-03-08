package ytmusic

import (
	"net/http"
)

const YTMUSIC_BASE_URL = "https://music.youtube.com"

type RequstBody struct {
	Context `       json:"context"`
	Query   string `json:"query"`
}

type BasePayload struct {
	Body RequstBody
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
	r.Header.Add("origin", YTMUSIC_BASE_URL)
}

func sendGetRequest() (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", YTMUSIC_BASE_URL, nil)
	if err != nil {
		return nil, err
	}

	initHeaders(req)

	// Add additional headers
	req.Header.Add("content-encoding", "gzip")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

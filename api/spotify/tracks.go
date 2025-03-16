package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type playlistTracks struct {
	Tracks []struct {
		Album struct {
			Name string `json:"name"`
		} `json:"album"`
		Name string `json:"name"`
	} `json:"tracks"`
}

func (a *auth) GetTracks(tracksUrl string) (playlistTracks, error) {
	// setup auth herader
	authHeader := fmt.Sprintf("Bearer %s", a.accessToken)

	// prep request
	req, err := http.NewRequest("GET", tracksUrl, nil)
	if err != nil {
		return playlistTracks{}, nil
	}
	req.Header.Add("Authorization", authHeader)

	// send response
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return playlistTracks{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return playlistTracks{}, errors.New("unable to fetch songs")
	}

	// decode response
	tracks := playlistTracks{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tracks)
	if err != nil {
		return playlistTracks{}, err
	}

	return tracks, nil
}

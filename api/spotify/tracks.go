package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type PlaylistTracks struct {
	Total  int `json:"total"`
	Tracks []struct {
		Track struct {
			Album struct {
				Name string `json:"name"`
			} `json:"album"`
			Name string `json:"name"`
		} `json:"track"`
	} `json:"items"`
}

func (a *auth) GetTracks(tracksUrl string) (PlaylistTracks, error) {
	// setup auth herader
	authHeader := fmt.Sprintf("Bearer %s", a.accessToken)

	// prep request
	req, err := http.NewRequest("GET", tracksUrl, nil)
	if err != nil {
		return PlaylistTracks{}, nil
	}
	req.Header.Add("Authorization", authHeader)

	// send response
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return PlaylistTracks{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PlaylistTracks{}, errors.New("unable to fetch songs")
	}

	// decode response
	tracks := PlaylistTracks{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tracks)
	if err != nil {
		return PlaylistTracks{}, err
	}

	return tracks, nil
}

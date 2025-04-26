package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Track struct {
	Album struct {
		Name string `json:"name"`
	} `json:"album"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Name string `json:"name"`
}

type PlaylistTracks struct {
	Next   string `json:"next"`
	Total  int    `json:"total"`
	Tracks []struct {
		Track Track `json:"track"`
	} `json:"items"`
}

func (a *auth) GetPlaylistTracks(plId string) ([]string, error) {
	var finalSongs []string
	authHeader := fmt.Sprintf("Bearer %s", a.accessToken)

	for {
		if plId == "" {
			break
		}

		// prep request
		req, err := http.NewRequest("GET", plId, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", authHeader)

		// send response
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("unable to fetch songs")
		}

		// decode response
		tracks := PlaylistTracks{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&tracks)
		if err != nil {
			return nil, err
		}

		trackNames := getTracksString(tracks)
		finalSongs = append(finalSongs, trackNames...)

		plId = tracks.Next
	}
	return finalSongs, nil
}

func getTracksString(pl PlaylistTracks) []string {
	var songs []string
	for _, val := range pl.Tracks {
		songName := fmt.Sprintf("%s By %s", val.Track.Name, val.Track.Artists[0].Name)
		songs = append(songs, songName)
	}

	return songs
}

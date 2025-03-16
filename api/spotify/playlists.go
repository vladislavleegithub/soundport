package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type spfyPlaylists struct {
	Total     int `json:"total"`
	Playlists []struct {
		Desc   string `json:"description"`
		Link   string `json:"href"`
		Id     string `json:"id"`
		Name   string `json:"name"`
		Tracks struct {
			Link  string `json:"href"`
			Total int    `json:"total"`
		} `json:"tracks"`
	} `json:"items"`
}

func (a *auth) GetPlaylists() (spfyPlaylists, error) {
	// setup auth herader
	authHeader := fmt.Sprintf("Bearer %s", a.accessToken)

	// Prep request
	req, err := http.NewRequest("GET", playlist_url, nil)
	if err != nil {
		return spfyPlaylists{}, err
	}
	req.Header.Add("Authorization", authHeader)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return spfyPlaylists{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("ERROR: ", resp.Status)
		return spfyPlaylists{}, errors.New("unsuccessful request")
	}

	playlists := spfyPlaylists{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&playlists)
	if err != nil {
		return spfyPlaylists{}, err
	}

	return playlists, nil
}

func (a *auth) RefreshSession() error {
	// If auth token is close to expiry, we still refresh it.
	checkTime := a.expiresAt.Add(-5 * time.Minute)

	if time.Now().Before(checkTime) {
		return nil
	}

	// Set up request body
	body := url.Values{}
	body.Add("grant_type", "refresh_token")
	body.Add("refresh_token", a.refreshToken)
	encodedBody := strings.NewReader(body.Encode())

	req, err := http.NewRequest("POST", token_url, encodedBody)
	if err != nil {
		return err
	}

	authHeader := a.creds.getAuthorizationHeader()
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err = handleAuthResponse(req)
	if err != nil {
		return err
	}

	return nil
}

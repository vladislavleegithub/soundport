package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Samarthbhat52/soundport/api"
	"github.com/charmbracelet/bubbles/list"
)

type Playlist struct {
	Desc   string `json:"description"`
	Name   string `json:"name"`
	Tracks struct {
		Link  string `json:"href"`
		Total int    `json:"total"`
	} `json:"tracks"`
}

// Make it compatable with bubbletea
func (p Playlist) FilterValue() string { return p.Name }
func (p Playlist) Title() string       { return p.Name }
func (p Playlist) Description() string { return fmt.Sprintf("Num tracks: %d", p.Tracks.Total) }
func (p Playlist) GetPlaylistDetails() *api.PlaylistDetails {
	selected := api.PlaylistDetails{}

	selected.PlId = p.Tracks.Link
	selected.PlName = p.Name
	selected.PlDesc = p.Desc
	selected.TotalTracks = p.Tracks.Total

	return &selected
}

// Api return struct
type spfyPlaylists struct {
	Total         int        `json:"total"`
	ItemPlaylists []Playlist `json:"items"`
}

func (s *spfyPlaylists) GetPlaylistItems() []list.Item {
	items := make([]list.Item, len(s.ItemPlaylists))

	for i := range len(s.ItemPlaylists) {
		items[i] = s.ItemPlaylists[i]
	}

	return items
}

func (a *auth) GetPlaylists() ([]list.Item, error) {
	// setup auth herader
	authHeader := fmt.Sprintf("Bearer %s", a.accessToken)

	// Prep request
	req, err := http.NewRequest("GET", playlist_url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("ERROR: ", resp.Status)
		return nil, errors.New("unsuccessful request")
	}

	playlists := spfyPlaylists{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&playlists)
	if err != nil {
		return nil, err
	}

	playlistItems := playlists.GetPlaylistItems()

	return playlistItems, nil
}

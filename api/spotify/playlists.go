package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/bubbles/list"
	"github.com/vladislavleegithub/soundport/api"
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
	Next          string     `json:"next"` // URL следующей страницы (если есть)
}

func (s *spfyPlaylists) GetPlaylistItems() []list.Item {
	items := make([]list.Item, len(s.ItemPlaylists))

	for i := range len(s.ItemPlaylists) {
		items[i] = s.ItemPlaylists[i]
	}

	return items
}

func (a *auth) GetPlaylists() ([]list.Item, error) {
	var allPlaylists []Playlist
	url := playlist_url // Начинаем с базового URL

	for {
		// 1. Формируем запрос
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.accessToken))

		// 2. Отправляем запрос
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API error: %s", resp.Status)
		}

		// 3. Парсим ответ
		var response spfyPlaylists
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}

		// 4. Добавляем плейлисты в общий список
		allPlaylists = append(allPlaylists, response.ItemPlaylists...)

		// 5. Если следующей страницы нет — выходим
		if response.Next == "" {
			break
		}

		// 6. Иначе переходим на следующую страницу
		url = response.Next
	}

	// 7. Конвертируем в формат для bubbletea
	items := make([]list.Item, len(allPlaylists))
	for i, pl := range allPlaylists {
		items[i] = pl
	}

	return items, nil
}

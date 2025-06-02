package api

import (
	"github.com/charmbracelet/bubbles/list"
)

type PlaylistDetails struct {
	PlId        string
	PlName      string
	PlDesc      string
	TotalTracks int
}

// To adhere to bubbble tea's list interface
func (p PlaylistDetails) FilterValue() string { return p.PlName }

type SongDetails struct {
	Name  string
	Id    string
	Found bool
}

type Playlist interface {
	list.DefaultItem
	GetPlaylistDetails() *PlaylistDetails
}

type Source interface {
	GetPlaylists() ([]list.Item, error)
	GetPlaylistTracks(string) ([]string, error)
}

type Destination interface {
	CreatePlaylist(name string, desc string) (string, error)
	AddTracks(plId string, tracks []string) (int, bool)
}

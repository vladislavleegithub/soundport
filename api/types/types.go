package types

import (
	"github.com/charmbracelet/bubbles/list"
)

type PlaylistDetails struct {
	PlId        string
	PlName      string
	PlDesc      string
	TotalTracks int
}

func (p PlaylistDetails) FilterValue() string { return p.PlName }

type SongDetails struct {
	Name  string
	Id    string
	Found bool
}

type Source interface {
	GetPlaylists() ([]list.Item, error)
	GetPlaylistTracks(string) ([]string, error)
}

type Playlist interface {
	list.DefaultItem
	GetPlaylistDetails() *PlaylistDetails
}

type Destination interface {
	FindTracks([]string, chan<- SongDetails)
	AddPlaylist(string, []string) error
}

type SetPlaylist interface {
	GetPlaylistDetails() *PlaylistDetails
}

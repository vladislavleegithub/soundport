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

type ensureAuth interface {
	EnsureInit()
	EnsureLogin()
}

type Source interface {
	// ensureAuth
	GetPlaylists() ([]list.Item, error)
	GetPlaylistTracks(string) ([]string, error)
}

type Playlist interface {
	list.DefaultItem
	GetPlaylistDetails() *PlaylistDetails
}

type Destination interface {
	ensureAuth
	CreatePlaylist(name string, desc string) (string, error)
	AddTracks(plId string, tracks []string) bool
}

type SetPlaylist interface {
	GetPlaylistDetails() *PlaylistDetails
}

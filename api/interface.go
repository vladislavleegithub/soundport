package api

import (
	"github.com/charmbracelet/bubbles/list"
)

type Playlist interface {
	list.DefaultItem
	GetPlaylistDetails() *PlaylistDetails
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

type Destination interface {
	ensureAuth
	CreatePlaylist(name string, desc string) (string, error)
	AddTracks(plId string, tracks []string) bool
}

package types

import "github.com/charmbracelet/bubbles/list"

/*
The main port command:
	1. Should take in destination.
	2. Fetch source playlists.
	3. Select a playlist.
	4. Pass the selected playlist to screen two.

Screen 2:
	1. accepts a playlist interface.
	2. fetches the available tracks.
	3. Does a batch submit.

playlist interface:
	1. GetTracks() []string
	2. CreatePlaylist(songs []string) bool
	2. <-- More methods will be added later -->

*/

type Playlists interface {
	GetPlaylists() []list.DefaultItem
}

package ui

import (
	"fmt"

	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/Samarthbhat52/soundport/api/ytmusic"
	tea "github.com/charmbracelet/bubbletea"
)

func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spotify.PlaylistTracks:
		songs := constructSongsList(msg)
		m.foundSongs = make(chan string, len(songs))

		return m, tea.Batch(waitForActivity(m.foundSongs), listenForActivity(m.foundSongs, songs))

	case responseMessage:
		m.foundSongsCount++
		return m, waitForActivity(m.foundSongs)

	default:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}
}

func viewChosen(m model) string {
	s := fmt.Sprintf(
		"\n %s Songs found: %d\n\nPress `q` to exit\n",
		m.Spinner.View(),
		m.foundSongsCount,
	)
	if m.Quitting {
		s += "\n"
	}

	return s
}

func constructSongsList(track spotify.PlaylistTracks) []string {
	songsList := []string{}

	for _, val := range track.Tracks {
		sn := val.Track.Name
		sn += " " + val.Track.Artists[0].Name
		sn += " " + val.Track.Album.Name
		songsList = append(songsList, sn)
	}

	return songsList
}

type responseMessage string

func waitForActivity(sub chan string) tea.Cmd {
	return func() tea.Msg {
		return responseMessage(<-sub)
	}
}

func listenForActivity(sub chan string, songs []string) tea.Cmd {
	return func() tea.Msg {
		for _, sn := range songs {
			go ytmusic.GetSongInfo(sn, sub)
		}

		return tea.Quit
	}
}

package ui

import (
	"fmt"

	"github.com/Samarthbhat52/soundport/api/spotify"
	tea "github.com/charmbracelet/bubbletea"
)

func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.List.SelectedItem().(spotify.Playlists)
			if ok {
				m.Choice = i
			}

			cmd := getSongsCmd(m)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func viewChoices(m model) string {
	return "\n" + m.List.View()
}

func getSongsCmd(m model) tea.Cmd {
	return func() tea.Msg {
		a, _ := spotify.NewAuth()
		songs, err := a.GetTracks(m.Choice.Tracks.Link)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return nil
		}

		return spotify.PlaylistTracks(songs)
	}
}

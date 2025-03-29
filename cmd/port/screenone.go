package port

import (
	"fmt"

	"github.com/Samarthbhat52/soundport/api/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func plSelected() tea.Cmd {
	return func() tea.Msg {
		return playlistSelected{}
	}
}

func (m *portModel) updatePlaylists(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.playlists.SetSize(msg.Width-h, msg.Height/2-v)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item, ok := m.playlists.SelectedItem().(types.Playlist)
			if !ok {
				fmt.Println("Improper types in selected item")
				// TODO: Add proper error handling
				return m, tea.Quit
			}

			selected := item.GetPlaylistDetails()
			m.ch = make(chan types.SongDetails, selected.TotalTracks)
			m.selected = selected

			return m, plSelected()
		}
	}
	var cmd tea.Cmd
	m.playlists, cmd = m.playlists.Update(msg)
	return m, cmd
}

func (m *portModel) viewPlaylists() string {
	return docStyle.Render(m.playlists.View())
}

package port

import (
	"github.com/Samarthbhat52/soundport/api"
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
			item, ok := m.playlists.SelectedItem().(api.Playlist)
			if !ok {
				m.quitting = true
				glbLogger.Println("Selected playlist does not implement the `Playlist` interface.")
				return m, tea.Quit
			}

			selected := item.GetPlaylistDetails()
			plId, err := m.dst.CreatePlaylist(selected.PlName, selected.PlDesc)
			if err != nil {
				glbLogger.Println("Error creating playlist: ", err.Error())
				m.quitting = true
				return m, tea.Quit
			}

			m.createdPlId = plId
			m.selectedPlId = selected.PlId

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

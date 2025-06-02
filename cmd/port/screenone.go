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
		case "enter": // Select a playlist from source
			item, ok := m.playlists.SelectedItem().(api.Playlist)
			if !ok {
				glbLogger.Println("Selected playlist does not implement the `Playlist` interface.")
				m.quitting = true
				return m, tea.Quit
			}

			selected := item.GetPlaylistDetails()

			// Create a destination playlist first which we will add to later.
			plId, err := m.dst.CreatePlaylist(selected.PlName, selected.PlDesc)
			if err != nil {
				glbLogger.Println("Error creating playlist: ", err.Error())
				m.quitting = true
				return m, tea.Quit
			}

			m.createdPlId = plId
			// Will be used to pull out songs belonging to the source playlist.
			m.selectedPlId = selected.PlId

			return m, plSelected()
		}
	}

	var cmd tea.Cmd
	m.playlists, cmd = m.playlists.Update(msg)
	return m, cmd
}

func (m *portModel) viewPlaylists() string {
	// Render all the playlists to be selected from,
	// in a bubble tea default list view.
	return docStyle.Render(m.playlists.View())
}

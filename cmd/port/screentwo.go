package port

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vladislavleegithub/soundport/ui"
)

func (m *portModel) updatePortProgress(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case playlistSelected:
		var message strings.Builder

		// Get all the tracks present in source playlist.
		tracks, err := m.src.GetPlaylistTracks(m.selectedPlId)
		if err != nil {
			glbLogger.Println("err: ", err.Error())
			m.quitting = true
			return m, tea.Quit
		}

		if addedTracks, ok := m.dst.AddTracks(m.createdPlId, tracks); ok {
			message.WriteString(fmt.Sprintf("Added %d out of %d total songs.", addedTracks, len(tracks)))

			if addedTracks != len(tracks) {
				message.WriteString("\nRun " + ui.Accent.Render("cat '/tmp/sp_notfound.log'") + " to get the list of songs that failed to port.")
			}

			m.statusMessage = message.String()
			m.successful = true
			m.quitting = true
			// Add a something that will show that playlist successfully created.
			return m, tea.Quit
		}

		// Cat the log files to know what happened here
		glbLogger.Println("unable to add tracks")
		m.quitting = true
		return m, tea.Quit

	default:
		return m, nil
	}
}

func (m *portModel) viewPortProgress() string {
	var s string
	s += "Porting playlist....\nFinding songs and adding to the new playlist..."
	s += "\n\n(Press 'q' or 'ctrl+c' to quit)"

	if m.successful {
		s = m.statusMessage
	}

	if m.quitting {
		s += "\n"
	}

	return s
}

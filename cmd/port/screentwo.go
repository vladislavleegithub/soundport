package port

import (
	"strings"

	"github.com/Samarthbhat52/soundport/api/types"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	padding  = 2
	maxWidth = 80
)

func plCreated() tea.Cmd {
	return func() tea.Msg {
		return playlistCreated{}
	}
}

func (m *portModel) updatePortProgress(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = min(msg.Width-padding*2-4, maxWidth)

		return m, nil

	case playlistSelected:
		// Get all the tracks present in source playlist.
		tracks, err := m.src.GetPlaylistTracks(m.selected.PlId)
		if err != nil {
			glbLogger.Println("err: ", err.Error())
			return m, tea.Quit
		}

		return m, tea.Batch(
			m.waitForActivity(),
			// Get equivalent track information from destination provider.
			m.listenForActivity(tracks),
		)

	case types.SongDetails:
		m.total += 1
		m.percent = float64(m.total) / float64(m.selected.TotalTracks)

		if msg.Found {
			m.songs = append(m.songs, msg.Id)
		} else {
			glbLogger.Printf("Not found: %s, ID: %s\n", msg.Name, msg.Id)
		}

		return m, m.waitForActivity()

	case songSearchComplete:
		m.dst.AddPlaylist(m.selected.PlName, m.songs)
		return m, plCreated()

	case playlistCreated:
		m.quitting = true
		return m, tea.Quit

	default:
		return m, nil
	}
}

func (m *portModel) viewPortProgress() string {
	s := "\n"
	s += strings.Repeat(" ", padding)

	s += "Finding songs..."
	s += m.progress.ViewAs(m.percent) + "\n\n"

	if m.quitting {
		s += "\n"
	}

	return s
}

func (m *portModel) waitForActivity() tea.Cmd {
	return func() tea.Msg {
		if val, ok := <-m.ch; ok {
			return val
		}

		return songSearchComplete(true)
	}
}

func (m *portModel) listenForActivity(songs []string) tea.Cmd {
	return func() tea.Msg {
		m.dst.FindTracks(songs, m.ch)
		return nil
	}
}

package ui

import (
	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Two views.
// 1. List playlists.
// 2. Show Number of songs in playlist.

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

var (
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	progressEmpty = subtleStyle.Render(progressEmptyChar)
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)
)

type model struct {
	List     list.Model
	Quitting bool
	Choice   spotify.Playlists
	Progress float64
	Spinner  spinner.Model
	SongLen  int
	msg      string
}

func InitModel(list list.Model) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		List:    list,
		Spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()

		if k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	if m.Choice == (spotify.Playlists{}) {
		return updateChoices(msg, m)
	}

	return updateChosen(msg, m)
}

func (m model) View() string {
	var s string

	if m.Choice == (spotify.Playlists{}) {
		s = viewChoices(m)
	} else {
		s = viewChosen(m)
	}

	return mainStyle.Render("\n" + s + "\n\n")
}

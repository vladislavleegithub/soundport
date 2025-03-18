package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case songsList:
		m.msg = fmt.Sprintf("TOTAL : %v\n", msg.Total)
		m.SongLen = msg.Total
		m.Quitting = true

		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}
}

func viewChosen(m model) string {
	if m.Quitting {
		return fmt.Sprintf("Songs found: %v\n", m.msg)
	}

	return fmt.Sprintf("\n\n %s Finding songs...press q to quit\n", m.Spinner.View())
}

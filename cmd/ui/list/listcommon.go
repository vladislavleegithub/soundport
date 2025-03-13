package listcommon

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	platlists []string
	cursor    int
	selected  int
}

func InitialModel(playlists []string) model {
	return model{
		platlists: playlists,
	}
}

func (m model) Init() tea.Cmd {
	return nil // Add i/o later.
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl + c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.platlists)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = m.cursor
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Select a playlist\n\n"

	for i, pl := range m.platlists {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if i == m.selected {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, pl)
	}

	s += "\nPress q to quit.\n"

	return s
}

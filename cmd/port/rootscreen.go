package port

import (
	"github.com/Samarthbhat52/soundport/logger"
	tea "github.com/charmbracelet/bubbletea"
)

var glbLogger = logger.GetInstance()

type rootScreenModal struct {
	model tea.Model
}

func rootScreen() *rootScreenModal {
	return &rootScreenModal{
		model: screenOne(),
	}
}

func (m rootScreenModal) Init() tea.Cmd {
	return m.model.Init()
}

func (m rootScreenModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m rootScreenModal) View() string {
	return m.model.View()
}

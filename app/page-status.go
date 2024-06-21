package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusPage struct {
}

func (p *StatusPage) Update(_ *Model, msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q", "esc", "enter", "space":
		return tea.Quit
	default:
		return tea.Quit
	}
}

func (p *StatusPage) Render(m *Model) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"STATUS",
			m.style.StatusStyle.Render(m.err),
		),
	)
}

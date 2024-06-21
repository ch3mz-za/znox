package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ErrorPage struct {
}

func (p *ErrorPage) Render(m *Model) string {
	return renderPageError(m)
}

func (p *ErrorPage) Update(m *Model, msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return tea.Quit
	case "enter", "space":
		m.currentPage = m.lastPage
		return nil
	}
	return cmd
}

type FatalErrorPage struct {
}

func (p *FatalErrorPage) Update(_ *Model, msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q", "esc", "enter", "space":
		return tea.Quit
	default:
		return tea.Quit
	}
}

func (p *FatalErrorPage) Render(m *Model) string {
	return renderPageError(m)
}

func renderPageError(m *Model) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"ERROR",
			m.style.ErrorStyle.Render(m.err),
		),
	)
}

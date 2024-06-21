package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EncryptionPage struct {
	list list.Model
}

func (p *EncryptionPage) Update(m *Model, msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return tea.Quit

	case "enter":
		i, ok := p.list.SelectedItem().(item)
		if ok {
			m.choice = string(i)
			m.currentPage = m.pages[passwordPage]
		}

		if m.choice == actionDecrypt && strings.HasSuffix(m.sourceFile, ".enc") {
			m.destinationDir = strings.TrimSuffix(m.destinationDir, ".enc")
		} else {
			m.destinationDir += ".enc"
		}
	}
	p.list, cmd = p.list.Update(msg)
	return cmd
}

func (p *EncryptionPage) Render(m *Model) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		strings.Join([]string{
			lipgloss.JoinVertical(
				lipgloss.Center,
				"Znox v3.0.0",
			),
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.style.ListStyle.Render("\n"+p.list.View()),
			),
		}, "\n"),
	)
}

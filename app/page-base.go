package app

import tea "github.com/charmbracelet/bubbletea"

type PageID int

const (
	encryptionPage PageID = iota
	passwordPage
	statusPage
	errorPage
	fatalErrorPage
)

type Page interface {
	Render(m *Model) string
	Update(m *Model, msg tea.KeyMsg) tea.Cmd
}

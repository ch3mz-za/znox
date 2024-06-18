package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Page int

const (
	PasswordPage Page = iota
	OptionsPage
	ErrorPage
	FatalErrorPage
	StatusPage
)

func renderPageOptions(m *Model) string {
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
				m.style.ListStyle.Render("\n"+m.list.View()),
			),
		}, "\n"),
	)
}

func renderPagePasswords(m *Model) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"Enter password:",
			m.style.InputField.Render(m.passwField.View()),
		),
	)
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

func renderPageStatus(m *Model) string {
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

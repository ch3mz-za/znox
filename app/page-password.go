package app

import (
	"fmt"

	"github.com/ch3mz-za/znox/znox"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PasswordPage struct {
	passwordField textinput.Model
	passwords     [2]string
}

func (p *PasswordPage) Update(m *Model, msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return tea.Quit
	case "enter":

		if m.choice == actionEncrypt {
			if p.passwords[0] == "" {
				// TODO: Display error if validation fails
				p.passwords[0] = p.passwordField.Value()
				p.passwordField.SetValue("")
			} else if p.passwords[1] == "" {
				// TODO: Display error if validation fails
				p.passwords[1] = p.passwordField.Value()

				if p.passwords[0] == p.passwords[1] {

					password := []byte(p.passwordField.Value())
					if err := znox.Encryption(m.sourceFile, m.destinationDir, password, password); err != nil {
						setError(m, err.Error())
						return cmd
					}

					m.err = fmt.Sprintf("Successfully %sed!", m.choice)
					m.currentPage = m.pages[statusPage]
					return cmd

				} else {
					setError(m, "Passwords do not match!")
					p.passwords[0], p.passwords[1] = "", ""
					p.passwordField.SetValue("")
				}
			}

		} else {
			if err := znox.Decryption(m.sourceFile, m.destinationDir, []byte(p.passwordField.Value())); err != nil {
				setError(m, err.Error())
				return cmd
			}
			m.err = fmt.Sprintf("Successfully %sed!", m.choice)
			m.currentPage = m.pages[statusPage]
		}

	}

	p.passwordField, cmd = p.passwordField.Update(msg)
	return cmd
}

func (p *PasswordPage) Render(m *Model) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"Enter password:",
			m.style.InputField.Render(p.passwordField.View()),
		),
	)
}

package app

import (
	"fmt"
	"strings"

	"github.com/ch3mz-za/znox/znox"
	tea "github.com/charmbracelet/bubbletea"
)

func updatePageOptions(m *Model, msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return tea.Quit

	case "enter":
		i, ok := m.list.SelectedItem().(item)
		if ok {
			m.choice = string(i)
			m.currentPage = PasswordPage
		}

		if m.choice == actionDecrypt && strings.HasSuffix(m.sourceFile, ".enc") {
			m.destinationDir = strings.TrimSuffix(m.destinationDir, ".enc")
		} else {
			m.destinationDir += ".enc"
		}
	}
	m.list, cmd = m.list.Update(msg)
	return cmd
}

func updatePagePassword(m *Model, msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return tea.Quit
	case "enter":

		if m.choice == actionEncrypt {
			if m.passwords[0] == "" {
				// TODO: Display error if validation fails
				m.passwords[0] = m.passwField.Value()
				m.passwField.SetValue("")
			} else if m.passwords[1] == "" {
				// TODO: Display error if validation fails
				m.passwords[1] = m.passwField.Value()

				if m.passwords[0] == m.passwords[1] {

					password := []byte(m.passwField.Value())
					if err := znox.Encryption(m.sourceFile, m.destinationDir, password, password); err != nil {
						setError(m, err.Error())
						return cmd
					}

					m.err = fmt.Sprintf("Successfully %sed!", m.choice)
					m.currentPage = StatusPage

				} else {
					setError(m, "Passwords do not match!")
					m.passwords[0], m.passwords[1] = "", ""
					m.passwField.SetValue("")
				}
			}

		} else {
			if err := znox.Decryption(m.sourceFile, m.destinationDir, []byte(m.passwField.Value())); err != nil {
				setError(m, err.Error())
				return cmd
			}
			m.err = fmt.Sprintf("Successfully %sed!", m.choice)
			m.currentPage = StatusPage
		}

	}

	m.passwField, cmd = m.passwField.Update(msg)
	return cmd
}

func updatePageError(m *Model, msg tea.KeyMsg) tea.Cmd {
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

func updatePageFatalError(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q", "esc", "enter", "space":
		return tea.Quit
	default:
		return tea.Quit
	}
}

func updatePageStatus(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q", "esc", "enter", "space":
		return tea.Quit
	default:
		return tea.Quit
	}
}

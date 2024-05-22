package app

import (
	"fmt"

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

		if len(m.passwords) < 2 {
			// TODO: Validate password before appending
			m.passwords = append(m.passwords, m.passwField.Value())
			m.passwField.SetValue("")

		} else if len(m.passwords) > 1 {
			if m.passwords[0] == m.passwords[1] {

				switch m.choice {
				case actionDecrypt:
					if err := znox.Decryption(m.srcFile, m.dstDir, []byte(m.passwField.Value())); err != nil {
						m.err = err.Error()
						m.currentPage = ErrorPage
						return cmd
					}
				case actionEncrypt:
					password := []byte(m.passwField.Value())
					if err := znox.Encryption(m.srcFile, m.dstDir, password, password); err != nil {
						m.err = err.Error()
						m.currentPage = ErrorPage
						return cmd
					}
				}

				m.err = fmt.Sprintf("Successfully %sed!", m.choice)
				m.currentPage = StatusPage
			}
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

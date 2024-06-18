package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultWidth = 20

	actionEncrypt = "Encrypt"
	actionDecrypt = "Decrypt"
)

type Model struct {
	width  int
	height int

	sourceFile     string
	destinationDir string
	choice         string
	err            string

	currentPage Page
	lastPage    Page
	style       *Style
	passwField  textinput.Model
	passwords   [2]string

	list list.Model
}

func New(src, dst string) *Model {
	currentPage := OptionsPage
	s := DefaultStyle()

	// validate source and destination
	srcFile, dstDir, err := validateSrcAndDstPaths(src, dst)
	if err != "" {
		currentPage = FatalErrorPage
	}

	// password field
	passwField := textinput.New()
	passwField.Placeholder = "Your password here"
	passwField.EchoMode = textinput.EchoPassword
	passwField.EchoCharacter = '•'
	passwField.Focus()

	// encrypt & decrypt list
	items := []list.Item{
		item(actionDecrypt),
		item(actionEncrypt),
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = fmt.Sprintf("Target file: %s", src)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = s.TitleStyle
	l.Styles.PaginationStyle = s.PaginationStyle
	l.Styles.HelpStyle = s.HelpStyle

	return &Model{
		err:            err,
		style:          s,
		list:           l,
		sourceFile:     srcFile,
		destinationDir: dstDir,
		passwField:     passwField,
		currentPage:    currentPage,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch m.currentPage {
		case OptionsPage:
			cmd = updatePageOptions(m, msg)
		case PasswordPage:
			cmd = updatePagePassword(m, msg)
		case ErrorPage:
			cmd = updatePageError(m, msg)
		case FatalErrorPage:
			cmd = updatePageFatalError(msg)
		case StatusPage:
			cmd = updatePageStatus(msg)
		}
	}

	return m, cmd
}

func (m *Model) View() string {
	switch m.currentPage {
	case OptionsPage:
		return renderPageOptions(m)
	case PasswordPage:
		return renderPagePasswords(m)
	case ErrorPage:
		return renderPageError(m)
	case FatalErrorPage:
		return renderPageError(m)
	case StatusPage:
		return renderPageStatus(m)
	}
	return "you should not be here... better quit"
}

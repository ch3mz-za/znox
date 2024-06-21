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

	pages       map[PageID]Page
	currentPage Page
	lastPage    Page
	style       *Style
}

func New(src, dst string) *Model {
	s := DefaultStyle()

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

	// password field
	passwordField := textinput.New()
	passwordField.Placeholder = "Your password here"
	passwordField.EchoMode = textinput.EchoPassword
	passwordField.EchoCharacter = '•'
	passwordField.Focus()

	// TODO: Load pages
	var appPages = map[PageID]Page{
		encryptionPage: &EncryptionPage{list: l},
		passwordPage:   &PasswordPage{passwordField: passwordField},
		statusPage:     &StatusPage{},
		errorPage:      &ErrorPage{},
		fatalErrorPage: &FatalErrorPage{},
	}

	currentPage := appPages[encryptionPage]

	// validate source and destination
	srcFile, dstDir, err := validateSrcAndDstPaths(src, dst)
	if err != "" {
		currentPage = appPages[fatalErrorPage]
	}

	return &Model{
		err:            err,
		style:          s,
		sourceFile:     srcFile,
		destinationDir: dstDir,
		currentPage:    currentPage,
		pages:          appPages,
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

		p := m.pages[encryptionPage].(*EncryptionPage)
		p.list.SetWidth(msg.Width)

		return m, nil

	case tea.KeyMsg:
		cmd = m.currentPage.Update(m, msg)
	}

	return m, cmd
}

func (m *Model) View() string {
	return m.currentPage.Render(m)
}

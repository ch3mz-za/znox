package iface

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	quit        bool
	width       int
	height      int
	index       int
	styles      *Styles
	answerField textinput.Model
	prompts     []Prompt
}

func New(prompts []Prompt) *model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Placeholder = "Your password here"
	answerField.Focus()
	return &model{
		prompts:     prompts,
		answerField: answerField,
		styles:      styles,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.prompts[m.index]

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quit = true
			return m, tea.Quit
		case "enter":
			current.response = m.answerField.Value()
			m.answerField.SetValue("")
			m.Next()
			return m, nil
		}
	}
	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.prompts[m.index].prompt,
			m.styles.InputField.Render(m.answerField.View()),
		),
	)
}

func (m *model) Next() {
	if m.index < len(m.prompts)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

type Prompt struct {
	prompt   string
	response string
}

func NewPrompt(prompt string) Prompt {
	return Prompt{prompt: prompt}
}

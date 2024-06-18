package app

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("36"))
	borderColor       = lipgloss.Color("36")
	errorColor        = lipgloss.Color("#FF2D00")
)

type Style struct {
	InputField      lipgloss.Style
	ListStyle       lipgloss.Style
	ErrorStyle      lipgloss.Style
	StatusStyle     lipgloss.Style
	TitleStyle      lipgloss.Style
	PaginationStyle lipgloss.Style
	HelpStyle       lipgloss.Style
}

func DefaultStyle() *Style {
	return &Style{
		InputField:      lipgloss.NewStyle().BorderForeground(borderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80),
		ListStyle:       lipgloss.NewStyle().BorderForeground(borderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(46),
		ErrorStyle:      lipgloss.NewStyle().BorderForeground(errorColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(46).AlignHorizontal(lipgloss.Center),
		StatusStyle:     lipgloss.NewStyle().BorderForeground(borderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(46).AlignHorizontal(lipgloss.Center),
		PaginationStyle: list.DefaultStyles().PaginationStyle.BorderForeground(borderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80),
		HelpStyle:       list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1),
	}
}

package ui

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor = lipgloss.Color("#7D56F4")
	mutedColor   = lipgloss.Color("#777777")
	accentColor  = lipgloss.Color("#04B575")

	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(primaryColor).
			Padding(0, 1)
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)
	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(primaryColor).
			Bold(true).
			Padding(0, 1)
	DetailLabelStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true)
	HelpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)
	StatusStyle = lipgloss.NewStyle().
			Foreground(accentColor)
)

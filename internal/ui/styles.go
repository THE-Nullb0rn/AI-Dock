package ui

import "github.com/charmbracelet/lipgloss"

// Palette — cohesive violet/teal/neutral theme
var (
	colorPrimary  = lipgloss.Color("#7D56F4") // violet  — brand
	colorAccent   = lipgloss.Color("#04B575") // mint    — positive/active
	colorMuted    = lipgloss.Color("#555566") // dark slate — secondary text
	colorSubtle   = lipgloss.Color("#9090A8") // lighter muted — hints/counts
	colorSelected = lipgloss.Color("#9F79FF") // lighter violet — selected highlight bg
	colorDivider  = lipgloss.Color("#2A2A3A") // near-black — subtle separators
	colorBorder   = lipgloss.Color("#6A44D4") // slightly deeper violet border
	colorText     = lipgloss.Color("#E8E6F0") // near-white — primary text
	colorBg       = lipgloss.Color("#1A1A2E") // very dark blue-black — selected bg
)

var (
	// AppStyle — outer shell with generous breathing room
	AppStyle = lipgloss.NewStyle().
		Padding(1, 2)

	// TitleBadgeStyle — "aidock" pill in the header
	TitleBadgeStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(colorPrimary).
		Padding(0, 2).
		MarginBottom(1)

	// SectionHeaderStyle — "Tools", "Add Tool" section titles
	SectionHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary).
		MarginBottom(0)

	// PanelStyle — main content box with rounded violet border + inner padding
	PanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1, 3)

	// SelectedStyle — list-item highlight: filled bg + bold white text
	SelectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(colorBg).
		Bold(true).
		Padding(0, 1)

	// ToolNameStyle — primary bright name in list
	ToolNameStyle = lipgloss.NewStyle().
		Foreground(colorText).
		Bold(true)

	// VariantCountStyle — muted "N variants" in list
	VariantCountStyle = lipgloss.NewStyle().
		Foreground(colorSubtle)

	// DetailToolNameStyle — tool heading inside detail view
	DetailToolNameStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorText).
		Background(colorPrimary).
		Padding(0, 1).
		MarginBottom(1)

	// DetailLabelStyle — variant label text (bold accent)
	DetailLabelStyle = lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true)

	// DetailSelectedMarker — "▶" prefix for the active variant
	DetailSelectedMarker = lipgloss.NewStyle().
		Foreground(colorSelected).
		Bold(true)

	// DetailMetaStyle — path / command lines (shown when details are expanded)
	DetailMetaStyle = lipgloss.NewStyle().
		Foreground(colorMuted).
		PaddingLeft(4)

	// DividerStyle — subtle horizontal rule between variants
	DividerStyle = lipgloss.NewStyle().
		Foreground(colorDivider)

	// HelpStyle — footer keybinding hints, clearly secondary
	HelpStyle = lipgloss.NewStyle().
		Foreground(colorMuted).
		Italic(true)

	// HintKeyStyle — the key name part of a hint, slightly brighter
	HintKeyStyle = lipgloss.NewStyle().
		Foreground(colorSubtle).
		Bold(true)

	// StatusStyle — status/info messages
	StatusStyle = lipgloss.NewStyle().
		Foreground(colorAccent)

	// TitleStyle kept for backward compat inside add/detail text renders
	TitleStyle = DetailToolNameStyle
)

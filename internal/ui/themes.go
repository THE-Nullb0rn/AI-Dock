package ui

import "github.com/charmbracelet/lipgloss"

// Theme holds the full color palette for a visual theme.
type Theme struct {
	Name           string
	Border         lipgloss.Color // card/section borders
	Accent         lipgloss.Color // active / highlighted text
	Selected       lipgloss.Color // selected marker / text
	Text           lipgloss.Color // primary content text
	MutedText      lipgloss.Color // secondary / dimmed text
	Background     lipgloss.Color // terminal/app bg
	CardBg         lipgloss.Color // dim background fill for unselected cards
	CardBgSelected lipgloss.Color // bright background fill for selected card
	TextSelected   lipgloss.Color // high-contrast text color on selected card
	MutedSelected  lipgloss.Color // high-contrast muted text on selected card
	SurfaceBg      lipgloss.Color // slightly lighter surface for card interiors
}

// Themes is the ordered list of built-in themes.
var Themes = []Theme{
	{
		Name:           "Violet Night",
		Border:         lipgloss.Color("#6A44D4"),
		Accent:         lipgloss.Color("#9D7CE8"),
		Selected:       lipgloss.Color("#BFA7FF"),
		Text:           lipgloss.Color("#E8E6F0"),
		MutedText:      lipgloss.Color("#9090A8"),
		Background:     lipgloss.Color("#1A1A2E"),
		CardBg:         lipgloss.Color("#1E1B2E"),
		CardBgSelected: lipgloss.Color("#6A44D4"),
		TextSelected:   lipgloss.Color("#FFFFFF"),
		MutedSelected:  lipgloss.Color("#D4C4FF"),
		SurfaceBg:      lipgloss.Color("#252240"),
	},
	{
		Name:           "Cyberpunk",
		Border:         lipgloss.Color("#FF2E88"),
		Accent:         lipgloss.Color("#00F5FF"),
		Selected:       lipgloss.Color("#FF80C0"),
		Text:           lipgloss.Color("#F0F0F0"),
		MutedText:      lipgloss.Color("#888899"),
		Background:     lipgloss.Color("#0D001A"),
		CardBg:         lipgloss.Color("#1A0E1A"),
		CardBgSelected: lipgloss.Color("#FF2E88"),
		TextSelected:   lipgloss.Color("#FFFFFF"),
		MutedSelected:  lipgloss.Color("#FFD6E8"),
		SurfaceBg:      lipgloss.Color("#1F0F25"),
	},
	{
		Name:           "Nord",
		Border:         lipgloss.Color("#5E81AC"),
		Accent:         lipgloss.Color("#88C0D0"),
		Selected:       lipgloss.Color("#81A1C1"),
		Text:           lipgloss.Color("#ECEFF4"),
		MutedText:      lipgloss.Color("#4C566A"),
		Background:     lipgloss.Color("#2E3440"),
		CardBg:         lipgloss.Color("#2E3440"),
		CardBgSelected: lipgloss.Color("#5E81AC"),
		TextSelected:   lipgloss.Color("#ECEFF4"),
		MutedSelected:  lipgloss.Color("#D8DEE9"),
		SurfaceBg:      lipgloss.Color("#3B4252"),
	},
	{
		Name:           "Dracula",
		Border:         lipgloss.Color("#BD93F9"),
		Accent:         lipgloss.Color("#FF79C6"),
		Selected:       lipgloss.Color("#CFA9FF"),
		Text:           lipgloss.Color("#F8F8F2"),
		MutedText:      lipgloss.Color("#6272A4"),
		Background:     lipgloss.Color("#282A36"),
		CardBg:         lipgloss.Color("#282A36"),
		CardBgSelected: lipgloss.Color("#BD93F9"),
		TextSelected:   lipgloss.Color("#282A36"),
		MutedSelected:  lipgloss.Color("#44475A"),
		SurfaceBg:      lipgloss.Color("#313345"),
	},
	{
		Name:           "Monochrome",
		Border:         lipgloss.Color("#CCCCCC"),
		Accent:         lipgloss.Color("#FFFFFF"),
		Selected:       lipgloss.Color("#EEEEEE"),
		Text:           lipgloss.Color("#F0F0F0"),
		MutedText:      lipgloss.Color("#666666"),
		Background:     lipgloss.Color("#111111"),
		CardBg:         lipgloss.Color("#1A1A1A"),
		CardBgSelected: lipgloss.Color("#FFFFFF"),
		TextSelected:   lipgloss.Color("#000000"),
		MutedSelected:  lipgloss.Color("#555555"),
		SurfaceBg:      lipgloss.Color("#222222"),
	},
}

// ThemeByName returns the theme with the given name, defaulting to Violet Night.
func ThemeByName(name string) Theme {
	for _, t := range Themes {
		if t.Name == name {
			return t
		}
	}
	return Themes[0]
}

// ThemeIndex returns the index of a theme by name, or 0 if not found.
func ThemeIndex(name string) int {
	for i, t := range Themes {
		if t.Name == name {
			return i
		}
	}
	return 0
}

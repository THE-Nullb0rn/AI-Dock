package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

// ToolIcons maps known tool names to Nerd Font glyphs.
// Falls back to a generic glyph for unknown tools.
var ToolIcons = map[string]string{
	"Cursor":            "\ue795", // nf-dev-terminal
	"Google Antigravity": "\uf135", // nf-fa-rocket
	"GitHub Copilot":    "\ue709", // nf-dev-github
	"ChatGPT Desktop":   "\uf4fb", // nf-md-robot (U+F4FB)
	"Ollama":            "\uf1b2", // nf-fa-cube
	"Claude":            "\U000F0760", // nf-md-brain-outline
	"VS Code":           "\ue70c", // nf-dev-visualstudio
	"Neovim":            "\ue6ae", // nf-dev-vim
	"Zed":               "\uf489", // nf-md-alpha-z-box (approx)
	"Windsurf":          "\uf72c", // nf-md-waves (approx)
	"Cline":             "\ue795", // nf-dev-terminal (generic CLI)
	"Continue":          "\uf138", // nf-fa-chevron-circle-right
	"Aider":             "\uf013", // nf-fa-cog
}

// genericIcon is the fallback for unknown tool names.
const genericIcon = "▪"

// iconAccentColors is the rotating palette for icon colorization.
var iconAccentColors = []lipgloss.Color{
	lipgloss.Color("#7D56F4"), // violet  (primary)
	lipgloss.Color("#04B575"), // mint    (accent)
	lipgloss.Color("#F4A556"), // amber
	lipgloss.Color("#56C1F4"), // sky blue
	lipgloss.Color("#F45678"), // rose
	lipgloss.Color("#A8F456"), // lime
}

// ToolIcon returns the Nerd Font glyph for the given tool name,
// styled with a deterministic accent color derived from the name.
func ToolIcon(toolName string) string {
	glyph, ok := ToolIcons[toolName]
	if !ok {
		glyph = genericIcon
	}
	color := accentColorFor(toolName)
	return lipgloss.NewStyle().Foreground(color).Render(glyph)
}

// accentColorFor hashes the tool name to pick a stable color from the palette.
func accentColorFor(name string) lipgloss.Color {
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	idx := int(h.Sum32()) % len(iconAccentColors)
	if idx < 0 {
		idx = -idx
	}
	return iconAccentColors[idx]
}

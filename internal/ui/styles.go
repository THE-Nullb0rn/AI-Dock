package ui

import "github.com/charmbracelet/lipgloss"

// Styles is a compiled set of lipgloss styles built from a Theme.
type Styles struct {
	// Outer layout
	App   lipgloss.Style
	Panel lipgloss.Style

	// Header Bar
	HeaderIcon     lipgloss.Style // large/bold icon
	HeaderTitle    lipgloss.Style // "AIDock" bold text
	HeaderVersion  lipgloss.Style // "v0.1.0"
	HeaderSubtitle lipgloss.Style // "Your AI tools, in one place."
	HeaderDate     lipgloss.Style // date string
	HeaderTime     lipgloss.Style // time string
	HeaderTagline  lipgloss.Style // "Build faster with AI."
	Divider        lipgloss.Style // repeated "─"

	// Section Label
	SectionTitle lipgloss.Style // icon + "YOUR AI TOOLS" bold accent
	SectionStats lipgloss.Style // "<N> tools | <M> interfaces" muted

	// Cards – bordered approach (no solid background fills)
	Card           lipgloss.Style // unselected card (subtle border)
	CardSelected   lipgloss.Style // selected card (accent border)
	CardTitle      lipgloss.Style // tool name on unselected card
	CardTitleSel   lipgloss.Style // tool name on selected card
	CardDesc       lipgloss.Style // description on unselected card
	CardDescSel    lipgloss.Style // description on selected card
	CardTagPill    lipgloss.Style // plain-text tag for variant label
	CardTagPillSel lipgloss.Style // plain-text tag on selected card
	CardCount      lipgloss.Style // "<N> interfaces"
	CardCountSel   lipgloss.Style // "<N> interfaces" on selected
	CardArrow      lipgloss.Style // "→"
	CardArrowSel   lipgloss.Style // "→" on selected
	StarMarker     lipgloss.Style // "★" on selected card

	// Footer
	FooterKeyPill lipgloss.Style // key label in footer
	FooterAction  lipgloss.Style // action text next to key
	FooterTagline lipgloss.Style // "AIDock | Terminal powered. AI everywhere."

	// Detail view
	DetailTitle          lipgloss.Style
	DetailLabel          lipgloss.Style
	DetailLabelSel       lipgloss.Style
	DetailSelectedMarker lipgloss.Style
	DetailMeta           lipgloss.Style

	// Theme picker
	PickerCurrent  lipgloss.Style // current theme: bold + accent
	PickerSelected lipgloss.Style // marker "▶ "
	PickerItem     lipgloss.Style // other theme names

	// General hints / status
	Help          lipgloss.Style
	HintKey       lipgloss.Style
	Status        lipgloss.Style
	SectionHeader lipgloss.Style

	// Legacy aliases
	CardName     lipgloss.Style
	CardNameSel  lipgloss.Style
	CardCountOld lipgloss.Style
}

// ActiveStyles is the package-level singleton rebuilt by ApplyTheme.
var ActiveStyles Styles

// ApplyTheme rebuilds ActiveStyles from the given theme.
func ApplyTheme(t Theme) {
	ActiveStyles = Styles{
		// ── outer layout ──────────────────────────────────────────────────
		App:   lipgloss.NewStyle().Padding(1, 2),
		Panel: lipgloss.NewStyle(), // no outer border

		// ── header bar ────────────────────────────────────────────────────
		HeaderIcon: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),

		HeaderTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Text),

		HeaderVersion: lipgloss.NewStyle().
			Foreground(t.MutedText),

		HeaderSubtitle: lipgloss.NewStyle().
			Foreground(t.MutedText),

		HeaderDate: lipgloss.NewStyle().
			Foreground(t.Text),

		HeaderTime: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),

		HeaderTagline: lipgloss.NewStyle().
			Foreground(t.MutedText),

		Divider: lipgloss.NewStyle().
			Foreground(t.MutedText),

		// ── section label ─────────────────────────────────────────────────
		SectionTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),

		SectionStats: lipgloss.NewStyle().
			Foreground(t.MutedText),

		SectionHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),

		// ── cards (bordered, NO solid background fills) ───────────────────
		// Unselected: subtle muted border, no background
		Card: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.MutedText).
			Padding(0, 1),

		// Selected: accent border, no solid background fill
		CardSelected: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Accent).
			Padding(0, 1),

		CardTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Text),

		CardTitleSel: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),

		CardDesc: lipgloss.NewStyle().
			Foreground(t.MutedText),

		CardDescSel: lipgloss.NewStyle().
			Foreground(t.Text),

		// Tags: plain-text bracket style [IDE] — no nested borders
		CardTagPill: lipgloss.NewStyle().
			Foreground(t.MutedText),

		CardTagPillSel: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true),

		CardCount: lipgloss.NewStyle().
			Foreground(t.MutedText),

		CardCountSel: lipgloss.NewStyle().
			Foreground(t.Text),

		CardArrow: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.MutedText),

		CardArrowSel: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),

		StarMarker: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),

		// ── footer (compact inline style) ────────────────────────────────
		FooterKeyPill: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true),

		FooterAction: lipgloss.NewStyle().
			Foreground(t.MutedText),

		FooterTagline: lipgloss.NewStyle().
			Foreground(t.MutedText),

		// ── detail view ──────────────────────────────────────────────────
		DetailTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Text),

		DetailLabel: lipgloss.NewStyle().
			Foreground(t.Text),

		DetailLabelSel: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),

		DetailSelectedMarker: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Selected),

		DetailMeta: lipgloss.NewStyle().
			Foreground(t.MutedText).
			PaddingLeft(4),

		// ── theme picker ─────────────────────────────────────────────────
		PickerCurrent: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),

		PickerSelected: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Selected),

		PickerItem: lipgloss.NewStyle().
			Foreground(t.MutedText),

		// ── misc ─────────────────────────────────────────────────────────
		Help: lipgloss.NewStyle().
			Foreground(t.MutedText),

		HintKey: lipgloss.NewStyle().
			Foreground(t.Text).
			Bold(true),

		Status: lipgloss.NewStyle().
			Foreground(t.Accent),

		// Aliases
		CardName: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Text),

		CardNameSel: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Accent),
	}
}

func init() {
	ApplyTheme(Themes[0])
}

package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/THE-Nullb0rn/AI-Dock/internal/config"
	"github.com/THE-Nullb0rn/AI-Dock/internal/launcher"
	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ─── constants ────────────────────────────────────────────────────────────────

const (
	cardMinW = 36 // minimum card width (chars, including border)
	cardGap  = 2  // horizontal gap between cards

	// cardW is kept as a package-level alias used by calcCols; it equals
	// cardMinW so existing column-calculation tests (if any) stay green.
	cardW = cardMinW
)

// ─── tick message ─────────────────────────────────────────────────────────────

type tickMsg time.Time

func tickEverySecond() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// ─── view / step enums ────────────────────────────────────────────────────────

type view int

const (
	dashboardView view = iota
	detailView
	helpView
)

type addStep int

const (
	addToolName addStep = iota
	addToolDescription
	addVariantLabel
	addLaunchCommand
	addTerminalMode
	addAnotherVariant
	addConfirmation
)

type deleteTarget int

const (
	deleteNone deleteTarget = iota
	deleteTool
	deleteVariant
)

// ─── Model ────────────────────────────────────────────────────────────────────

type Model struct {
	tools         []model.Tool
	selected      model.Tool
	variantIndex  int
	view          view
	width, height int
	cfg           *config.Config

	// grid navigation
	gridCursor int // linear index of selected card
	cols       int // columns (recalculated on resize)

	// live clock
	currentTime time.Time

	// add flow
	adding      bool
	addStep     addStep
	input       textinput.Model
	tempTool    model.Tool
	tempVariant model.Variant

	// delete flow
	deleteTarget  deleteTarget
	deleteToolIdx int
	launchErr     error

	// detail view toggle
	showDetails bool

	// theme picker overlay
	themePicking bool
	themeIdx     int // cursor in Themes slice
}

func NewModel(cfg *config.Config) Model {
	ApplyTheme(ThemeByName(cfg.Theme))

	ti := textinput.New()
	ti.CharLimit = 256

	m := Model{
		tools:       cfg.Tools,
		cfg:         cfg,
		input:       ti,
		cols:        1,
		themeIdx:    ThemeIndex(cfg.Theme),
		currentTime: time.Now(),
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tickEverySecond())
}

// LaunchError is reported by main after Bubble Tea restores the terminal.
func (m Model) LaunchError() error { return m.launchErr }

// ─── column calculation ───────────────────────────────────────────────────────

func calcCols(termWidth int) int {
	available := termWidth - 4 // outer padding: 2 on left, 2 on right
	each := cardMinW + cardGap
	c := (available + cardGap) / each
	if c < 1 {
		c = 1
	}
	return c
}

// ─── Update ───────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.currentTime = time.Time(msg)
		return m, tickEverySecond()

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.cols = calcCols(msg.Width)
		m.gridCursor = clamp(m.gridCursor, 0, max(0, len(m.tools)-1))
		return m, nil

	case tea.KeyMsg:
		if m.adding {
			return m.updateAdd(msg)
		}
		if m.deleteTarget != deleteNone {
			return m.updateDelete(msg)
		}
		if m.themePicking {
			return m.updateThemePicker(msg)
		}
		if m.view == helpView {
			if msg.String() == "esc" || msg.String() == "?" || msg.String() == "q" {
				m.view = dashboardView
				return m, nil
			}
		}
		return m.updateNormal(msg)
	}
	return m, nil
}

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	// ── dashboard navigation ──────────────────────────────────────────────
	case "left":
		if m.view == dashboardView {
			m.gridCursor = clamp(m.gridCursor-1, 0, len(m.tools)-1)
		}
	case "right":
		if m.view == dashboardView {
			m.gridCursor = clamp(m.gridCursor+1, 0, len(m.tools)-1)
		}
	case "up":
		if m.view == dashboardView {
			m.gridCursor = clamp(m.gridCursor-m.cols, 0, len(m.tools)-1)
		}
		if m.view == detailView && m.variantIndex > 0 {
			m.variantIndex--
		}
	case "down":
		if m.view == dashboardView {
			m.gridCursor = clamp(m.gridCursor+m.cols, 0, len(m.tools)-1)
		}
		if m.view == detailView && m.variantIndex < len(m.selected.Variants)-1 {
			m.variantIndex++
		}

	// ── actions ───────────────────────────────────────────────────────────
	case "enter":
		if m.view == dashboardView && len(m.tools) > 0 {
			m.selected = m.tools[m.gridCursor]
			m.variantIndex = 0
			m.view = detailView
			m.showDetails = false
		} else if m.view == detailView && len(m.selected.Variants) > 0 {
			variant := m.selected.Variants[m.variantIndex]
			m.launchErr = launcher.Launch(variant, resolveLaunchMode(variant))
			return m, tea.Quit
		}

	case "esc":
		if m.view == detailView {
			m.view = dashboardView
			m.showDetails = false
		}

	case "i":
		if m.view == detailView {
			m.showDetails = !m.showDetails
		}

	case "a":
		if m.view == dashboardView {
			return m, m.startAdd()
		}

	case "d":
		if m.view == dashboardView && len(m.tools) > 0 {
			m.deleteTarget = deleteTool
			m.deleteToolIdx = m.gridCursor
		} else if m.view == detailView && len(m.selected.Variants) > 0 {
			m.deleteTarget = deleteVariant
		}

	case "t":
		if m.view == dashboardView {
			m.themePicking = true
			m.themeIdx = ThemeIndex(m.cfg.Theme)
		}

	case "?":
		if m.view == dashboardView {
			m.view = helpView
		}
	}
	return m, nil
}

// ─── Theme picker ─────────────────────────────────────────────────────────────

func (m Model) updateThemePicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.themePicking = false
	case "up":
		if m.themeIdx > 0 {
			m.themeIdx--
		}
	case "down":
		if m.themeIdx < len(Themes)-1 {
			m.themeIdx++
		}
	case "enter":
		chosen := Themes[m.themeIdx]
		ApplyTheme(chosen)
		m.cfg.Theme = chosen.Name
		_ = m.cfg.Save()
		m.themePicking = false
	}
	return m, nil
}

// ─── Add flow ─────────────────────────────────────────────────────────────────

func (m Model) updateAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		m.cancelAdd()
		return m, nil
	}
	switch m.addStep {
	case addTerminalMode:
		switch msg.String() {
		case "y":
			m.tempVariant.Type, m.tempVariant.LaunchMode = "cli", "new_terminal"
			m.addStep = addAnotherVariant
		case "n":
			m.tempVariant.Type, m.tempVariant.LaunchMode = "app", ""
			m.addStep = addAnotherVariant
		}
		return m, nil
	case addAnotherVariant:
		switch msg.String() {
		case "y":
			m.tempTool.Variants = append(m.tempTool.Variants, m.tempVariant)
			m.tempVariant, m.addStep = model.Variant{}, addVariantLabel
			return m, m.resetInput("")
		case "n":
			m.tempTool.Variants = append(m.tempTool.Variants, m.tempVariant)
			m.cfg.Tools = append(m.cfg.Tools, m.tempTool)
			if err := m.cfg.Save(); err != nil {
				m.launchErr = err
				return m, tea.Quit
			}
			m.refreshTools()
			m.addStep = addConfirmation
		}
		return m, nil
	case addConfirmation:
		m.cancelAdd()
		return m, nil
	}
	if msg.String() == "enter" {
		value := strings.TrimSpace(m.input.Value())
		switch m.addStep {
		case addToolName:
			if value == "" {
				return m, nil
			}
			m.tempTool = model.Tool{Name: value}
			m.addStep = addToolDescription
			return m, m.resetInput("e.g. Next-generation AI code editor (optional)")
		case addToolDescription:
			m.tempTool.Description = value
			m.addStep = addVariantLabel
			return m, m.resetInput("e.g. IDE, CLI, Agent")
		case addVariantLabel:
			if value == "" {
				return m, nil
			}
			m.tempVariant = model.Variant{Label: value}
			m.addStep = addLaunchCommand
			return m, m.resetInput("e.g. cursor, /usr/bin/code")
		case addLaunchCommand:
			if value == "" {
				return m, nil
			}
			m.tempVariant.Path = value
			m.addStep = addTerminalMode
			return m, m.resetInput("")
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) updateDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "n", "esc":
		m.deleteTarget = deleteNone
	case "y":
		switch m.deleteTarget {
		case deleteTool:
			if idx := m.deleteToolIdx; idx >= 0 && idx < len(m.tools) {
				m.tools = append(m.tools[:idx], m.tools[idx+1:]...)
				m.cfg.Tools = m.tools
				_ = m.cfg.Save()
				m.refreshTools()
			}
		case deleteVariant:
			if toolIdx := m.selectedToolIndex(); toolIdx >= 0 && m.variantIndex < len(m.tools[toolIdx].Variants) {
				tool := m.tools[toolIdx]
				tool.Variants = append(tool.Variants[:m.variantIndex], tool.Variants[m.variantIndex+1:]...)
				if len(tool.Variants) == 0 {
					m.tools = append(m.tools[:toolIdx], m.tools[toolIdx+1:]...)
					m.view = dashboardView
				} else {
					m.tools[toolIdx], m.selected = tool, tool
					if m.variantIndex >= len(tool.Variants) {
						m.variantIndex = len(tool.Variants) - 1
					}
				}
				m.cfg.Tools = m.tools
				_ = m.cfg.Save()
				m.refreshTools()
			}
		}
		m.deleteTarget = deleteNone
	}
	return m, nil
}

// ─── helper mutators ──────────────────────────────────────────────────────────

func (m *Model) startAdd() tea.Cmd {
	m.adding, m.addStep = true, addToolName
	m.tempTool, m.tempVariant = model.Tool{}, model.Variant{}
	return m.resetInput("e.g. Cursor")
}

func (m *Model) cancelAdd() {
	m.adding, m.addStep = false, addToolName
	m.tempTool, m.tempVariant = model.Tool{}, model.Variant{}
	m.input.Blur()
	m.input.SetValue("")
	m.input.Placeholder = ""
}

func (m *Model) resetInput(placeholder string) tea.Cmd {
	m.input.SetValue("")
	m.input.Placeholder = placeholder
	m.input.CursorEnd()
	return m.input.Focus()
}

func (m *Model) refreshTools() {
	m.tools = m.cfg.Tools
	m.gridCursor = clamp(m.gridCursor, 0, max(0, len(m.tools)-1))
}

func (m Model) selectedToolIndex() int {
	for i, tool := range m.tools {
		if tool.Name == m.selected.Name {
			return i
		}
	}
	return -1
}

// ─── View ─────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	s := ActiveStyles

	var content string
	switch {
	case m.adding:
		content = m.addView()
	case m.themePicking:
		content = m.themePickerView()
	case m.view == helpView:
		content = m.helpView()
	case m.view == detailView:
		content = m.renderDetailView()
	default:
		content = m.renderDashboard()
	}

	return s.App.Render(content)
}

// ─── Header Bar ───────────────────────────────────────────────────────────────

func (m Model) renderHeader(availW int) string {
	s := ActiveStyles

	// Left: Icon + Title + Version
	icon := s.HeaderIcon.Render("◈")
	title := s.HeaderTitle.Render("AIDock")
	version := s.HeaderVersion.Render("v0.1.0")
	subtitle := s.HeaderSubtitle.Render("Your AI tools, in one place.")
	leftPart := icon + " " + title + " " + version + "  " + subtitle

	// Right: Date + Time + Tagline
	now := m.currentTime
	if now.IsZero() {
		now = time.Now()
	}
	dateStr := s.HeaderDate.Render(now.Format("Mon, 02 Jan 2006"))
	timeStr := s.HeaderTime.Render(now.Format("15:04:05"))
	tagline := s.HeaderTagline.Render("Build faster with AI.")
	rightPart := tagline + "  " + dateStr + " " + timeStr

	// Single row header
	gap := availW - lipgloss.Width(leftPart) - lipgloss.Width(rightPart)
	if gap < 2 {
		gap = 2
	}
	headerRow := leftPart + strings.Repeat(" ", gap) + rightPart

	// Horizontal divider line
	divider := s.Divider.Render(strings.Repeat("─", availW))

	return headerRow + "\n" + divider
}

// ─── Section Label ────────────────────────────────────────────────────────────

func (m Model) renderSectionLabel(availW int) string {
	s := ActiveStyles

	left := s.SectionTitle.Render("❖  YOUR AI TOOLS")

	totalInterfaces := 0
	for _, t := range m.tools {
		totalInterfaces += len(t.Variants)
	}
	toolsWord := "tools"
	if len(m.tools) == 1 {
		toolsWord = "tool"
	}
	interfacesWord := "interfaces"
	if totalInterfaces == 1 {
		interfacesWord = "interface"
	}
	right := s.SectionStats.Render(fmt.Sprintf("%d %s | %d %s", len(m.tools), toolsWord, totalInterfaces, interfacesWord))

	gap := availW - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 2 {
		gap = 2
	}
	return left + strings.Repeat(" ", gap) + right
}

// ─── Dashboard ────────────────────────────────────────────────────────────────

func (m Model) renderDashboard() string {
	s := ActiveStyles

	availW := m.width - 4
	if availW < 40 {
		availW = 40
	}

	header := m.renderHeader(availW)
	section := m.renderSectionLabel(availW)

	var mainContent string
	if len(m.tools) == 0 {
		noTools := s.CardTitle.Render("No tools added yet.")
		hint := s.Help.Render("Press 'a' to add one, or 'q' to quit.")
		mainContent = noTools + "\n" + hint
	} else {
		cols := m.cols
		if cols < 1 {
			cols = 1
		}

		// Calculate card width: fill available space evenly
		cw := m.calcCardWidth(cols, availW)

		cards := make([]string, len(m.tools))
		for i, tool := range m.tools {
			cards[i] = m.renderCard(tool, i == m.gridCursor, cw)
		}

		var rows []string
		for rowStart := 0; rowStart < len(cards); rowStart += cols {
			end := rowStart + cols
			if end > len(cards) {
				end = len(cards)
			}
			var rowCards []string
			for c := rowStart; c < end; c++ {
				if c > rowStart {
					rowCards = append(rowCards, strings.Repeat(" ", cardGap))
				}
				rowCards = append(rowCards, cards[c])
			}
			row := lipgloss.JoinHorizontal(lipgloss.Top, rowCards...)
			rows = append(rows, row)
		}
		mainContent = strings.Join(rows, "\n")
	}

	footer := m.renderFooter(availW)

	var b strings.Builder
	b.WriteString(header + "\n")
	b.WriteString(section + "\n\n")
	b.WriteString(mainContent + "\n")
	if m.deleteTarget == deleteTool {
		b.WriteString("\n" + m.deletePrompt() + "\n")
	}
	b.WriteString("\n" + footer)

	return b.String()
}

// calcCardWidth computes the width each card should use to fill the available
// terminal space evenly across the given number of columns.
func (m Model) calcCardWidth(cols, availW int) int {
	if cols <= 0 {
		cols = 1
	}
	totalGaps := (cols - 1) * cardGap
	cw := (availW - totalGaps) / cols
	if cw < cardMinW {
		cw = cardMinW
	}
	return cw
}

// ─── Card Rendering ───────────────────────────────────────────────────────────
//
// Key design: Cards use Lip Gloss borders only. No background fills. All child
// content is composed as plain styled text (no nested bordered/background
// Lip Gloss containers) to avoid ANSI escape corruption.

func (m Model) renderCard(tool model.Tool, isSelected bool, cardWidth ...int) string {
	s := ActiveStyles

	// Determine card width. The variadic cardWidth lets the existing test
	// signature `renderCard(tool, bool)` keep working without a width arg.
	cw := cardMinW
	if len(cardWidth) > 0 && cardWidth[0] > 0 {
		cw = cardWidth[0]
	}

	// The border takes 2 chars on each side (left+right), padding is 1 each side
	// Total overhead for content = border(2) + padding(2) = 4
	// But lipgloss Width() on a border style sets the *outer* width including
	// the border. Padding(0,1) adds 1 char each side inside the border.
	// So inner content width = cw - 2 (border) - 2 (padding) = cw - 4
	innerW := cw - 4
	if innerW < 10 {
		innerW = 10
	}

	// ── Row 1: Icon + Tool Name ──────────────────────────────────────────
	var iconGlyph string
	if isSelected {
		iconGlyph = s.CardTitleSel.Render(ToolGlyph(tool.Name))
	} else {
		iconGlyph = ToolIcon(tool.Name)
	}

	var nameText string
	if isSelected {
		nameText = s.CardTitleSel.Render(tool.Name)
	} else {
		nameText = s.CardTitle.Render(tool.Name)
	}

	leftTitle := iconGlyph + " " + nameText

	star := " "
	if isSelected {
		star = s.StarMarker.Render("★")
	}

	gapTitle := innerW - lipgloss.Width(leftTitle) - lipgloss.Width(star)
	if gapTitle < 1 {
		maxNameLen := innerW - lipgloss.Width(iconGlyph) - 4
		if maxNameLen < 4 {
			maxNameLen = 4
		}
		truncatedName := truncate(tool.Name, maxNameLen)
		if isSelected {
			nameText = s.CardTitleSel.Render(truncatedName)
		} else {
			nameText = s.CardTitle.Render(truncatedName)
		}
		leftTitle = iconGlyph + " " + nameText
		gapTitle = innerW - lipgloss.Width(leftTitle) - lipgloss.Width(star)
		if gapTitle < 1 {
			gapTitle = 1
		}
	}
	titleRow := leftTitle + strings.Repeat(" ", gapTitle) + star

	// ── Row 2: Description ───────────────────────────────────────────────
	desc := tool.Description
	if desc == "" {
		desc = "AI tool"
	}
	desc = truncate(desc, innerW)
	var descRow string
	if isSelected {
		descRow = s.CardDescSel.Render(desc)
	} else {
		descRow = s.CardDesc.Render(desc)
	}

	// ── Row 3: Variant tags as plain text [ TAG ] ────────────────────────
	tagStyle := s.CardTagPill
	if isSelected {
		tagStyle = s.CardTagPillSel
	}

	var tagParts []string
	curW := 0
	for _, v := range tool.Variants {
		lbl := variantLabel(v)
		rendered := tagStyle.Render("[ " + lbl + " ]")
		tw := lipgloss.Width(rendered)
		if curW > 0 && curW+1+tw > innerW {
			break
		}
		if curW > 0 {
			tagParts = append(tagParts, " ")
			curW += 1
		}
		tagParts = append(tagParts, rendered)
		curW += tw
	}
	if len(tagParts) == 0 {
		tagParts = append(tagParts, tagStyle.Render("[ default ]"))
	}
	tagsRow := strings.Join(tagParts, "")

	// ── Row 4: Interface count + arrow ───────────────────────────────────
	n := len(tool.Variants)
	interfaceWord := "interfaces"
	if n == 1 {
		interfaceWord = "interface"
	}
	countStr := fmt.Sprintf("%d %s", n, interfaceWord)
	var leftCount string
	if isSelected {
		leftCount = s.CardCountSel.Render(countStr)
	} else {
		leftCount = s.CardCount.Render(countStr)
	}

	var rightArrow string
	if isSelected {
		rightArrow = s.CardArrowSel.Render("→")
	} else {
		rightArrow = s.CardArrow.Render("→")
	}

	gapBottom := innerW - lipgloss.Width(leftCount) - lipgloss.Width(rightArrow)
	if gapBottom < 1 {
		gapBottom = 1
	}
	bottomRow := leftCount + strings.Repeat(" ", gapBottom) + rightArrow

	// ── Assemble inner lines with uniform width ──────────────────────────
	padLine := func(line string) string {
		w := lipgloss.Width(line)
		if w < innerW {
			return line + strings.Repeat(" ", innerW-w)
		}
		return line
	}

	cardContent := strings.Join([]string{
		padLine(titleRow),
		padLine(descRow),
		padLine(tagsRow),
		padLine(bottomRow),
	}, "\n")

	// ── Apply outer card style (border only, no fill) ────────────────────
	cardStyle := s.Card.Width(cw - 2) // -2 for border chars
	if isSelected {
		cardStyle = s.CardSelected.Width(cw - 2)
	}

	return cardStyle.Render(cardContent)
}

// ─── Footer ───────────────────────────────────────────────────────────────────

func (m Model) renderFooter(availW int) string {
	s := ActiveStyles

	// Divider above footer
	divider := s.Divider.Render(strings.Repeat("─", availW))

	// Compact key hints: key action · key action · ...
	type keyAction struct {
		key, action string
	}
	keys := []keyAction{
		{"↑↓←→", "navigate"},
		{"enter", "open"},
		{"a", "add"},
	}
	if len(m.tools) > 0 {
		keys = append(keys, keyAction{"d", "delete"})
	}
	keys = append(keys, keyAction{"t", "theme"})
	keys = append(keys, keyAction{"?", "help"})
	keys = append(keys, keyAction{"q", "quit"})

	var parts []string
	for _, ka := range keys {
		parts = append(parts, s.FooterKeyPill.Render(ka.key)+" "+s.FooterAction.Render(ka.action))
	}
	left := strings.Join(parts, s.Divider.Render(" · "))

	tagline := s.FooterTagline.Render("AIDock | Terminal powered. AI everywhere.")

	gap := availW - lipgloss.Width(left) - lipgloss.Width(tagline)
	var footerRow string
	if gap >= 2 {
		footerRow = left + strings.Repeat(" ", gap) + tagline
	} else {
		footerRow = left + "\n" + tagline
	}

	return divider + "\n" + footerRow
}

// ─── Help View ────────────────────────────────────────────────────────────────

func (m Model) helpView() string {
	s := ActiveStyles
	var b strings.Builder
	b.WriteString(s.SectionHeader.Render("AIDock Help") + "\n\n")

	keys := [][]string{
		{"↑ / ↓ / ← / →", "Navigate between tools in the grid dashboard"},
		{"Enter", "Open selected tool to choose launch variant"},
		{"a", "Add a new tool (name, description, launch commands)"},
		{"d", "Delete the currently selected tool (with confirmation)"},
		{"t", "Open theme picker to choose from 5 distinct visual styles"},
		{"i", "Toggle path and launch command info in detail view"},
		{"Esc", "Go back / cancel current action"},
		{"q / Ctrl+C", "Quit AIDock"},
	}

	for _, k := range keys {
		pill := s.FooterKeyPill.Render(k[0])
		desc := s.DetailLabel.Render(k[1])
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, pill, "  ", desc) + "\n\n")
	}

	b.WriteString(s.Help.Render("Press Esc or ? to return to dashboard."))
	return b.String()
}

// ─── Theme picker (plain text list, NO border) ────────────────────────────────

func (m Model) themePickerView() string {
	s := ActiveStyles
	var b strings.Builder
	b.WriteString(s.SectionHeader.Render("Select Theme") + "\n\n")

	for i, t := range Themes {
		isNavSelected := i == m.themeIdx
		isCurrent := t.Name == m.cfg.Theme

		prefix := "    "
		if isNavSelected {
			prefix = "  " + s.PickerSelected.Render("▶ ")
		}

		var nameStr string
		if isCurrent {
			nameStr = s.PickerCurrent.Render(t.Name)
		} else if isNavSelected {
			nameStr = s.CardTitle.Render(t.Name)
		} else {
			nameStr = s.PickerItem.Render(t.Name)
		}

		b.WriteString(prefix + nameStr + "\n")
	}

	b.WriteString("\n")
	b.WriteString(s.Help.Render("↑↓ choose  ·  ") +
		s.HintKey.Render("enter") + " " + s.Help.Render("apply  ·  ") +
		s.HintKey.Render("esc") + " " + s.Help.Render("cancel"))
	return b.String()
}

// ─── Add view ─────────────────────────────────────────────────────────────────

func (m Model) addView() string {
	s := ActiveStyles
	var header, line, footer string
	switch m.addStep {
	case addToolName:
		header = s.SectionHeader.Render("Add Tool")
		line = "Tool Name:"
		footer = s.Help.Render("Type the name · Enter to continue · Esc to cancel")
	case addToolDescription:
		header = s.SectionHeader.Render(fmt.Sprintf("Add Tool: %s — Description", m.tempTool.Name))
		line = "Description (optional):"
		footer = s.Help.Render("Type description (or press Enter to skip) · Esc to cancel")
	case addVariantLabel:
		header = s.SectionHeader.Render(fmt.Sprintf("Add Tool: %s — New Variant", m.tempTool.Name))
		line = "Variant Label (e.g. IDE, CLI, Agentic Mode):"
		footer = s.Help.Render("Type the label · Enter to continue · Esc to cancel")
	case addLaunchCommand:
		header = s.SectionHeader.Render(fmt.Sprintf("Add Tool: %s — %s", m.tempTool.Name, m.tempVariant.Label))
		line = "Launch command or path (e.g. cursor, /usr/bin/code):"
		footer = s.Help.Render("Type the command · Enter to continue · Esc to cancel")
	case addTerminalMode:
		header = s.SectionHeader.Render(fmt.Sprintf("Add Tool: %s — %s", m.tempTool.Name, m.tempVariant.Label))
		line = "Does this need a new terminal window? " + s.HintKey.Render("(y/n)")
		footer = s.Help.Render("Press y or n · Esc to cancel")
	case addAnotherVariant:
		header = s.SectionHeader.Render(fmt.Sprintf("Add Tool: %s", m.tempTool.Name))
		line = fmt.Sprintf("Variant '%s' added. Add another variant? %s", variantLabel(m.tempVariant), s.HintKey.Render("(y/n)"))
		footer = s.Help.Render("Press y or n")
	case addConfirmation:
		msg := s.Status.Render(fmt.Sprintf("✓  Tool '%s' saved with %d variant(s).", m.tempTool.Name, len(m.tempTool.Variants)))
		return msg + "\n\n" + s.Help.Render("Press any key to return to the dashboard.")
	}
	if m.addStep == addTerminalMode || m.addStep == addAnotherVariant {
		return strings.Join([]string{header, "", line, "", footer}, "\n")
	}
	return strings.Join([]string{header, "", line, m.input.View(), "", footer}, "\n")
}

// ─── Detail view ──────────────────────────────────────────────────────────────

func (m Model) renderDetailView() string {
	s := ActiveStyles
	var b strings.Builder

	icon := ToolIcon(m.selected.Name)
	heading := s.DetailTitle.Render(icon + "  " + m.selected.Name)
	b.WriteString(heading + "\n")
	if m.selected.Description != "" {
		b.WriteString(s.CardDesc.Render(m.selected.Description) + "\n")
	}
	b.WriteString("\n")

	for i, variant := range m.selected.Variants {
		isSel := i == m.variantIndex
		marker := "  "
		if isSel {
			marker = s.DetailSelectedMarker.Render("▶ ")
		}
		vIcon := ToolIcon(m.selected.Name)
		label := variantLabel(variant)
		var labelStr string
		if isSel {
			labelStr = s.DetailLabelSel.Render(label)
		} else {
			labelStr = s.DetailLabel.Render(label)
		}
		b.WriteString(marker + vIcon + "  " + labelStr + "\n")

		if isSel && m.showDetails {
			b.WriteString(s.DetailMeta.Render("path:    "+variant.Path) + "\n")
			if variant.LaunchCmd != "" {
				b.WriteString(s.DetailMeta.Render("command: "+variant.LaunchCmd) + "\n")
			}
		}
		if i < len(m.selected.Variants)-1 {
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	navHint := s.Help.Render("↑↓ choose variant  ·  ") +
		s.HintKey.Render("enter") + " " + s.Help.Render("launch  ·  ") +
		s.HintKey.Render("d") + " " + s.Help.Render("delete  ·  ") +
		s.HintKey.Render("esc") + " " + s.Help.Render("back")
	var detailHint string
	if m.showDetails {
		detailHint = "  ·  " + s.HintKey.Render("i") + " " + s.Help.Render("hide details")
	} else {
		detailHint = "  ·  " + s.HintKey.Render("i") + " " + s.Help.Render("toggle details")
	}
	b.WriteString(navHint + detailHint + "\n")

	if m.deleteTarget == deleteVariant {
		b.WriteString("\n" + fmt.Sprintf("Delete variant '%s'? (y/n):", variantLabel(m.selected.Variants[m.variantIndex])))
	}
	return b.String()
}

func (m Model) deletePrompt() string {
	if m.deleteToolIdx >= 0 && m.deleteToolIdx < len(m.tools) {
		return fmt.Sprintf("Delete tool '%s' and all its variants? (y/n):", m.tools[m.deleteToolIdx].Name)
	}
	return ""
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func resolveLaunchMode(v model.Variant) string {
	if v.LaunchMode != "" {
		return v.LaunchMode
	}
	if strings.ToLower(strings.TrimSpace(v.Type)) == "cli" {
		return "new_terminal"
	}
	return ""
}

func variantLabel(v model.Variant) string {
	if v.Label != "" {
		return v.Label
	}
	return v.Type
}

func truncate(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	if maxRunes <= 1 {
		return "…"
	}
	return string(runes[:maxRunes-1]) + "…"
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

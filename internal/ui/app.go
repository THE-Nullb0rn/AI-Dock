package ui

import (
	"fmt"
	"strings"

	"github.com/THE-Nullb0rn/AI-Dock/internal/config"
	"github.com/THE-Nullb0rn/AI-Dock/internal/launcher"
	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type view int

const (
	listView view = iota
	detailView
)

type addStep int

const (
	addToolName addStep = iota
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

type listItem struct{ tool model.Tool }

func (i listItem) FilterValue() string { return i.tool.Name }
func (i listItem) Title() string       { return i.tool.Name }
func (i listItem) Description() string { return fmt.Sprintf("%d variant(s)", len(i.tool.Variants)) }

type Model struct {
	list          list.Model
	tools         []model.Tool
	selected      model.Tool
	variantIndex  int
	view          view
	width, height int
	cfg           *config.Config

	adding      bool
	addStep     addStep
	input       textinput.Model
	tempTool    model.Tool
	tempVariant model.Variant

	deleteTarget  deleteTarget
	deleteToolIdx int
	launchErr     error
}

func NewModel(cfg *config.Config) Model {
	items := make([]list.Item, len(cfg.Tools))
	for i, tool := range cfg.Tools {
		items[i] = listItem{tool: tool}
	}
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = SelectedStyle
	delegate.Styles.SelectedDesc = SelectedStyle
	l := list.New(items, delegate, 0, 0)
	l.Title = "Tools"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	ti := textinput.New()
	ti.CharLimit = 256
	return Model{list: l, tools: cfg.Tools, cfg: cfg, input: ti}
}

func (m Model) Init() tea.Cmd { return textinput.Blink }

// LaunchError is reported by main after Bubble Tea has restored the terminal.
func (m Model) LaunchError() error { return m.launchErr }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(msg.Width-6, msg.Height-8)
		return m, nil
	case tea.KeyMsg:
		if m.adding {
			return m.updateAdd(msg)
		}
		if m.deleteTarget != deleteNone {
			return m.updateDelete(msg)
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "a":
			if m.view == listView {
				return m, m.startAdd()
			}
		case "d":
			if m.view == listView && len(m.tools) > 0 {
				m.deleteTarget, m.deleteToolIdx = deleteTool, m.list.Index()
				return m, nil
			}
			if m.view == detailView && len(m.selected.Variants) > 0 {
				m.deleteTarget = deleteVariant
				return m, nil
			}
		case "esc":
			if m.view == detailView {
				m.view = listView
				return m, nil
			}
		case "enter":
			if m.view == listView && len(m.tools) > 0 {
				m.selected, m.variantIndex, m.view = m.tools[m.list.Index()], 0, detailView
				return m, nil
			}
			if m.view == detailView && len(m.selected.Variants) > 0 {
				variant := m.selected.Variants[m.variantIndex]
				m.launchErr = launcher.Launch(variant, resolveLaunchMode(variant))
				return m, tea.Quit
			}
		case "up":
			if m.view == detailView && m.variantIndex > 0 {
				m.variantIndex--
				return m, nil
			}
		case "down":
			if m.view == detailView && m.variantIndex < len(m.selected.Variants)-1 {
				m.variantIndex++
				return m, nil
			}
		}
	}
	if m.view == listView {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

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
		if value == "" {
			return m, nil
		}
		switch m.addStep {
		case addToolName:
			m.tempTool, m.addStep = model.Tool{Name: value}, addVariantLabel
		case addVariantLabel:
			m.tempVariant, m.addStep = model.Variant{Label: value}, addLaunchCommand
		case addLaunchCommand:
			m.tempVariant.Path, m.addStep = value, addTerminalMode
		}
		return m, m.resetInput("")
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
					m.view = listView
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

func (m *Model) startAdd() tea.Cmd {
	m.adding, m.addStep = true, addToolName
	m.tempTool, m.tempVariant = model.Tool{}, model.Variant{}
	return m.resetInput("e.g. Antigravity")
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
	items := make([]list.Item, len(m.tools))
	for i, tool := range m.tools {
		items[i] = listItem{tool: tool}
	}
	m.list.SetItems(items)
}

func (m Model) selectedToolIndex() int {
	for i, tool := range m.tools {
		if tool.Name == m.selected.Name {
			return i
		}
	}
	return -1
}

func (m Model) View() string {
	var content string
	switch {
	case m.adding:
		content = m.addView()
	case m.view == detailView:
		content = m.detailView()
	default:
		content = m.listView()
	}
	return AppStyle.Render(lipgloss.JoinVertical(lipgloss.Left, TitleStyle.Render(" aidock "), PanelStyle.Width(max(20, m.width-6)).Render(content)))
}

func (m Model) listView() string {
	if len(m.tools) == 0 {
		return "No tools added yet.\nPress 'a' to add one, or 'q' to quit."
	}
	content := m.list.View()
	if m.deleteTarget == deleteTool {
		content += "\n\n" + m.deletePrompt()
	}
	return content
}

func (m Model) addView() string {
	var header, line, footer string
	switch m.addStep {
	case addToolName:
		header, line, footer = "Add Tool", "Tool Name:", "Type the name, press Enter to continue. Esc to cancel."
	case addVariantLabel:
		header, line, footer = fmt.Sprintf("Add Tool: %s — New Variant", m.tempTool.Name), "Variant Label (e.g. IDE, CLI, Agentic Mode):", "Type the label, press Enter to continue. Esc to cancel."
	case addLaunchCommand:
		header, line, footer = fmt.Sprintf("Add Tool: %s — New Variant: %s", m.tempTool.Name, m.tempVariant.Label), "Launch Command or Path (e.g. cursor, /usr/bin/code, code .):", "Type the command, press Enter to continue. Esc to cancel."
	case addTerminalMode:
		header, line, footer = fmt.Sprintf("Add Tool: %s — New Variant: %s", m.tempTool.Name, m.tempVariant.Label), "Does this need a terminal window to run? (y/n):", "Press y or n. Esc to cancel."
	case addAnotherVariant:
		header, line, footer = fmt.Sprintf("Add Tool: %s", m.tempTool.Name), fmt.Sprintf("Variant '%s' added. Add another variant? (y/n):", variantLabel(m.tempVariant)), "Press y or n."
	case addConfirmation:
		return fmt.Sprintf("Tool '%s' saved with %d variant(s).\n\nPress any key to return to the list.", m.tempTool.Name, len(m.tempTool.Variants))
	}
	if m.addStep == addTerminalMode || m.addStep == addAnotherVariant {
		return strings.Join([]string{header, "", line, "", footer}, "\n")
	}
	return strings.Join([]string{header, "", line, m.input.View(), "", footer}, "\n")
}

func (m Model) detailView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(m.selected.Name))
	b.WriteString("\n\n")
	for i, variant := range m.selected.Variants {
		marker := " "
		if i == m.variantIndex {
			marker = ">"
		}
		fmt.Fprintf(&b, "%s %s %s\n", marker, DetailLabelStyle.Render(fmt.Sprintf("Variant %d", i+1)), variantLabel(variant))
		fmt.Fprintf(&b, "  path: %s\n  command: %s\n", variant.Path, variant.LaunchCmd)
	}
	b.WriteString("\nUse up/down to choose a variant. Press enter to launch it.\n")
	if m.deleteTarget == deleteVariant {
		fmt.Fprintf(&b, "\nDelete variant '%s'? (y/n):", variantLabel(m.selected.Variants[m.variantIndex]))
	}
	return b.String()
}

func (m Model) deletePrompt() string {
	if m.deleteToolIdx >= 0 && m.deleteToolIdx < len(m.tools) {
		return fmt.Sprintf("Delete tool '%s' and all its variants? (y/n):", m.tools[m.deleteToolIdx].Name)
	}
	return ""
}

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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

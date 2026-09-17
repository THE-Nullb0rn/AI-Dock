package ui

import (
	"fmt"
	"strings"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
	"github.com/THE-Nullb0rn/AI-Dock/internal/launcher"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type view int

const (
	listView view = iota
	detailView
)

type listItem struct {
	tool model.Tool
}

func (i listItem) FilterValue() string { return i.tool.Name }
func (i listItem) Title() string       { return i.tool.Name }
func (i listItem) Description() string { return fmt.Sprintf("%d variant(s)", len(i.tool.Variants)) }

type Model struct {
	list     list.Model
	tools    []model.Tool
	selected model.Tool
	view     view
	width    int
	height   int
}

func NewModel(tools []model.Tool) Model {
	items := make([]list.Item, len(tools))
	for i, tool := range tools {
		items[i] = listItem{tool: tool}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = SelectedStyle
	delegate.Styles.SelectedDesc = SelectedStyle
	l := list.New(items, delegate, 0, 0)
	l.Title = "Tools"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	return Model{list: l, tools: tools}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(msg.Width-6, msg.Height-8)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.view == listView && len(m.tools) > 0 {
				m.selected = m.tools[m.list.Index()]
				m.view = detailView
				return m, nil
			}
			if m.view == detailView && len(m.selected.Variants) > 0 {
				variant := m.selected.Variants[0]
				if err := launcher.Launch(variant); err != nil {
					fmt.Println("launch error:", err)
				}
				return m, nil
			}
		case "esc":
			if m.view == detailView {
				m.view = listView
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

func (m Model) View() string {
	var content string
	if m.view == listView {
		content = m.list.View()
	} else {
		content = m.detailView()
	}

	title := TitleStyle.Render(" aidock ")
	panel := PanelStyle.Width(max(20, m.width-6)).Render(content)
	help := HelpStyle.Render("↑/↓ navigate  enter select/launch  esc back  q quit")
	return AppStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, panel, help))
}

func (m Model) detailView() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(m.selected.Name))
	b.WriteString("\n\n")
	if len(m.selected.Variants) == 0 {
		b.WriteString("No variants detected.\n")
	} else {
		for i, variant := range m.selected.Variants {
			fmt.Fprintf(&b, "%s %s\n", DetailLabelStyle.Render(fmt.Sprintf("Variant %d", i+1)), variant.Type)
			fmt.Fprintf(&b, "  path: %s\n  command: %s\n", variant.Path, variant.LaunchCmd)
		}
		b.WriteString("\nPress enter to launch the first variant (stub).\n")
	}
	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

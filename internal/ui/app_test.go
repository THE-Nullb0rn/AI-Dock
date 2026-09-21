package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/THE-Nullb0rn/AI-Dock/internal/config"
	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func TestDashboardLayout(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	for _, theme := range Themes {
		cfg.Theme = theme.Name
		m := NewModel(cfg)
		tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
		view := tm.View()

		// 1. Header Bar
		if !strings.Contains(view, "AIDock") {
			t.Errorf("Theme %s: dashboard missing 'AIDock' title", theme.Name)
		}
		if !strings.Contains(view, "v0.1.0") {
			t.Errorf("Theme %s: dashboard missing version 'v0.1.0'", theme.Name)
		}
		if !strings.Contains(view, "Your AI tools, in one place.") {
			t.Errorf("Theme %s: dashboard missing subtitle", theme.Name)
		}
		if !strings.Contains(view, "Build faster with AI.") {
			t.Errorf("Theme %s: dashboard missing tagline", theme.Name)
		}
		if !strings.Contains(view, "─") {
			t.Errorf("Theme %s: dashboard missing divider line", theme.Name)
		}

		// 2. Section Label
		if !strings.Contains(view, "YOUR AI TOOLS") {
			t.Errorf("Theme %s: dashboard missing 'YOUR AI TOOLS' section label", theme.Name)
		}
		if !strings.Contains(view, "tools") || !strings.Contains(view, "interfaces") {
			t.Errorf("Theme %s: dashboard missing tools/interfaces summary stats", theme.Name)
		}

		// 3. Card Elements
		if !strings.Contains(view, "→") {
			t.Errorf("Theme %s: dashboard missing arrow '→' on cards", theme.Name)
		}
		if !strings.Contains(view, "★") {
			t.Errorf("Theme %s: dashboard missing star '★' on selected card", theme.Name)
		}

		// 4. Footer Key Pills & Tagline
		if !strings.Contains(view, "navigate") || !strings.Contains(view, "open") || !strings.Contains(view, "theme") {
			t.Errorf("Theme %s: dashboard missing key actions in footer", theme.Name)
		}
		if !strings.Contains(view, "Terminal powered. AI everywhere.") {
			t.Errorf("Theme %s: dashboard missing footer tagline", theme.Name)
		}

		// 5. Detail View
		detailM, _ := tm.Update(tea.KeyMsg{Type: tea.KeyEnter})
		detailView := detailM.View()
		if !strings.Contains(detailView, "choose variant") {
			t.Errorf("Theme %s: detail view missing navigation hints", theme.Name)
		}

		// 6. 'i' toggle in Detail View
		detailWithInfo, _ := detailM.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
		infoView := detailWithInfo.View()
		if !strings.Contains(infoView, "path:") {
			t.Errorf("Theme %s: detail view missing 'path:' when info toggled", theme.Name)
		}

		// 7. Theme Picker View
		pickerM, _ := tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
		pickerView := pickerM.View()
		if !strings.Contains(pickerView, "Select Theme") {
			t.Errorf("Theme %s: theme picker missing 'Select Theme' title", theme.Name)
		}
	}
}

func TestSelectedCardHasVisibleForecolors(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	cfg := config.NewConfigForTest([]model.Tool{{Name: "Cursor", Description: "Test tool", Variants: []model.Variant{{Label: "IDE", Type: "app"}}}}, "Violet Night")

	m := NewModel(cfg)
	// force selected grid cursor 0 and window size
	_, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Render single selected card in isolation by calling renderCard
	// Use ActiveStyles and ensure theme applied
	card := m.renderCard(cfg.Tools[0], true)

	// verify each non-empty line contains a truecolor foreground escape (38;2;...)
	lines := strings.Split(card, "\n")
	found := false
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		// lipgloss uses ANSI with 24-bit sequences that include '38;2' for foreground
		if strings.Contains(ln, "38;2") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected selected card rendering to include truecolor foreground codes, got:\n%s", card)
	}
}

func TestCardLineWidthsUniform(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	cfg := config.NewConfigForTest([]model.Tool{{Name: "Cursor", Description: "Test tool", Variants: []model.Variant{{Label: "IDE", Type: "app"}}}}, "Violet Night")

	m := NewModel(cfg)
	_, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Check both selected and unselected rendering
	for _, sel := range []bool{true, false} {
		card := m.renderCard(cfg.Tools[0], sel)
		lines := strings.Split(card, "\n")
		var widths []int
		for _, ln := range lines {
			if strings.TrimSpace(ln) == "" {
				continue
			}
			widths = append(widths, lipgloss.Width(ln))
		}
		if len(widths) == 0 {
			t.Fatalf("no non-empty lines in card render")
		}
		for i := 1; i < len(widths); i++ {
			if widths[i] != widths[0] {
				t.Fatalf("card lines have non-uniform widths (first=%d, idx=%d val=%d)\ncard:\n%s", widths[0], i, widths[i], card)
			}
		}
	}
}

func TestLiveClockTick(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	cfg := config.NewConfigForTest([]model.Tool{
		{Name: "Cursor", Variants: []model.Variant{{Label: "IDE", Type: "app"}}},
	}, "Violet Night")

	m := NewModel(cfg)
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	testTime := time.Date(2026, time.September, 21, 14, 30, 45, 0, time.UTC)
	tm, cmd := tm.Update(tickMsg(testTime))
	if cmd == nil {
		t.Errorf("expected non-nil tea.Cmd for next tick")
	}

	view := tm.View()
	if !strings.Contains(view, "21 Sep 2026") {
		t.Errorf("expected view to contain formatted date '21 Sep 2026', got:\n%s", view)
	}
	if !strings.Contains(view, "14:30:45") {
		t.Errorf("expected view to contain formatted time '14:30:45', got:\n%s", view)
	}
}

func TestAddFlowWithDescription(t *testing.T) {
	t.Setenv("AIDOCK_CONFIG", t.TempDir()+"/tools.json")
	lipgloss.SetColorProfile(termenv.TrueColor)
	cfg := config.NewConfigForTest([]model.Tool{}, "Violet Night")

	m := NewModel(cfg)
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Press 'a' to start adding
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	modelObj := tm.(Model)
	if modelObj.addStep != addToolName {
		t.Fatalf("expected step addToolName, got %d", modelObj.addStep)
	}

	// 1. Enter Tool Name
	for _, r := range "SuperAI" {
		tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	modelObj = tm.(Model)
	if modelObj.addStep != addToolDescription {
		t.Fatalf("expected step addToolDescription, got %d", modelObj.addStep)
	}

	// 2. Enter Tool Description
	for _, r := range "Autonomous coding agent" {
		tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	modelObj = tm.(Model)
	if modelObj.addStep != addVariantLabel {
		t.Fatalf("expected step addVariantLabel, got %d", modelObj.addStep)
	}
	if modelObj.tempTool.Description != "Autonomous coding agent" {
		t.Errorf("expected description 'Autonomous coding agent', got %q", modelObj.tempTool.Description)
	}

	// 3. Enter Variant Label
	for _, r := range "CLI" {
		tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	modelObj = tm.(Model)
	if modelObj.addStep != addLaunchCommand {
		t.Fatalf("expected step addLaunchCommand, got %d", modelObj.addStep)
	}

	// 4. Enter Launch Command
	for _, r := range "superai" {
		tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	modelObj = tm.(Model)
	if modelObj.addStep != addTerminalMode {
		t.Fatalf("expected step addTerminalMode, got %d", modelObj.addStep)
	}

	// 5. Terminal mode (y)
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	modelObj = tm.(Model)
	if modelObj.addStep != addAnotherVariant {
		t.Fatalf("expected step addAnotherVariant, got %d", modelObj.addStep)
	}

	// 6. Another variant? (n)
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	modelObj = tm.(Model)
	if modelObj.addStep != addConfirmation {
		t.Fatalf("expected step addConfirmation, got %d, launchErr=%v", modelObj.addStep, modelObj.launchErr)
	}

	// Verify tool was saved to config with description
	if len(cfg.Tools) != 1 {
		t.Fatalf("expected 1 tool in cfg, got %d", len(cfg.Tools))
	}
	if cfg.Tools[0].Description != "Autonomous coding agent" {
		t.Errorf("expected saved tool description 'Autonomous coding agent', got %q", cfg.Tools[0].Description)
	}
}

func TestAddFlowSkipDescription(t *testing.T) {
	t.Setenv("AIDOCK_CONFIG", t.TempDir()+"/tools.json")
	lipgloss.SetColorProfile(termenv.TrueColor)
	cfg := &config.Config{
		Tools: []model.Tool{},
		Theme: "Nord",
	}

	m := NewModel(cfg)
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Press 'a'
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	// Name
	for _, r := range "QuickTool" {
		tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Skip description by pressing Enter immediately
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	modelObj := tm.(Model)
	if modelObj.addStep != addVariantLabel {
		t.Fatalf("expected step addVariantLabel after skipping description, got %d", modelObj.addStep)
	}
	if modelObj.tempTool.Description != "" {
		t.Errorf("expected empty description, got %q", modelObj.tempTool.Description)
	}
}

func TestThemeSwitching(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	cfg := &config.Config{
		Tools: []model.Tool{
			{Name: "Cursor", Variants: []model.Variant{{Label: "IDE", Type: "app"}}},
		},
		Theme: "Violet Night",
	}

	m := NewModel(cfg)
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Press 't' to open theme picker
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})

	// Down arrow to second theme (Cyberpunk)
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Enter to apply
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cfg.Theme != "Cyberpunk" {
		t.Fatalf("expected theme to be Cyberpunk, got %s", cfg.Theme)
	}
}

func TestEdgeCases(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)

	// 1. Empty state (0 tools)
	cfg0 := config.NewConfigForTest([]model.Tool{}, "Nord")
	m0 := NewModel(cfg0)
	tm0, _ := m0.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	view0 := tm0.View()
	if !strings.Contains(view0, "No tools added yet.") {
		t.Errorf("expected empty state message, got:\n%s", view0)
	}
	if !strings.Contains(view0, "0 tools | 0 interfaces") {
		t.Errorf("expected '0 tools | 0 interfaces' in empty state, got:\n%s", view0)
	}

	// 2. Exactly 1 tool (singular nouns)
	cfg1 := config.NewConfigForTest([]model.Tool{
		{
			Name:        "SingleTool",
			Description: "Solitary helper",
			Variants:    []model.Variant{{Label: "App", Type: "app"}},
		},
	}, "Dracula")
	m1 := NewModel(cfg1)
	tm1, _ := m1.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	view1 := tm1.View()
	if !strings.Contains(view1, "1 tool | 1 interface") {
		t.Errorf("expected singular '1 tool | 1 interface', got:\n%s", view1)
	}
	if !strings.Contains(view1, "1 interface") {
		t.Errorf("expected card to show '1 interface', got:\n%s", view1)
	}

	// 3. Narrow terminal (e.g. 50 columns)
	tmNarrow, _ := m1.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
	viewNarrow := tmNarrow.View()
	if len(viewNarrow) == 0 {
		t.Errorf("expected non-empty view on narrow terminal")
	}

	// 4. Help view with '?'
	mHelp := NewModel(cfg1)
	tmHelp, _ := mHelp.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	tmHelp, _ = tmHelp.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	viewHelp := tmHelp.View()
	if !strings.Contains(viewHelp, "AIDock Help") {
		t.Errorf("expected help view, got:\n%s", viewHelp)
	}

	// Esc returns to dashboard
	tmBack, _ := tmHelp.Update(tea.KeyMsg{Type: tea.KeyEsc})
	viewBack := tmBack.View()
	if !strings.Contains(viewBack, "YOUR AI TOOLS") {
		t.Errorf("expected return to dashboard from help view")
	}
}

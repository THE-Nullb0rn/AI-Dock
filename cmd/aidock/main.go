package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/THE-Nullb0rn/AI-Dock/internal/config"
	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
	"github.com/THE-Nullb0rn/AI-Dock/internal/ui"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	tools := cfg.Tools
	if len(tools) == 0 {
		tools = []model.Tool{
			{
				Name: "Claude",
				Variants: []model.Variant{
					{Type: "cli", Path: "claude", LaunchCmd: "claude", LaunchMode: "new_terminal"},
					{Type: "app", Path: "/Applications/Claude.app", LaunchCmd: "open Claude"},
				},
			},
			{
				Name: "Cursor",
				Variants: []model.Variant{
					{Type: "ide", Path: "cursor", LaunchCmd: "cursor .", LaunchMode: "same_terminal"},
				},
			},
		}
	}

	p := tea.NewProgram(ui.NewModel(tools), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("error running aidock:", err)
	}
}

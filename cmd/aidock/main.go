package main

import (
	"fmt"
	"log"
	"os"

	"github.com/THE-Nullb0rn/AI-Dock/internal/config"
	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
	"github.com/THE-Nullb0rn/AI-Dock/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(ui.NewModel(cfg), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("error running aidock:", err)
		return
	}
	if model, ok := finalModel.(ui.Model); ok && model.LaunchError() != nil {
		fmt.Fprintln(os.Stderr, model.LaunchError())
	}
}

func filterTools(tools []model.Tool) []model.Tool {
	filtered := make([]model.Tool, 0, len(tools))
	for _, tool := range tools {
		if len(tool.Variants) == 0 {
			continue
		}
		filtered = append(filtered, tool)
	}
	return filtered
}

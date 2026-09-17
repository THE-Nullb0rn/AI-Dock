package detector

import "github.com/THE-Nullb0rn/AI-Dock/internal/model"

type Detector interface {
	Name() string
	Detect() []model.Tool
}

var registered = []Detector{
	ClaudeDetector{},
	CursorDetector{},
	VSCodeDetector{},
}

func DetectAll() []model.Tool {
	merged := make([]model.Tool, 0, len(registered))
	indexByName := make(map[string]int)

	for _, detector := range registered {
		for _, tool := range detector.Detect() {
			if len(tool.Variants) == 0 {
				continue
			}

			if index, ok := indexByName[tool.Name]; ok {
				merged[index].Variants = append(merged[index].Variants, tool.Variants...)
				continue
			}

			indexByName[tool.Name] = len(merged)
			merged = append(merged, tool)
		}
	}

	return merged
}

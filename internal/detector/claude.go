package detector

import "github.com/THE-Nullb0rn/AI-Dock/internal/model"

type ClaudeDetector struct{}

func (ClaudeDetector) Name() string {
	return "Claude"
}

func (ClaudeDetector) Detect() []model.Tool {
	variants := make([]model.Variant, 0, 2)

	if path, err := execLookPath("claude"); err == nil {
		variants = append(variants, model.Variant{
			Type:      "cli",
			Path:      path,
			LaunchCmd: "claude",
		})
	}

	variants = append(variants, claudeDesktopVariants()...)

	return []model.Tool{{Name: "Claude", Variants: variants}}
}

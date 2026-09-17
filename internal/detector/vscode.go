package detector

import "github.com/THE-Nullb0rn/AI-Dock/internal/model"

type VSCodeDetector struct{}

func (VSCodeDetector) Name() string {
	return "VS Code"
}

func (VSCodeDetector) Detect() []model.Tool {
	variants := make([]model.Variant, 0, 2)

	if path, err := execLookPath("code"); err == nil {
		variants = append(variants, model.Variant{
			Type:      "cli",
			Path:      path,
			LaunchCmd: "code",
		})
	}

	variants = append(variants, vscodeDesktopVariants()...)

	return []model.Tool{{Name: "VS Code", Variants: variants}}
}

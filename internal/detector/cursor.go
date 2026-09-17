package detector

import "github.com/THE-Nullb0rn/AI-Dock/internal/model"

type CursorDetector struct{}

func (CursorDetector) Name() string {
	return "Cursor"
}

func (CursorDetector) Detect() []model.Tool {
	variants := make([]model.Variant, 0, 2)

	if path, err := execLookPath("cursor"); err == nil {
		variants = append(variants, model.Variant{
			Type:      "cli",
			Path:      path,
			LaunchCmd: "cursor",
		})
	}

	variants = append(variants, cursorDesktopVariants()...)

	return []model.Tool{{Name: "Cursor", Variants: variants}}
}

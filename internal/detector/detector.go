package detector

import "github.com/THE-Nullb0rn/AI-Dock/internal/model"

type Detector interface {
	Detect() []model.Tool
}

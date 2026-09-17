//go:build windows

package detector

import (
	"os"
	"path/filepath"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

func claudeDesktopVariants() []model.Variant {
	// Check the standard Windows application folders where Claude is typically installed.
	patterns := []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Claude"),
		filepath.Join(os.Getenv("APPDATA"), "Claude"),
	}

	return collectVariants(patterns, "", false, "app", "")
}

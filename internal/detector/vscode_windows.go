//go:build windows

package detector

import (
	"os"
	"path/filepath"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

func vscodeDesktopVariants() []model.Variant {
	// Check the standard Windows application folders where Visual Studio Code is typically installed.
	patterns := []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Microsoft VS Code"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Code"),
	}

	return collectVariants(patterns, "", false, "app", "")
}

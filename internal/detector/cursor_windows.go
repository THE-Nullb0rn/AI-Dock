//go:build windows

package detector

import (
	"os"
	"path/filepath"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

func cursorDesktopVariants() []model.Variant {
	// Check the standard Windows application folder where Cursor is typically installed.
	patterns := []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "cursor"),
	}

	return collectVariants(patterns, "", false, "ide", "")
}

//go:build linux

package detector

import (
	"os"
	"path/filepath"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

func vscodeDesktopVariants() []model.Variant {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	// Check desktop launchers for Visual Studio Code before probing common install locations.
	desktopPatterns := []string{
		filepath.Join(home, ".local/share/applications", "*.desktop"),
		"/usr/share/applications/*code*.desktop",
	}
	variants := collectVariants(desktopPatterns, "code", true, "app", "")

	// Check common package install directories for Visual Studio Code binaries or bundles.
	variants = append(variants, collectVariants([]string{"/opt/*code*"}, "", false, "app", "")...)
	return variants
}

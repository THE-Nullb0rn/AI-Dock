//go:build linux

package detector

import (
	"os"
	"path/filepath"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

func cursorDesktopVariants() []model.Variant {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	// Check desktop launchers for Cursor before probing common install locations.
	desktopPatterns := []string{
		filepath.Join(home, ".local/share/applications", "*.desktop"),
		"/usr/share/applications/*cursor*.desktop",
	}
	variants := collectVariants(desktopPatterns, "cursor", true, "ide", "")

	// Check common package install directories for Cursor binaries or bundles.
	variants = append(variants, collectVariants([]string{"/opt/*cursor*"}, "", false, "ide", "")...)
	return variants
}

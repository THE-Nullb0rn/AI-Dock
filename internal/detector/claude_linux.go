//go:build linux

package detector

import (
	"os"
	"path/filepath"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

func claudeDesktopVariants() []model.Variant {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	// Check the user's desktop launchers for Claude before probing common install paths.
	desktopPatterns := []string{
		filepath.Join(home, ".local/share/applications", "*.desktop"),
		"/usr/share/applications/*claude*.desktop",
	}
	variants := collectVariants(desktopPatterns, "claude", true, "app", "")

	// Check common package install directories for Claude binaries or bundles.
	variants = append(variants, collectVariants([]string{"/opt/*claude*"}, "", false, "app", "")...)
	return variants
}

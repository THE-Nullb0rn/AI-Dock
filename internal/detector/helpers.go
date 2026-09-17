package detector

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

func collectVariants(patterns []string, needle string, requireContent bool, variantType string, launchCmd string) []model.Variant {
	variants := make([]model.Variant, 0)
	seen := make(map[string]struct{})
	lowerNeedle := strings.ToLower(needle)

	for _, pattern := range patterns {
		for _, match := range matchingPaths(pattern) {
			if _, ok := seen[match]; ok {
				continue
			}

			info, err := os.Stat(match)
			if err != nil {
				continue
			}

			if requireContent {
				if info.IsDir() {
					continue
				}

				data, err := os.ReadFile(match)
				if err != nil {
					continue
				}
				if !strings.Contains(strings.ToLower(string(data)), lowerNeedle) {
					continue
				}
			}

			seen[match] = struct{}{}
			command := launchCmd
			if command == "" {
				command = match
			}
			variants = append(variants, model.Variant{
				Type:      variantType,
				Path:      match,
				LaunchCmd: command,
			})
		}
	}

	return variants
}

func matchingPaths(pattern string) []string {
	if strings.ContainsAny(pattern, "*?[") {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil
		}
		return matches
	}

	if _, err := os.Stat(pattern); err != nil {
		return nil
	}

	return []string{pattern}
}

func execLookPath(file string) (string, error) {
	return exec.LookPath(file)
}

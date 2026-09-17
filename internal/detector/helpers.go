package detector

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

var excludedDesktopPatterns = []string{
	"-url-handler",
	"-uninstall",
	"-devel",
	"kcm_",
	"-settings",
	"-preferences",
}

type desktopVariantCandidate struct {
	variant model.Variant
	target  string
	score   int
}

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

func collectDesktopVariants(patterns []string, canonicalName, variantType string) []model.Variant {
	candidates := make([]desktopVariantCandidate, 0)
	indexByTarget := make(map[string]int)
	canonicalLower := strings.ToLower(strings.TrimSpace(canonicalName))

	for _, pattern := range patterns {
		for _, match := range matchingPaths(pattern) {
			baseName := strings.ToLower(filepath.Base(match))
			if isExcludedDesktopFile(baseName) {
				continue
			}

			candidate, ok := readDesktopCandidate(match, canonicalLower, variantType)
			if !ok {
				continue
			}

			if existingIndex, found := indexByTarget[candidate.target]; found {
				if candidate.score < candidates[existingIndex].score {
					candidates[existingIndex] = candidate
				}
				continue
			}

			indexByTarget[candidate.target] = len(candidates)
			candidates = append(candidates, candidate)
		}
	}

	variants := make([]model.Variant, 0, len(candidates))
	for _, candidate := range candidates {
		variants = append(variants, candidate.variant)
	}

	return variants
}

func isExcludedDesktopFile(baseName string) bool {
	for _, pattern := range excludedDesktopPatterns {
		if strings.Contains(baseName, pattern) {
			return true
		}
	}
	return false
}

func readDesktopCandidate(path, canonicalName, variantType string) (desktopVariantCandidate, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		if !desktopFilenameMatches(path, canonicalName) {
			return desktopVariantCandidate{}, false
		}
		return desktopVariantCandidate{
			variant: model.Variant{Type: variantType, Path: path, LaunchCmd: path},
			target:  path,
			score:   desktopFilenameScore(path, canonicalName),
		}, true
	}

	name, hasName := desktopNameField(data)
	if hasName {
		if !desktopNameMatches(name, canonicalName) {
			return desktopVariantCandidate{}, false
		}
	} else if !desktopFilenameMatches(path, canonicalName) {
		return desktopVariantCandidate{}, false
	}

	execTarget, execCommand := desktopExecField(data)
	if execTarget == "" {
		execTarget = path
	}
	if execCommand == "" {
		execCommand = execTarget
	}

	return desktopVariantCandidate{
		variant: model.Variant{Type: variantType, Path: path, LaunchCmd: execCommand},
		target:  execTarget,
		score:   desktopFilenameScore(path, canonicalName),
	}, true
}

func desktopNameField(data []byte) (string, bool) {
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) < 6 {
			continue
		}

		lower := strings.ToLower(trimmed)
		if !strings.HasPrefix(lower, "name") {
			continue
		}

		equals := strings.Index(trimmed, "=")
		if equals < 0 {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(trimmed[:equals]))
		if key != "name" && !strings.HasPrefix(key, "name[") {
			continue
		}

		value := strings.TrimSpace(trimmed[equals+1:])
		if value == "" {
			continue
		}
		return value, true
	}

	return "", false
}

func desktopExecField(data []byte) (string, string) {
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(strings.ToLower(trimmed), "exec=") {
			continue
		}

		command := strings.TrimSpace(trimmed[len("Exec="):])
		if command == "" {
			return "", ""
		}

		tokens := splitDesktopCommand(command)
		for i := 0; i < len(tokens); i++ {
			token := strings.TrimSpace(tokens[i])
			if token == "" || strings.HasPrefix(token, "%") {
				continue
			}
			if token == "env" {
				continue
			}
			if strings.Contains(token, "=") && !strings.ContainsAny(token, "/\\") {
				continue
			}

			resolved := resolveDesktopTarget(token)
			return resolved, token
		}

		return "", command
	}

	return "", ""
}

func splitDesktopCommand(command string) []string {
	tokens := make([]string, 0)
	var current strings.Builder
	inSingle := false
	inDouble := false
	escape := false

	for _, r := range command {
		switch {
		case escape:
			current.WriteRune(r)
			escape = false
		case r == '\\' && !inSingle:
			escape = true
		case r == '\'' && !inDouble:
			inSingle = !inSingle
		case r == '"' && !inSingle:
			inDouble = !inDouble
		case unicode.IsSpace(r) && !inSingle && !inDouble:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func resolveDesktopTarget(token string) string {
	if token == "" {
		return ""
	}

	if filepath.IsAbs(token) || strings.ContainsAny(token, string(os.PathSeparator)+"/") {
		if resolved, err := filepath.EvalSymlinks(token); err == nil {
			return resolved
		}
		return filepath.Clean(token)
	}

	if resolved, err := exec.LookPath(token); err == nil {
		return resolved
	}

	return token
}

func desktopFilenameMatches(path, canonicalName string) bool {
	base := strings.TrimSuffix(strings.ToLower(filepath.Base(path)), ".desktop")
	canonical := normalizeDesktopName(canonicalName)
	return base == canonical || strings.Contains(base, canonical)
}

func desktopNameMatches(name, canonicalName string) bool {
	value := normalizeDesktopName(name)
	canonical := normalizeDesktopName(canonicalName)
	return value == canonical || strings.HasPrefix(value, canonical)
}

func normalizeDesktopName(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func desktopFilenameScore(path, canonicalName string) int {
	base := strings.TrimSuffix(strings.ToLower(filepath.Base(path)), ".desktop")
	score := len(base)

	if base == normalizeDesktopName(canonicalName) {
		score -= 10
	}

	for _, suffix := range excludedDesktopPatterns {
		if strings.Contains(base, suffix) {
			score += 100
		}
	}

	if strings.HasPrefix(base, "kcm_") {
		score += 100
	}

	if strings.Contains(base, "url-handler") {
		score += 50
	}

	return score
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

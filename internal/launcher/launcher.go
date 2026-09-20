package launcher

import (
	"errors"
	"fmt"
	"strings"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

func Launch(v model.Variant, mode string) error {
	t := strings.ToLower(strings.TrimSpace(v.Type))
	switch t {
	case "cli":
		if mode != "new_terminal" && mode != "same_terminal" {
			return fmt.Errorf("unsupported launch mode %q for cli variant", mode)
		}
		return launchCLI(v, mode)
	case "ide", "app":
		// IDEs and apps always launch detached; the mode is intentionally ignored.
		return launchDetached(v)
	default:
		return fmt.Errorf("unsupported variant type %q", v.Type)
	}
}

func parseCommandLine(command string) (string, []string, error) {
	parts := splitCommandLine(command)
	if len(parts) == 0 {
		return "", nil, errors.New("empty launch command")
	}

	return parts[0], parts[1:], nil
}

func splitCommandLine(command string) []string {
	parts := make([]string, 0)
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
		case r == ' ' || r == '\t' || r == '\n' || r == '\r':
			if inSingle || inDouble {
				current.WriteRune(r)
				continue
			}
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

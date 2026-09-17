package launcher

import (
	"fmt"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

func Launch(v model.Variant) error {
	fmt.Printf("launch stub: %s (%s)\n", v.LaunchCmd, v.LaunchMode)
	return nil
}

package config

import (
	"os"
	"testing"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

func TestSaveNeverWritesWithoutPath(t *testing.T) {
	cfg := NewConfigForTest([]model.Tool{{Name: "X"}}, "Violet Night")
	// Record state of real config path (if any) and ensure Save() does not
	// modify it. This makes the test safe even when a real user config exists.
	path, err := configPath()
	if err != nil {
		t.Fatalf("configPath failed: %v", err)
	}
	beforeInfo, beforeErr := os.Stat(path)

	// Ensure Save is a no-op (does not error)
	if err := cfg.Save(); err != nil {
		t.Fatalf("expected Save to be no-op for test config, got err: %v", err)
	}

	afterInfo, afterErr := os.Stat(path)

	if beforeErr == nil {
		// file existed before; it should still exist and have unchanged modtime
		if afterErr != nil {
			t.Fatalf("expected existing config file to remain, but it was removed: %v", afterErr)
		}
		if !beforeInfo.ModTime().Equal(afterInfo.ModTime()) {
			t.Fatalf("expected existing config file modtime to be unchanged, before=%v after=%v", beforeInfo.ModTime(), afterInfo.ModTime())
		}
	} else {
		// file did not exist before; it must still not exist
		if afterErr == nil {
			t.Fatalf("expected Save() to not create real config file at %s", path)
		}
	}
}

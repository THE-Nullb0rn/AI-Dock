package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/THE-Nullb0rn/AI-Dock/internal/model"
)

type Config struct {
	Tools []model.Tool `json:"tools"`
	Theme string       `json:"theme,omitempty"`
	// path is the resolved file path used for Save(). It is set by LoadConfig
	// and left empty for test-constructed configs to prevent accidental writes.
	path string
}

func configPath() (string, error) {
	if env := os.Getenv("AIDOCK_CONFIG"); env != "" {
		return env, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("find home directory: %w", err)
	}
	return filepath.Join(home, ".config", "aidock", "tools.json"), nil
}

func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		cfg := &Config{Tools: []model.Tool{}, Theme: "Violet Night"}
		cfg.path = path
		if err := cfg.Save(); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	// store the resolved path so subsequent Save() calls use the same file
	cfg.path = path
	if cfg.Tools == nil {
		cfg.Tools = []model.Tool{}
	}
	// Default theme for existing configs that predate the theme field.
	if cfg.Theme == "" {
		cfg.Theme = "Violet Night"
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	// If path is empty, this is likely a test-constructed Config — do not write.
	if c.path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(c.path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

// NewConfigForTest constructs a Config intended for tests. The returned
// Config has an empty internal path so calls to Save() become a no-op and
// cannot overwrite the user's real config file.
func NewConfigForTest(tools []model.Tool, theme string) *Config {
	if tools == nil {
		tools = []model.Tool{}
	}
	if theme == "" {
		theme = "Violet Night"
	}
	return &Config{Tools: tools, Theme: theme, path: ""}
}

// Test helper: returns whether this config would write to disk on Save().
// Not exported; used by tests below if needed.
func (c *Config) hasPath() bool {
	return c.path != ""
}

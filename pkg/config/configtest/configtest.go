package configtest

import (
	"path/filepath"
	"testing"

	"github.com/Darkclainer/japwords/pkg/config"
)

// New creates config.Manager with specified config (by making temporary dir with config file)
func New(tb testing.TB, userConfig *config.UserConfig) *config.Manager {
	tempDir := tb.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	err := config.SaveConfig(configPath, userConfig)
	if err != nil {
		tb.Fatalf("save config failed: %s", err)
	}
	manager, err := config.New(configPath)
	if err != nil {
		tb.Fatalf("manager creation failed: %s", err)
	}
	return manager
}

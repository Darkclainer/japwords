package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Darkclainer/japwords/pkg/anki"
	"github.com/Darkclainer/japwords/pkg/config"
)

// Test_Anki_ConfigReloader_DefaultConfig tests that default configuration for anki is valid
func Test_Anki_ConfigReloader_DefaultConfig(t *testing.T) {
	reloader := &anki.ConfigReloader{}
	_, err := reloader.Config(config.DefaultUserConfig())
	require.NoError(t, err)
}

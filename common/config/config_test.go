package config_test

import (
	"salami/common/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigGetters(t *testing.T) {
	t.Run("GetSourceDir", func(t *testing.T) {
		setConfigFile(t, "valid.yaml")
		config.LoadConfig()
		require.Equal(t, "testdata/source_dir", config.GetSourceDir())
	})

	t.Run("GetTargetDir", func(t *testing.T) {
		setConfigFile(t, "valid.yaml")
		config.LoadConfig()
		require.Equal(t, "terraform", config.GetTargetDir())
	})
}

package config_test

import (
	"path/filepath"
	"salami/common/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigGetters(t *testing.T) {
	t.Run("GetSourceDir", func(t *testing.T) {
		setConfigFile(t, "valid.yaml")
		config.LoadConfig()

		expectedSourceDir, err := filepath.Abs("testdata/source_dir")
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, expectedSourceDir, config.GetSourceDir())
	})

	t.Run("GetTargetDir", func(t *testing.T) {
		setConfigFile(t, "valid.yaml")
		config.LoadConfig()

		expectedTargetDir, err := filepath.Abs("terraform")
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, expectedTargetDir, config.GetTargetDir())
	})

	
}

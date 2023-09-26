package config_test

import (
	"path/filepath"
	"salami/common/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateConfig(t *testing.T) {
	testCases := getTestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setConfigFile(t, tc.fileName)
			err := config.ValidateConfig()
			require.Equal(t, err != nil, tc.wantErr, "unexpected error status: got error = %v, wantErr %v", err, tc.wantErr)
			if err != nil {
				require.Equal(
					t,
					tc.expectedErrorMessage,
					err.Error(),
					"unexpected error message: got = %v, want = %v",
					err.Error(),
					tc.expectedErrorMessage,
				)
			}
		})
	}
}

func setConfigFile(t *testing.T, fileName string) {
	fixturePath := filepath.Join("testdata", "config_files", fileName)
	config.SetConfigFilePath(fixturePath)
}

type testCase struct {
	name                 string
	fileName             string
	wantErr              bool
	expectedErrorMessage string
}

func getTestCases() []testCase {
	return []testCase{
		{
			"Valid config with all required fields",
			"valid.yaml",
			false,
			"",
		},
		{
			"Non-existing source directory",
			"invalid_source_dir.yaml",
			true,
			"config error: 'testdata/non_existent_dir' directory could not be resolved",
		},
		{
			"Missing target",
			"missing_target.yaml",
			true,
			"config error: invalid target configuration",
		},
		{
			"Invalid target platform",
			"invalid_platform.yaml",
			true,
			"config error: invalid target configuration",
		},
		{
			"Missing target platform",
			"missing_target_platform.yaml",
			true,
			"config error: invalid target configuration",
		},
		{
			"Missing llm",
			"missing_llm.yaml",
			true,
			"config error: invalid llm configuration",
		},
		{
			"Invalid llm provider",
			"invalid_llm_provider.yaml",
			true,
			"config error: invalid llm configuration",
		},
		{
			"Missing llm provider",
			"missing_llm_provider.yaml",
			true,
			"config error: invalid llm configuration",
		},
		{
			"Invalid llm model",
			"invalid_llm_model.yaml",
			true,
			"config error: invalid llm configuration",
		},
		{
			"Missing llm model",
			"missing_llm_model.yaml",
			true,
			"config error: invalid llm configuration",
		},
		{
			"Missing llm api key",
			"missing_llm_api_key.yaml",
			true,
			"config error: invalid llm configuration",
		},
		{
			"Invalid yaml format",
			"invalid_yaml.yaml",
			true,
			"config error: could not parse config file. Ensure it is valid yaml format",
		},
		{
			"Target directory outside of program's root directory",
			"target_dir_outside_root.yaml",
			true,
			"config error: target directory must be a subdirectory inside the root of the project",
		},
		{
			"Target directory equals the program's root directory",
			"target_dir_equals_root.yaml",
			true,
			"config error: target directory must be a subdirectory inside the root of the project",
		},
	}
}

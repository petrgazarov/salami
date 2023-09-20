package lock_file_manager_test

import (
	"salami/common/lock_file_manager"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLockFileValidate(t *testing.T) {
	testCases := getValidateTestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setLockFile(t, tc.fileName)
			err := lock_file_manager.ValidateLockFile()
			require.Equal(t, err != nil, tc.wantErr, "unexpected error status: got error = %v, wantErr %v", err, tc.wantErr)
			if err != nil {
				require.Equal(
					t,
					err.Error(),
					tc.expectedErrorMessage,
					"unexpected error message: got = %v, want = %v",
					err.Error(),
					tc.expectedErrorMessage,
				)
			}
		})
	}
}

type validateTestCase struct {
	name                 string
	fileName             string
	wantErr              bool
	expectedErrorMessage string
}

func getValidateTestCases() []validateTestCase {
	return []validateTestCase{
		{
			"Valid lock file",
			"valid.toml",
			false,
			"",
		},
		{
			"Empty lock file",
			"empty.toml",
			false,
			"",
		},
		{
			"Lock file doesn't exist",
			"nonexistent.toml",
			false,
			"",
		},
		{
			"Missing version",
			"missing_version.toml",
			true,
			"lock file error: missing or invalid lock file version",
		},
		{
			"Invalid semver",
			"invalid_semver.toml",
			true,
			"lock file error: '0.0.1a' is not a valid semver",
		},
		{
			"Missing target file path",
			"missing_target_file_path.toml",
			true,
			"lock file error: missing or invalid target file path",
		},
		{
			"Missing target file checksum",
			"missing_target_file_checksum.toml",
			true,
			"lock file error: missing or invalid target file checksum",
		},
		{
			"Missing source file path",
			"missing_source_file_path.toml",
			true,
			"lock file error: missing or invalid source file path",
		},
		{
			"Missing parsed resource object type",
			"missing_object_type.toml",
			true,
			"lock file error: missing or invalid parsed object type",
		},
		{
			"Missing parsed resource type",
			"missing_resource_resource_type.toml",
			true,
			"lock file error: missing or invalid parsed object resource type",
		},
		{
			"Missing parsed resource logical name",
			"missing_resource_logical_name.toml",
			true,
			"lock file error: missing or invalid parsed object logical name",
		},
		{
			"Missing variable name",
			"missing_variable_name.toml",
			true,
			"lock file error: missing or invalid parsed object name",
		},
		{
			"Invalid variable type",
			"missing_variable_type.toml",
			true,
			"lock file error: missing or invalid parsed object variable type",
		},
		{
			"Invalid variable type",
			"invalid_variable_type.toml",
			true,
			"lock file error: 'unsupported' is not a valid value",
		},
		{
			"Missing code segment type",
			"missing_segment_type.toml",
			true,
			"lock file error: missing or invalid code segment type",
		},
		{
			"Missing code segment target file path",
			"missing_code_segment_target_file_path.toml",
			true,
			"lock file error: missing or invalid code segment target file path",
		},
		{
			"Missing code segment content",
			"missing_code_segment_content.toml",
			true,
			"lock file error: missing or invalid code segment content",
		},
	}
}

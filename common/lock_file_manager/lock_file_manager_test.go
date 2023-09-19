package lock_file_manager_test

import (
	"path/filepath"
	"salami/common/lock_file_manager"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLockFileValidate(t *testing.T) {
	testCases := getTestCases()

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

func setLockFile(t *testing.T, fileName string) {
	fixturePath := filepath.Join("testdata", "lock_files", fileName)
	lock_file_manager.SetLockFilePath(fixturePath)
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
	}
}

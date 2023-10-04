package target_file_manager_test

import (
	"salami/backend/target_file_manager"
	"salami/common/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifyChecksums(t *testing.T) {
	testCases := []struct {
		name                 string
		filesAndChecksums    []types.TargetFileMeta
		wantErr              bool
		expectedErrorMessage string
	}{
		{
			"When checksums match",
			[]types.TargetFileMeta{
				{FilePath: "target_file_1.tf", Checksum: "415d4b5a48f2a887fffc285382fc4db1"},
				{FilePath: "target_file_2.tf", Checksum: "6700c5970a3183c2ecdc06900f7b30d4"}},
			false,
			"",
		},
		{
			"When one of the checksums doesn't match",
			[]types.TargetFileMeta{
				{FilePath: "target_file_1.tf", Checksum: "415d4b5a48f2a887fffc285382fc4db1"},
				{FilePath: "target_file_2.tf", Checksum: "invalid_checksum"}},
			true,
			"target file error: checksum mismatch for file target_file_2.tf",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := target_file_manager.VerifyChecksums(tc.filesAndChecksums, "testdata")
			require.Equal(t, err != nil, tc.wantErr, "unexpected error status: got error = %v, wantErr %v", err, tc.wantErr)
			if err != nil {
				require.Equal(
					t,
					err[0].Error(),
					tc.expectedErrorMessage,
					"unexpected error message: got = %v, want = %v",
					err[0].Error(),
					tc.expectedErrorMessage,
				)
			}
		})
	}
}

func TestGenerateTargetFileMetas(t *testing.T) {
	targetFiles := []*types.TargetFile{
		{FilePath: "testdata/target_file_1.tf", Content: "some content123"},
		{FilePath: "testdata/target_file_2.tf", Content: "another content456"},
	}
	expectedFileMetas := []types.TargetFileMeta{
		{FilePath: "testdata/target_file_1.tf", Checksum: "7fb90ebc6a8f51aedc1568b9f709ddf0"},
		{FilePath: "testdata/target_file_2.tf", Checksum: "d673f200b33c4c5f92bd7d1a1ca3b27f"},
	}

	t.Run("should compute checksums and return TargetFileMetas", func(t *testing.T) {
		fileMetas := target_file_manager.GenerateTargetFileMetas(targetFiles)
		require.Equal(t, expectedFileMetas, fileMetas, "unexpected file metas: got = %v, want = %v", fileMetas, expectedFileMetas)
	})
}

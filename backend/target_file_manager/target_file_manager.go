package target_file_manager

import "salami/common/types"

func VerifyChecksums(targetFilesMeta []types.TargetFileMeta) error {
	// Verify checksums
	return nil
}

func GenerateTargetFilesMeta([]*types.TargetFile) ([]types.TargetFileMeta, error) {
	// Compute checksums
	return []types.TargetFileMeta{}, nil
}

func WriteTargetFile(targetFile *types.TargetFile) error {
	// Write target file
	return nil
}

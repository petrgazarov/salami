package target_file_manager

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"

	"salami/common/types"

	"golang.org/x/sync/errgroup"
)

func VerifyChecksums(targetFileMetas []types.TargetFileMeta) error {
	var g errgroup.Group

	for _, meta := range targetFileMetas {
		meta := meta
		g.Go(func() error {
			data, err := os.ReadFile(meta.FilePath)
			if err != nil {
				return err
			}

			md5Checksum := fmt.Sprintf("%x", md5.Sum(data))
			if md5Checksum != meta.Checksum {
				return &TargetFileError{
					Message: fmt.Sprintf("checksum mismatch for file %s", meta.FilePath),
				}
			}

			return nil
		})
	}

	return g.Wait()
}

func GenerateTargetFileMetas(targetFiles []*types.TargetFile) ([]types.TargetFileMeta, error) {
	var g errgroup.Group
	targetFileMetas := make([]types.TargetFileMeta, len(targetFiles))

	for i, targetFile := range targetFiles {
		i, targetFile := i, targetFile
		g.Go(func() error {
			data := []byte(targetFile.Content)

			md5Checksum := fmt.Sprintf("%x", md5.Sum(data))
			targetFileMetas[i] = types.TargetFileMeta{
				FilePath: targetFile.FilePath,
				Checksum: md5Checksum,
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return targetFileMetas, nil
}

func DeleteTargetFiles(filePaths []string, targetDir string) []error {
	var errs []error

	for _, filePath := range filePaths {
		fullRelativePath := filepath.Join(targetDir, filePath)
		err := os.Remove(fullRelativePath)
		if err != nil {
			if !os.IsNotExist(err) {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

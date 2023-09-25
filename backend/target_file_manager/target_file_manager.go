package target_file_manager

import (
	"crypto/md5"
	"fmt"
	"os"

	"salami/common/types"

	"golang.org/x/sync/errgroup"
)

func VerifyChecksums(targetFilesMeta []types.TargetFileMeta) error {
	var g errgroup.Group

	for _, meta := range targetFilesMeta {
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

func GenerateTargetFilesMeta(targetFiles []*types.TargetFile) ([]types.TargetFileMeta, error) {
	var g errgroup.Group
	targetFilesMeta := make([]types.TargetFileMeta, len(targetFiles))

	for i, targetFile := range targetFiles {
		i, targetFile := i, targetFile
		g.Go(func() error {
			data := []byte(targetFile.Content)

			md5Checksum := fmt.Sprintf("%x", md5.Sum(data))
			targetFilesMeta[i] = types.TargetFileMeta{
				FilePath: targetFile.FilePath,
				Checksum: md5Checksum,
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return targetFilesMeta, nil
}

package target_file_manager

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"salami/common/types"
)

func VerifyChecksums(targetFileMetas []types.TargetFileMeta, targetDir string) []error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(targetFileMetas))

	for _, meta := range targetFileMetas {
		meta := meta
		wg.Add(1)
		go func() {
			defer wg.Done()
			fullRelativePath := filepath.Join(targetDir, meta.FilePath)
			data, err := os.ReadFile(fullRelativePath)
			if err != nil {
				errChan <- err
				return
			}

			md5Checksum := fmt.Sprintf("%x", md5.Sum(data))
			if md5Checksum != meta.Checksum {
				errChan <- &TargetFileError{
					Message: fmt.Sprintf("checksum mismatch for file %s", meta.FilePath),
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	return errs
}

func GenerateTargetFileMetas(targetFiles []*types.TargetFile) []types.TargetFileMeta {
	var wg sync.WaitGroup
	targetFileMetas := make([]types.TargetFileMeta, len(targetFiles))

	for i, targetFile := range targetFiles {
		i, targetFile := i, targetFile
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := []byte(targetFile.Content)

			md5Checksum := fmt.Sprintf("%x", md5.Sum(data))
			targetFileMetas[i] = types.TargetFileMeta{
				FilePath: targetFile.FilePath,
				Checksum: md5Checksum,
			}
		}()
	}

	wg.Wait()

	return targetFileMetas
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

package target_file_manager

import (
	"os"
	"path/filepath"
	"salami/common/types"
	"salami/common/utils/file_utils"
	"sync"
)

func WriteTargetFiles(targetFiles []*types.TargetFile, targetDir string) []error {
	var wg sync.WaitGroup

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return []error{err}
		}
	}

	errorChannel := make(chan error, len(targetFiles))
	for _, targetFile := range targetFiles {
		targetFile := targetFile
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := writeTargetFile(targetFile, targetDir); err != nil {
				errorChannel <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	errors := make([]error, 0, len(targetFiles))
	for err := range errorChannel {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func writeTargetFile(targetFile *types.TargetFile, targetDir string) error {
	fullRelativeFilePath := filepath.Join(targetDir, targetFile.FilePath)
	dir := filepath.Dir(fullRelativeFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	if err := file_utils.WriteFileIfChanged(fullRelativeFilePath, targetFile.Content); err != nil {
		return err
	}

	return nil
}

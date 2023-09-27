package target_file_manager

import (
	"io"
	"os"
	"path/filepath"
	"salami/common/types"

	"golang.org/x/sync/errgroup"
)

func WriteTargetFiles(targetFiles []*types.TargetFile, targetDir string) []error {
	var g errgroup.Group

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return []error{err}
		}
	}

	errs := make([]error, 0, len(targetFiles))
	for _, targetFile := range targetFiles {
		targetFile := targetFile
		g.Go(func() error {
			if err := writeTargetFile(targetFile, targetDir); err != nil {
				errs = append(errs, err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errs
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

	file, err := os.OpenFile(fullRelativeFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	oldContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if string(oldContent) != targetFile.Content {
		file.Truncate(0)
		file.Seek(0, 0)
		_, err = file.WriteString(targetFile.Content)
		if err != nil {
			return err
		}
	}

	return nil
}

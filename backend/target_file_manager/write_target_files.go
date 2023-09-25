package target_file_manager

import (
	"io"
	"os"
	"path/filepath"
	"salami/common/types"
	"salami/common/utils/file_utils"
	"strings"

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
			if err := writeTargetFile(targetFile); err != nil {
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

	if err := deleteOldFiles(targetDir, targetFiles); err != nil {
		return []error{err}
	}

	return nil
}

func writeTargetFile(targetFile *types.TargetFile) error {
	dir := filepath.Dir(targetFile.FilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(targetFile.FilePath, os.O_RDWR|os.O_CREATE, 0644)
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

func deleteOldFiles(targetDir string, targetFiles []*types.TargetFile) error {
	oldFilePaths, err := getOldFilePaths(targetDir, targetFiles)
	if err != nil {
		return err
	}

	for _, file := range oldFilePaths {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	isAncestorOfAnyTargetFile := func(path string, targetFiles []*types.TargetFile) bool {
		for _, targetFile := range targetFiles {
			if strings.HasPrefix(targetFile.FilePath, path) {
				return true
			}
		}
		return false
	}

	emptySubdirectoryPaths, err := file_utils.GetSubdirectoryPaths(targetDir, func(path string) bool {
		return !isAncestorOfAnyTargetFile(path, targetFiles)
	})
	if err != nil {
		return err
	}

	for _, dir := range emptySubdirectoryPaths {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

func getOldFilePaths(targetDir string, targetFiles []*types.TargetFile) ([]string, error) {
	targetFilesMap := make(map[string]bool)
	for _, targetFile := range targetFiles {
		targetFilesMap[targetFile.FilePath] = true
	}

	filter := func(path string) bool {
		relativePath, err := filepath.Rel(targetDir, path)
		if err != nil {
			panic(err)
		}
		_, exists := targetFilesMap[relativePath]
		return !exists
	}

	oldFiles, err := file_utils.GetFilePaths(targetDir, filter)
	if err != nil {
		return nil, err
	}
	return oldFiles, nil
}

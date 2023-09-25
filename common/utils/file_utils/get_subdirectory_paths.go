package file_utils

import (
	"os"
	"path/filepath"
)

func GetSubdirectoryPaths(directory string, filter FilterFunc) ([]string, error) {
	var subdirectories []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != directory && filter(path) {
			subdirectories = append(subdirectories, path)
		}
		return nil
	})

	return subdirectories, err
}

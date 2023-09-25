package file_utils

import (
	"os"
	"path/filepath"
)

type FilterFunc func(string) bool

func GetFilePaths(directory string, filter FilterFunc) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filter(path) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

package file_utils

import "path/filepath"

func GetRelativeFilePaths(baseDir string, paths []string) ([]string, error) {
	relativePaths := make([]string, len(paths))
	for i, path := range paths {
		relativePath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return nil, err
		}
		relativePaths[i] = relativePath
	}
	return relativePaths, nil
}

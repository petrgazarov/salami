package target_file_manager

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

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

func WriteTargetFiles(targetFiles []*types.TargetFile, targetDir string) []error {
	var g errgroup.Group
	var mu sync.Mutex
	errs := make([]error, 0, len(targetFiles))

	for _, targetFile := range targetFiles {
		targetFile := targetFile
		g.Go(func() error {
			if err := writeTargetFile(targetFile); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
			return nil
		})
	}

	g.Go(func() error {
		if err := deleteOldFiles(targetDir, targetFiles); err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		mu.Lock()
		errs = append(errs, err)
		mu.Unlock()
	}

	return errs
}

func writeTargetFile(targetFile *types.TargetFile) error {
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
	targetFilesMap := make(map[string]bool)
	for _, targetFile := range targetFiles {
		targetFilesMap[targetFile.FilePath] = true
	}

	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if _, exists := targetFilesMap[path]; !exists {
			if info.IsDir() {
				return os.RemoveAll(path)
			} else {
				return os.Remove(path)
			}
		}
		return nil
	})

	return err
}

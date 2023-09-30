package file_utils

import (
	"io"
	"os"
)

func WriteFileIfChanged(fullRelativeFilePath string, newContent string) error {
	file, err := os.OpenFile(fullRelativeFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	oldContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if string(oldContent) != newContent {
		file.Truncate(0)
		file.Seek(0, 0)
		_, err = file.WriteString(newContent)
		if err != nil {
			return err
		}
	}

	return nil
}

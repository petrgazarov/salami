package lock_file_manager

import (
	"log"
	"os"
	commonTypes "salami/common/types"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
)

const lockFilePath = "salami-lock.toml"

var loadedLockFile *commonTypes.LockFile

func ValidateLockFile() error {
	lockFile := getLockFile()
	validate := validator.New()
	if err := validate.Struct(lockFile); err != nil {
		return err
	}
	for _, object := range lockFile.Objects {
		if err := object.Parsed.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func GetTargetFiles() []*commonTypes.TargetFile {
	targetFiles := getLockFile().TargetFiles
	ptrTargetFiles := make([]*commonTypes.TargetFile, len(targetFiles))
	for i := range targetFiles {
		ptrTargetFiles[i] = &targetFiles[i]
	}
	return ptrTargetFiles
}

func UpdateLockFile(changeSet *commonTypes.ChangeSet) error {
	// Merge changeSet into lockFile
	writeLockFile()
	return nil
}

func writeLockFile() error {
	file, err := os.OpenFile(lockFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(getLockFile()); err != nil {
		return err
	}
	return nil
}

func getLockFile() *commonTypes.LockFile {
	if loadedLockFile == nil {
		if _, err := toml.DecodeFile(lockFilePath, &loadedLockFile); err != nil {
			if err != nil && !os.IsNotExist(err) {
				log.Fatalf("failed to read lock file")
			}
		}
	}
	return loadedLockFile
}

package driver

import (
	"salami/backend/target_file_manager"
	"salami/common/config"
	"salami/common/lock_file_manager"
)

func runValidations() error {
	if err := config.ValidateConfig(); err != nil {
		return err
	}
	if err := lock_file_manager.ValidateLockFile(); err != nil {
		return err
	}
	targetFiles := lock_file_manager.GetTargetFiles()
	if err := target_file_manager.ValidateTargetFilesUnchanged(targetFiles); err != nil {
		return err
	}
	return nil
}

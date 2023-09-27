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
	targetFileMetas := lock_file_manager.GetTargetFileMetas()
	if err := target_file_manager.VerifyChecksums(targetFileMetas); err != nil {
		return err
	}
	return nil
}

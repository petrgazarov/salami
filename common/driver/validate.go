package driver

import (
	"salami/backend/target"
	"salami/backend/target_file_manager"
	"salami/common/config"
	"salami/common/lock_file_manager"
)

func runValidations() []error {
	if err := config.ValidateConfig(); err != nil {
		return []error{err}
	}

	if err := lock_file_manager.ValidateLockFile(); err != nil {
		return []error{err}
	}

	targetFileMetas := lock_file_manager.GetTargetFileMetas()
	targetDir := config.GetTargetDir()
	if errors := target_file_manager.VerifyChecksums(targetFileMetas, targetDir); len(errors) > 0 {
		return errors
	}

	target := target.ResolveTarget(config.GetTargetConfig())
	if err := target.VerifyPeerDependencies(); err != nil {
		return []error{err}
	}

	return nil
}

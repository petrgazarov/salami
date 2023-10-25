package driver

import (
	"fmt"
	"salami/common/lock_file_manager"
	"salami/common/logger"
	"salami/common/metrics"
	"strconv"
)

func Run(verbose bool) []error {
	logger.InitializeLogger(verbose)
	metrics.InitializeMetrics()

	if errors := runValidations(); len(errors) > 0 {
		return errors
	}

	symbolTable, errors := runFrontend()
	if len(errors) > 0 {
		return errors
	}

	newTargetFileMetas, newObjects, errors := runBackend(symbolTable)
	if len(errors) > 0 {
		return errors
	}

	if err := lock_file_manager.UpdateLockFile(newTargetFileMetas, newObjects); err != nil {
		return []error{err}
	}
	logCompletion()

	return nil
}

func logCompletion() {
	message := fmt.Sprintf(
		"✨ Done in %s. Changed objects: %s, added objects: %s, removed objects: %s, processed files: %s",
		metrics.GetDuration(),
		strconv.Itoa(metrics.GetMetrics().ObjectsChanged),
		strconv.Itoa(metrics.GetMetrics().ObjectsAdded),
		strconv.Itoa(metrics.GetMetrics().ObjectsRemoved),
		strconv.Itoa(metrics.GetMetrics().SourceFilesProcessed),
	)
	logger.Log(message)
}

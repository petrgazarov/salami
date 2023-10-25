package driver

import (
	"salami/backend/llm"
	"salami/backend/target"
	"salami/backend/target_file_manager"
	"salami/common/change_set"
	"salami/common/config"
	"salami/common/lock_file_manager"
	"salami/common/metrics"
	"salami/common/symbol_table"
	"salami/common/types"
	"sort"
)

func runBackend(
	symbolTable *symbol_table.SymbolTable,
) ([]types.TargetFileMeta, []*types.Object, []error) {
	newObjects, newTargetFiles, errors := generateCode(symbolTable)
	if len(errors) > 0 {
		return nil, nil, errors
	}

	targetDir := config.GetTargetDir()
	if errors := target_file_manager.WriteTargetFiles(newTargetFiles, targetDir); len(errors) > 0 {
		return nil, nil, errors
	}

	oldTargetFileMetas := lock_file_manager.GetTargetFileMetas()
	newTargetFileMetas := target_file_manager.GenerateTargetFileMetas(newTargetFiles)
	filePathsToRemove := getFilePathsToRemove(oldTargetFileMetas, newTargetFileMetas)
	if errors := target_file_manager.DeleteTargetFiles(filePathsToRemove, targetDir); len(errors) > 0 {
		return nil, nil, errors
	}

	metrics.SetMetric(metrics.SourceFilesProcessed, len(newTargetFiles))

	return newTargetFileMetas, newObjects, nil
}

func generateCode(
	symbolTable *symbol_table.SymbolTable,
) ([]*types.Object, []*types.TargetFile, []error) {
	previousObjects := lock_file_manager.GetObjects()
	changeSet := change_set.NewChangeSet(previousObjects, symbolTable)
	changeSetRepository := change_set.NewChangeSetRepository(changeSet)

	llm := llm.ResolveLlm(config.GetLlmConfig())
	target := target.ResolveTarget(config.GetTargetConfig())

	if errors := target.GenerateCode(symbolTable, changeSetRepository, llm); len(errors) > 0 {
		return nil, nil, errors
	}

	newObjects := computeNewObjects(previousObjects, changeSetRepository)
	if err := target.ValidateCode(newObjects, symbolTable, changeSetRepository, llm, 0); err != nil {
		return nil, nil, []error{err}
	}
	targetFiles := target.GetFilesFromObjects(newObjects)

	updateMetricsFromChangeSet(changeSet)

	return newObjects, targetFiles, nil
}

func getFilePathsToRemove(
	oldTargetFileMetas []types.TargetFileMeta,
	newTargetFileMetas []types.TargetFileMeta,
) []string {
	newMetaMap := make(map[string]bool)
	for _, meta := range newTargetFileMetas {
		newMetaMap[meta.FilePath] = true
	}

	oldFilePaths := make([]string, 0)
	for _, meta := range oldTargetFileMetas {
		if _, exists := newMetaMap[meta.FilePath]; !exists {
			oldFilePaths = append(oldFilePaths, meta.FilePath)
		}
	}

	return oldFilePaths
}

func computeNewObjects(
	previousObjects []*types.Object,
	changeSetRepository *change_set.ChangeSetRepository,
) []*types.Object {
	objects := make([]*types.Object, 0)
	for _, object := range previousObjects {
		if changeSetRepository.WasObjectChanged(object) {
			objects = append(objects, changeSetRepository.GetChangedObject(object))
			continue
		}
		if !changeSetRepository.WasObjectDeleted(object) {
			objects = append(objects, object)
		}
	}
	objects = append(objects, changeSetRepository.AddedObjects...)

	sort.Slice(objects, func(i, j int) bool {
		currentObject := objects[i]
		nextObject := objects[j]

		if currentObject.GetSourceFilePath() != nextObject.GetSourceFilePath() {
			return currentObject.GetSourceFilePath() < nextObject.GetSourceFilePath()
		} else {
			return currentObject.GetSourceFileLine() < nextObject.GetSourceFileLine()
		}
	})

	return objects
}

func updateMetricsFromChangeSet(changeSet *types.ChangeSet) {
	objectsAdded := 0
	objectsRemoved := 0
	objectsChanged := 0

	for _, diff := range changeSet.Diffs {
		if diff.IsAdd() {
			objectsAdded++
		} else if diff.IsRemove() {
			objectsRemoved++
		} else if diff.IsUpdate() || diff.IsMove() {
			objectsChanged++
		}
	}

	metrics.SetMetric(metrics.ObjectsAdded, objectsAdded)
	metrics.SetMetric(metrics.ObjectsRemoved, objectsRemoved)
	metrics.SetMetric(metrics.ObjectsChanged, objectsChanged)
}

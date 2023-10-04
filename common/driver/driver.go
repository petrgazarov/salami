package driver

import (
	"path/filepath"
	"sort"

	"salami/backend/llm"
	"salami/backend/target"
	"salami/backend/target_file_manager"
	"salami/common/change_set"
	"salami/common/config"
	"salami/common/constants"
	"salami/common/lock_file_manager"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"
	"salami/common/utils/file_utils"
	"salami/frontend/semantic_analyzer"
)

func Run() []error {
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

	return nil
}

func runFrontend() (*symbol_table.SymbolTable, []error) {
	sourceFilePaths, err := getSourceFilePaths()
	if err != nil {
		return nil, []error{err}
	}

	allResources, allVariables, errors := parseFiles(sourceFilePaths, config.GetSourceDir())
	if len(errors) > 0 {
		return nil, errors
	}

	symbolTable, err := symbol_table.NewSymbolTable(allResources, allVariables)
	if err != nil {
		return nil, []error{err}
	}

	semanticAnalyzer := semantic_analyzer.NewSemanticAnalyzer(symbolTable)
	if err = semanticAnalyzer.Analyze(); err != nil {
		return nil, []error{err}
	}

	return symbolTable, nil
}

func runBackend(
	symbolTable *symbol_table.SymbolTable,
) ([]commonTypes.TargetFileMeta, []*commonTypes.Object, []error) {
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

	return newTargetFileMetas, newObjects, nil
}

func generateCode(
	symbolTable *symbol_table.SymbolTable,
) ([]*commonTypes.Object, []*commonTypes.TargetFile, []error) {
	previousObjects := lock_file_manager.GetObjects()
	changeSet := change_set.NewChangeSet(previousObjects, symbolTable)
	changeSetRepository := change_set.NewChangeSetRepository(changeSet)

	llm := llm.ResolveLlm(config.GetLlmConfig())
	target := target.ResolveTarget(config.GetTargetConfig())

	if errors := target.GenerateCode(symbolTable, changeSetRepository, llm); len(errors) > 0 {
		return nil, nil, errors
	}

	newObjects := computeNewObjects(previousObjects, changeSetRepository)
	if errors := target.ValidateCode(newObjects, symbolTable, changeSetRepository, llm); len(errors) > 0 {
		return nil, nil, errors
	}
	targetFiles := target.GetFilesFromObjects(newObjects)

	return newObjects, targetFiles, nil
}

func getSourceFilePaths() ([]string, error) {
	sourceFilePaths, err := file_utils.GetFilePaths(config.GetSourceDir(), func(path string) bool {
		return filepath.Ext(path) == constants.SalamiFileExtension
	})
	if err != nil {
		return nil, err
	}

	relativeSourceFilePaths, err := file_utils.GetRelativeFilePaths(
		config.GetSourceDir(),
		sourceFilePaths,
	)
	if err != nil {
		return nil, err
	}

	return relativeSourceFilePaths, nil
}

func computeNewObjects(
	previousObjects []*commonTypes.Object,
	changeSetRepository *change_set.ChangeSetRepository,
) []*commonTypes.Object {
	objects := make([]*commonTypes.Object, 0)
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

func getFilePathsToRemove(
	oldTargetFileMetas []commonTypes.TargetFileMeta,
	newTargetFileMetas []commonTypes.TargetFileMeta,
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

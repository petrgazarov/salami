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
	if err := runValidations(); err != nil {
		return []error{err}
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
	previousObjects := lock_file_manager.GetObjects()
	targetDir := config.GetTargetDir()
	newTargetFiles, newObjects, errors := generateCode(previousObjects, symbolTable)
	if len(errors) > 0 {
		return nil, nil, errors
	}
	if errors := target_file_manager.WriteTargetFiles(newTargetFiles, targetDir); len(errors) > 0 {
		return nil, nil, errors
	}

	newTargetFileMetas, err := target_file_manager.GenerateTargetFileMetas(newTargetFiles)
	if err != nil {
		return nil, nil, []error{err}
	}

	oldFilePaths := getOldFilePaths(newTargetFileMetas)
	if errors := target_file_manager.DeleteTargetFiles(oldFilePaths, targetDir); len(errors) > 0 {
		return nil, nil, errors
	}

	return newTargetFileMetas, newObjects, nil
}

func generateCode(
	previousObjects []*commonTypes.Object,
	symbolTable *symbol_table.SymbolTable,
) ([]*commonTypes.TargetFile, []*commonTypes.Object, []error) {
	changeSet := change_set.NewChangeSet(previousObjects, symbolTable)
	llm := llm.ResolveLlm(config.GetLlmConfig())
	target := target.ResolveTarget(config.GetTargetConfig())

	if errors := target.GenerateCode(changeSet, symbolTable, llm); len(errors) > 0 {
		return nil, nil, errors
	}

	newObjects := computeNewObjects(previousObjects, changeSet)
	targetFiles := target.GetFilesFromObjects(newObjects)

	return targetFiles, newObjects, nil
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
	changeSet *commonTypes.ChangeSet,
) []*commonTypes.Object {
	changeSetRepository := change_set.NewChangeSetRepository(changeSet)

	objects := make([]*commonTypes.Object, 0)
	for _, object := range previousObjects {
		if changeSetRepository.WasChanged(object) {
			objects = append(objects, changeSetRepository.GetChanged(object))
			continue
		}
		if !changeSetRepository.WasDeleted(object) {
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

func getOldFilePaths(newTargetFileMetas []commonTypes.TargetFileMeta) []string {
	oldTargetFileMetas := lock_file_manager.GetTargetFileMetas()
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

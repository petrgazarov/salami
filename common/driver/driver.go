package driver

import (
	"os"
	"path/filepath"
	"sync"

	"salami/backend/target"
	"salami/backend/target_file_manager"
	"salami/common/change_manager"
	"salami/common/config"
	"salami/common/lock_file_manager"
	"salami/common/symbol_table"
	"salami/common/types"
	"salami/frontend/semantic_analyzer"
)

const salamiFileExtension = ".sami"

func Run() []error {
	if err := runValidations(); err != nil {
		return []error{err}
	}
	symbolTable, errors := runFrontend()
	if len(errors) > 0 {
		return errors
	}
	newTargetFiles, objects, errors := generateCode(symbolTable)
	if len(errors) > 0 {
		return errors
	}
	if errors := writeTargetFiles(newTargetFiles); len(errors) > 0 {
		return errors
	}

	newTargetFilesMeta, err := target_file_manager.GenerateTargetFilesMeta(newTargetFiles)
	if err != nil {
		return []error{err}
	}
	if err := lock_file_manager.UpdateLockFile(newTargetFilesMeta, objects); err != nil {
		return []error{err}
	}
	return nil
}

func runFrontend() (*symbol_table.SymbolTable, []error) {
	sourceFilePaths, err := getSourceFilePaths()
	if err != nil {
		return nil, []error{err}
	}
	allResources, allVariables, errors := parseFiles(sourceFilePaths)
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

func getSourceFilePaths() ([]string, error) {
	sourceDir := config.GetSourceDir()
	var files []string

	error := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == salamiFileExtension {
			files = append(files, path)
		}
		return nil
	})

	return files, error
}

func generateCode(
	symbolTable *symbol_table.SymbolTable,
) ([]*types.TargetFile, []*types.Object, []error) {
	previousObjects := lock_file_manager.GetObjects()
	changeSet, err := change_manager.GenerateChangeSet(previousObjects, symbolTable)
	if err != nil {
		return nil, nil, []error{err}
	}
	targetConfig := config.GetTargetConfig()
	llmConfig := config.GetLlmConfig()
	resolvedTarget := target.ResolveTarget(targetConfig, llmConfig)
	if errors := resolvedTarget.GenerateCode(changeSet, symbolTable); len(errors) > 0 {
		return nil, nil, errors
	}
	objects := change_manager.ComputeNewObjects(previousObjects, changeSet)
	targetFiles := resolvedTarget.GetNewFiles(objects)
	return targetFiles, objects, nil
}

func writeTargetFiles(targetFiles []*types.TargetFile) []error {
	errChan := make(chan error, len(targetFiles))
	var wg sync.WaitGroup

	for _, targetFile := range targetFiles {
		wg.Add(1)
		go func(tf *types.TargetFile) {
			defer wg.Done()
			err := target_file_manager.WriteTargetFile(tf)
			if err != nil {
				errChan <- err
			}
		}(targetFile)
	}

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

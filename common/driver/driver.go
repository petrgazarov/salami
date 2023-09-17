package driver

import (
	"os"
	"path/filepath"

	"salami/backend/target"
	"salami/backend/target_file_manager"
	"salami/common/change_manager"
	"salami/common/config"
	"salami/common/lock_file_manager"
	"salami/common/symbol_table"
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
	previousObjects := lock_file_manager.GetObjects()
	changeSet, err := change_manager.GenerateChangeSet(previousObjects, symbolTable)
	if err != nil {
		return []error{err}
	}
	targetConfig := config.GetTargetConfig()
	llmConfig := config.GetLlmConfig()
	resolvedTarget := target.ResolveTarget(targetConfig, llmConfig)
	if errors := resolvedTarget.GenerateCode(changeSet, symbolTable); len(errors) > 0 {
		return errors
	}
	newObjects := change_manager.ComputeNewObjects(changeSet, previousObjects)
	targetFiles := resolvedTarget.GetNewFiles(newObjects)
	for _, targetFile := range targetFiles {
		err := target_file_manager.WriteTargetFile(targetFile)
		if err != nil {
			return []error{err}
		}
	}
	newTargetFilesMeta, err := target_file_manager.GenerateTargetFilesMeta(targetFiles)
	if err := lock_file_manager.UpdateLockFile(newTargetFilesMeta, newObjects); err != nil {
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

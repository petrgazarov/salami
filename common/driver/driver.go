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
	if errors := target_file_manager.WriteTargetFiles(newTargetFiles, config.GetTargetDir()); len(errors) > 0 {
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
	previousResources, previousVariables := getPreviousObjectsMaps()
	changeSet := change_manager.GenerateChangeSet(previousResources, previousVariables, symbolTable)
	target := resolveTarget()
	if errors := target.GenerateCode(changeSet, symbolTable); len(errors) > 0 {
		return nil, nil, errors
	}
	newObjects := change_manager.ComputeNewObjects(previousResources, previousVariables, changeSet)
	targetFiles := target.GetFilesFromObjects(newObjects)
	return targetFiles, newObjects, nil
}

func getPreviousObjectsMaps() (
	map[types.LogicalName]*types.Object,
	map[string]*types.Object,
) {
	previousObjects := lock_file_manager.GetObjects()
	resources := make(map[types.LogicalName]*types.Object)
	variables := make(map[string]*types.Object)

	for _, object := range previousObjects {
		switch v := object.Parsed.(type) {
		case *types.ParsedResource:
			resources[v.LogicalName] = object
		case *types.ParsedVariable:
			variables[v.Name] = object
		}
	}

	return resources, variables
}

func resolveTarget() target.Target {
	targetConfig := config.GetTargetConfig()
	llmConfig := config.GetLlmConfig()
	return target.ResolveTarget(targetConfig, llmConfig)
}

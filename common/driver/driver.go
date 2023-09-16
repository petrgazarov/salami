package driver

import (
	"os"
	"path/filepath"

	"salami/backend"
	"salami/common/config"
	"salami/common/lock_file_manager"
	"salami/common/change_set_manager"
	"salami/common/symbol_table"
	"salami/common/types"
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
	changeSet, err := change_set_manager.GenerateChangeSet(symbolTable)
	if err != nil {
		return []error{err}
	}
	backend := backend.NewBackend(symbolTable, changeSet)
	if errors := backend.GenerateCode(); len(errors) > 0 {
		return errors
	}
	err := backend.WriteTargetFiles()
	if err != nil {
		return []error{err}
	}
	if err := lock_file_manager.UpdateLockFile(changeSet); err != nil {
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
	compilerConfig := config.GetCompilerConfig()
	var files []string

	error := filepath.Walk(compilerConfig.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == types.SalamiFileExtension {
			files = append(files, path)
		}
		return nil
	})

	return files, error
}

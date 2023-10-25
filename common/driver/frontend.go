package driver

import (
	"os"
	"path/filepath"
	"salami/common/config"
	"salami/common/constants"
	"salami/common/symbol_table"
	"salami/common/types"
	"salami/common/utils/file_utils"
	"salami/frontend/lexer"
	"salami/frontend/parser"
	"salami/frontend/semantic_analyzer"
	"sync"
)

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

func parseFiles(filePaths []string, sourceDir string) ([]*types.ParsedResource, []*types.ParsedVariable, []error) {
	resourcesChan := make(chan []*types.ParsedResource, len(filePaths))
	variablesChan := make(chan []*types.ParsedVariable, len(filePaths))
	errorChan := make(chan error, len(filePaths))

	var wg sync.WaitGroup
	for _, filePath := range filePaths {
		wg.Add(1)
		go func(filePath string) {
			resources, variables, err := parseFile(filePath, sourceDir)
			if err != nil {
				errorChan <- err
			}
			if resources != nil {
				resourcesChan <- resources
			}
			if variables != nil {
				variablesChan <- variables
			}
			wg.Done()
		}(filePath)
	}
	wg.Wait()
	close(resourcesChan)
	close(variablesChan)
	close(errorChan)

	var allResources []*types.ParsedResource
	var allVariables []*types.ParsedVariable
	var allErrors []error
	for resources := range resourcesChan {
		allResources = append(allResources, resources...)
	}
	for variables := range variablesChan {
		allVariables = append(allVariables, variables...)
	}
	for err := range errorChan {
		allErrors = append(allErrors, err)
	}

	return allResources, allVariables, allErrors
}

func parseFile(filePath string, sourceDir string) (
	resources []*types.ParsedResource,
	variables []*types.ParsedVariable,
	err error,
) {
	fullRelativeFilePath := filepath.Join(sourceDir, filePath)
	content, err := os.ReadFile(fullRelativeFilePath)
	if err != nil {
		return nil, nil, err
	}

	lexerInstance := lexer.NewLexer(filePath, string(content))
	tokens, err := lexerInstance.Run()

	if err != nil {
		return nil, nil, err
	}
	parserInstance := parser.NewParser(tokens, filePath)
	resources, variables, parsingError := parserInstance.Parse()
	return resources, variables, parsingError
}

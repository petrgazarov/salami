package driver

import (
	"os"
	"path/filepath"

	"salami/backend"
	"salami/common/config"
	"salami/common/symbol_table"
	"salami/common/types"
	"salami/frontend/lexer"
	"salami/frontend/parser"
	"salami/frontend/semantic_analyzer"
	"sync"
)

func Run() []error {
	config.ValidateConfig()
	files, err := getFilePaths()
	if err != nil {
		return []error{err}
	}

	allResources, allVariables, errors := parseFiles(files)

	if len(errors) > 0 {
		return errors
	}
	symbolTable, err := symbol_table.NewSymbolTable(allResources, allVariables)
	if err != nil {
		return []error{err}
	}
	semanticAnalyzer := semantic_analyzer.NewSemanticAnalyzer(symbolTable)
	if err = semanticAnalyzer.Analyze(); err != nil {
		return []error{err}
	}
	backend := backend.NewBackend(symbolTable)
	targetFilePaths, backendObjects, errors := backend.Generate()
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func getFilePaths() ([]string, error) {
	compilerConfig := config.GetCompilerConfig()
	var files []string

	error := filepath.Walk(compilerConfig.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".sami" {
			files = append(files, path)
		}
		return nil
	})

	return files, error
}

func parseFiles(files []string) ([]*types.Resource, []*types.Variable, []error) {
	resourcesChan := make(chan []*types.Resource, len(files))
	variablesChan := make(chan []*types.Variable, len(files))
	errorChan := make(chan error, len(files))

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			resources, variables, err := parseFile(file)
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
		}(file)
	}
	wg.Wait()
	close(resourcesChan)
	close(variablesChan)
	close(errorChan)

	var allResources []*types.Resource
	var allVariables []*types.Variable
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

func parseFile(filePath string) (
	resources []*types.Resource,
	variables []*types.Variable,
	err error,
) {
	content, err := os.ReadFile(filePath)
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

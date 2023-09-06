package driver

import (
	"os"
	"path/filepath"
	"salami/compiler/config"
	"salami/compiler/lexer"
	"salami/compiler/parser"
	"salami/compiler/semantic_analyzer"
	"salami/compiler/symbol_table"
	"salami/compiler/types"
	"sync"
)

func Run() []error {
	files, err := getFilePaths()
	if err != nil {
		return []error{err}
	}

	allResources, allVariables, allErrors := processFiles(files)

	if len(allErrors) > 0 {
		return allErrors
	}
	symbolTable := symbol_table.NewSymbolTable(allResources, allVariables)
	semanticAnalyzer := semantic_analyzer.NewSemanticAnalyzer(symbolTable)
	semanticAnalyzer.Analyze()
	return nil
}

func getFilePaths() ([]string, error) {
	compilerConfig := config.GetConfig()
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

func processFiles(files []string) ([]*types.Resource, []*types.Variable, []error) {
	resourcesChan := make(chan []*types.Resource, len(files))
	variablesChan := make(chan []*types.Variable, len(files))
	errorChan := make(chan error, len(files))

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			resources, variables, err := processFile(file)
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

func processFile(filePath string) (
	resources []*types.Resource,
	variables []*types.Variable,
	err error,
) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	lexerInstance := lexer.NewLexer(filePath, string(content))
	tokens := lexerInstance.Run()
	parserInstance := parser.NewParser(tokens, filePath)
	resources, variables, parsingError := parserInstance.Parse()
	return resources, variables, parsingError
}

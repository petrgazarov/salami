package driver

import (
	"os"
	"path/filepath"
	"salami/common/types"
	"salami/frontend/lexer"
	"salami/frontend/parser"
	"sync"
)

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

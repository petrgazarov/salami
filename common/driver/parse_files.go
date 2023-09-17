package driver

import (
	"os"
	"salami/common/types"
	"salami/frontend/lexer"
	"salami/frontend/parser"
	"sync"
)

func parseFiles(files []string) ([]*types.ParsedResource, []*types.ParsedVariable, []error) {
	resourcesChan := make(chan []*types.ParsedResource, len(files))
	variablesChan := make(chan []*types.ParsedVariable, len(files))
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

func parseFile(filePath string) (
	resources []*types.ParsedResource,
	variables []*types.ParsedVariable,
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

package terraform

import (
	"salami/backend/types"
	"sync"
)

type TerraformTarget struct {
}

func (t *TerraformTarget) GenerateCodeFiles(
	objectsMap map[string][]types.ParsedObject,
) ([]*types.GeneratedCodeFile, []error) {
	var wg sync.WaitGroup
	resultsChan := make(chan *types.DestinationFile)
	errorsChan := make(chan error)

	for _, objects := range objectsMap {
		wg.Add(1)
		go func(objs []types.ParsedObject) {
			defer wg.Done()
			result, err := t.generateCodeFile(objs)
			if err != nil {
				errorsChan <- err
			} else {
				resultsChan <- result
			}
		}(objects)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()

	var results []*types.DestinationFile
	var errors []error
	for result := range resultsChan {
		results = append(results, result)
	}
	for err := range errorsChan {
		errors = append(errors, err)
	}

	return results, errors
}

func (t *TerraformTarget) generateCodeFile(objects []types.ParsedObject) ([]*types.ObjectGeneratedCode, error) {
	return nil, nil
}

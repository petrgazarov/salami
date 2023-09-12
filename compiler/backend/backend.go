package backend

import (
	"os"
	"path"
	"salami/compiler/backend/target"
	"salami/compiler/config"
	"salami/compiler/symbol_table"
	"salami/compiler/types"
	"sort"
	"sync"
)

type Backend struct {
	symbolTable *symbol_table.SymbolTable
}

func NewBackend(symbolTable *symbol_table.SymbolTable) *Backend {
	return &Backend{
		symbolTable: symbolTable,
	}
}

func (b *Backend) Generate() []error {
	compilerConfig := config.GetCompilerConfig()
	targetModule, err := target.ResolveTarget(compilerConfig.Target, compilerConfig.Llm)
	if err != nil {
		return []error{err}
	}
	objectsMap := b.objectsMap()
	codeFiles, errors := targetModule.GenerateCodeFiles(objectsMap)
	if len(errors) > 0 {
		return errors
	}
	errors = b.writeFiles(codeFiles)
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (b *Backend) writeFiles(codeFiles []*types.DestinationFile) []error {
	var wg sync.WaitGroup
	errorsChan := make(chan error)

	for _, codeFile := range codeFiles {
		wg.Add(1)
		go func(cf *types.DestinationFile) {
			defer wg.Done()
			err := b.writeFile(*cf)
			if err != nil {
				errorsChan <- err
			}
		}(codeFile)
	}

	go func() {
		wg.Wait()
		close(errorsChan)
	}()
	errors := []error{}
	for err := range errorsChan {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (b *Backend) writeFile(codeFile types.DestinationFile) error {
	compilerConfig := config.GetCompilerConfig()
	fullPath := path.Join(compilerConfig.DestinationDir, codeFile.FilePath)

	err := os.MkdirAll(path.Dir(fullPath), os.ModePerm)
	if err != nil {
		return err
	}
	if _, err := os.Stat(fullPath); err == nil {
		err = os.Remove(fullPath)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(codeFile.Content)
	if err != nil {
		return err
	}
	return nil
}

func (b *Backend) objectsMap() map[string][]types.Object {
	result := make(map[string][]types.Object)
	for _, resource := range b.symbolTable.ResourceTable {
		result[resource.GetSourceFilePath()] = append(result[resource.GetSourceFilePath()], resource)
	}
	for _, variable := range b.symbolTable.VariableTable {
		result[variable.GetSourceFilePath()] = append(result[variable.GetSourceFilePath()], variable)
	}
	for filePath, objects := range result {
		sort.Slice(objects, func(i, j int) bool {
			return objects[i].GetSourceFileLine() < objects[j].GetSourceFileLine()
		})
		result[filePath] = objects
	}
	return result
}

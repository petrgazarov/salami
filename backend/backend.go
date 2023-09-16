package backend

import (
	"salami/backend/target"
	backendTypes "salami/backend/types"
	"salami/common/config"
	"salami/common/symbol_table"
	"salami/common/types"
	commonTypes "salami/common/types"
	"sort"
)

type Backend struct {
	symbolTable *symbol_table.SymbolTable
}

func NewBackend(symbolTable *symbol_table.SymbolTable, changeSet *types.ChangeSet) *Backend {
	return &Backend{
		symbolTable: symbolTable,
	}
}

func (b *Backend) GenerateCode() ([]string, []*backendTypes.BackendObject, []error) {
	targetModule, err := b.resolveTarget()
	if err != nil {
		return nil, nil, []error{err}
	}
	objectsMap := b.objectsMap()
	codeFiles, errors := targetModule.GenerateCodeFiles(objectsMap)
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (b *Backend) resolveTarget() (target.Target, error) {
	compilerConfig := config.GetCompilerConfig()
	return target.ResolveTarget(compilerConfig.Target, compilerConfig.Llm)
}

func (b *Backend) objectsMap() map[string][]commonTypes.ParsedObject {
	result := make(map[string][]commonTypes.ParsedObject)
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

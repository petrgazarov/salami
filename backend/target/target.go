package target

import (
	"salami/common/errors"
	"salami/common/symbol_table"
	"salami/common/types"
	"salami/compiler/backend/targets/terraform"
)

type Target interface {
	GenerateCodeFiles(objectsMap map[string][]types.ParsedObject) ([]*types.DestinationFile, []error)
}

func ResolveTarget(
	compilerConfig types.CompilerTargetConfig,
	llmConfig types.CompilerLlmConfig,
) (target Target, error error) {
	if compilerConfig.Platform == types.TerraformPlatform &&
		llmConfig.Provider == types.LlmOpenaiProvider &&
		llmConfig.Model == types.LlmGpt4Model {
		target = &terraform.TerraformTarget{}
	} else {
		return nil, &errors.ConfigError{Message: "unhandled configuration error"}
	}
	return target, nil
}

func GenerateCode(previousObjects []*types.Object, symbolTable *symbol_table.SymbolTable) []error {
	return nil
}

// func (b *Backend) objectsMap() map[string][]commonTypes.ParsedObject {
// 	result := make(map[string][]commonTypes.ParsedObject)
// 	for _, resource := range b.symbolTable.ResourceTable {
// 		result[resource.GetSourceFilePath()] = append(result[resource.GetSourceFilePath()], resource)
// 	}
// 	for _, variable := range b.symbolTable.VariableTable {
// 		result[variable.GetSourceFilePath()] = append(result[variable.GetSourceFilePath()], variable)
// 	}
// 	for filePath, objects := range result {
// 		sort.Slice(objects, func(i, j int) bool {
// 			return objects[i].GetSourceFileLine() < objects[j].GetSourceFileLine()
// 		})
// 		result[filePath] = objects
// 	}
// 	return result
// }

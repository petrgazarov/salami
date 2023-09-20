package target

import (
	"salami/backend/target/terraform"
	"salami/common/symbol_table"
	"salami/common/types"
)

type Target interface {
	GenerateCode(*types.ChangeSet, *symbol_table.SymbolTable) []error
	GetFilesFromObjects([]*types.Object) []*types.TargetFile
}

func ResolveTarget(
	targetConfig types.TargetConfig,
	llmConfig types.LlmConfig,
) Target {
	var target Target
	if targetConfig.Platform == types.TerraformPlatform &&
		llmConfig.Provider == types.LlmOpenaiProvider &&
		llmConfig.Model == types.LlmGpt4Model {
		target = &terraform.TerraformTarget{}
	}
	return target
}

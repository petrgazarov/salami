package target

import (
	"salami/backend/target/terraform"
	backendTypes "salami/backend/types"
	commonTypes "salami/common/types"
)

var targetFuncsMap = map[string]backendTypes.Target{
	"terraform": {
		GenerateCode:        terraform.GenerateCode,
		GetFilesFromObjects: terraform.GetFilesFromObjects,
	},
}

func ResolveTarget(
	targetConfig commonTypes.TargetConfig,
	llmConfig commonTypes.LlmConfig,
) backendTypes.Target {
	var target backendTypes.Target
	if targetConfig.Platform == commonTypes.TerraformPlatform &&
		llmConfig.Provider == commonTypes.LlmOpenaiProvider &&
		llmConfig.Model == commonTypes.LlmGpt4Model {
		target = targetFuncsMap["terraform"]
	}
	return target
}

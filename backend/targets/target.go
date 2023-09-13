package target

import (
	"salami/compiler/backend/target/terraform"
	"salami/compiler/errors"
	"salami/compiler/types"
)

type Target interface {
	GenerateCodeFiles(objectsMap map[string][]types.Object) ([]*types.DestinationFile, []error)
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

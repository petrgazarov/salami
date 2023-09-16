package config

import (
	"fmt"
	"log"
	"os"
	"salami/common/errors"
	"salami/common/types"

	"gopkg.in/yaml.v3"
)

var loadedConfig *types.Config

func GetConfig() *types.Config {
	if loadedConfig != nil {
		return loadedConfig
	} else {
		yamlFile, err := os.ReadFile("salami.yaml")
		if err != nil {
			log.Fatalf("failed to read config file. Ensure 'salami.yaml' exists in the root directory")
		}

		err = yaml.Unmarshal(yamlFile, &loadedConfig)
		if err != nil {
			log.Fatalf("error parsing config file: %v", err)
		}
		return loadedConfig
	}
}

func GetCompilerConfig() types.CompilerConfig {
	return GetConfig().Compiler
}

func ValidateConfig() error {
	config := GetConfig()
	if config.Compiler.Target.Platform != types.TerraformPlatform {
		return &errors.ConfigError{
			Message: fmt.Sprintf(
				"invalid target platform configuration. Supported values: '%s'",
				types.TerraformPlatform,
			),
		}
	}
	if _, err := os.Stat(config.Compiler.SourceDir); os.IsNotExist(err) {
		return &errors.ConfigError{
			Message: fmt.Sprintf("source directory '%s' does not exist", config.Compiler.SourceDir),
		}
	}
	if config.Compiler.Llm.Provider != types.LlmOpenaiProvider {
		return &errors.ConfigError{
			Message: fmt.Sprintf(
				"invalid LLM provider configuration. Supported values: '%s'",
				types.LlmOpenaiProvider,
			),
		}
	}
	if config.Compiler.Llm.Model != types.LlmGpt4Model {
		return &errors.ConfigError{
			Message: fmt.Sprintf("invalid LLM model configuration. Supported values: '%s'",
				types.LlmGpt4Model,
			),
		}
	}
	return nil
}

package config

import (
	"log"
	"os"
	"salami/common/types"

	"gopkg.in/yaml.v3"
)

var loadedConfig *configType

func GetSourceDir() string {
	return getConfig().compiler.sourceDir
}

func GetTargetDir() string {
	return getConfig().compiler.targetDir
}

func GetTargetConfig() types.TargetConfig {
	compilerTargetConfig := getConfig().compiler.target
	return types.TargetConfig{
		Platform: compilerTargetConfig.platform,
	}
}

func GetLlmConfig() types.LlmConfig {
	compilerLlmConfig := getConfig().compiler.llm
	return types.LlmConfig{
		Provider: compilerLlmConfig.provider,
		Model:    compilerLlmConfig.model,
	}
}

func getConfig() *configType {
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

package config

import (
	"log"
	"salami/common/types"
)

var configFilePath = "salami.yaml"
var loadedConfig *ConfigType

func SetConfigFilePath(path string) {
	configFilePath = path
	loadedConfig = nil
}

func GetSourceDir() string {
	return getConfig().Compiler.SourceDir
}

func GetTargetDir() string {
	return getConfig().Compiler.TargetDir
}

func GetTargetConfig() types.TargetConfig {
	compilerTargetConfig := getConfig().Compiler.Target
	return types.TargetConfig{
		Platform: compilerTargetConfig.Platform,
	}
}

func GetLlmConfig() types.LlmConfig {
	compilerLlmConfig := getConfig().Compiler.Llm
	return types.LlmConfig{
		Provider: compilerLlmConfig.Provider,
		Model:    compilerLlmConfig.Model,
	}
}

func getConfig() *ConfigType {
	if loadedConfig == nil {
		log.Fatal("Config not loaded")
	}
	return loadedConfig
}

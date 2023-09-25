package config

import (
	"log"
	"path/filepath"
	"salami/common/types"
)

var configFilePath = "salami.yaml"
var loadedConfig *ConfigType

func SetConfigFilePath(path string) {
	configFilePath = path
	loadedConfig = nil
}

func GetSourceDir() string {
	absPath, err := filepath.Abs(getConfig().Compiler.SourceDir)
	if err != nil {
		log.Fatal(err)
	}
	return absPath
}

func GetTargetDir() string {
	absPath, err := filepath.Abs(getConfig().Compiler.TargetDir)
	if err != nil {
		log.Fatal(err)
	}
	return absPath
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

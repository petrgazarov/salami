package config

import (
	"log"
	"os"
	"path/filepath"
	"salami/common/types"

	"gopkg.in/yaml.v3"
)

var configFilePath = "salami.yaml"
var loadedConfig *ConfigType

func SetConfigFilePath(path string) {
	configFilePath = path
	loadedConfig = nil
}

func LoadConfig() error {
	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(yamlFile, &loadedConfig); err != nil {
		return &ConfigError{Message: "could not parse config file. Ensure it is valid yaml format"}
	}
	return nil
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

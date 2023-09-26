package config

import (
	"log"
	"os"
	"regexp"
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
	yamlFile = []byte(substituteEnvVars(string(yamlFile)))
	if err = yaml.Unmarshal(yamlFile, &loadedConfig); err != nil {
		return &ConfigError{Message: "could not parse config file. Ensure it is valid yaml format"}
	}
	return nil
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
		ApiKey:   compilerLlmConfig.ApiKey,
	}
}

func getConfig() *ConfigType {
	if loadedConfig == nil {
		log.Fatal("Config not loaded")
	}
	return loadedConfig
}

func substituteEnvVars(input string) string {
	re := regexp.MustCompile(`\$\{(.+?)\}`)
	return re.ReplaceAllStringFunc(input, func(s string) string {
		varName := s[2 : len(s)-1]
		return os.Getenv(varName)
	})
}

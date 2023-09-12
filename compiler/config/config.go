package config

import (
	"fmt"
	"log"
	"os"
	"salami/compiler/errors"
	"salami/compiler/types"

	"gopkg.in/yaml.v3"
)

type CompilerConfig struct {
	Target         types.CompilerTargetConfig `yaml:"target"`
	Llm            types.CompilerLlmConfig    `yaml:"llm"`
	SourceDir      string                     `yaml:"source_dir"`
	DestinationDir string                     `yaml:"destination_dir"`
}

type Config struct {
	Compiler CompilerConfig `yaml:"compiler"`
}

var config Config

func init() {
	yamlFile, err := os.ReadFile("salami.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file. Ensure 'salami.yaml' exists in the root directory")
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}

func GetConfig() Config {
	return config
}

func GetCompilerConfig() CompilerConfig {
	return config.Compiler
}

func ValidateConfig() {
	if config.Compiler.Target.Platform != types.TerraformPlatform {
		panic(&errors.ConfigError{
			Message: fmt.Sprintf(
				"invalid target platform configuration. Supported values: '%s'",
				types.TerraformPlatform,
			),
		})
	}
	if _, err := os.Stat(config.Compiler.SourceDir); os.IsNotExist(err) {
		panic(&errors.ConfigError{
			Message: fmt.Sprintf("source directory '%s' does not exist", config.Compiler.SourceDir),
		})
	}
	if config.Compiler.Llm.Provider != types.LlmOpenaiProvider {
		panic(&errors.ConfigError{
			Message: fmt.Sprintf(
				"invalid LLM provider configuration. Supported values: '%s'",
				types.LlmOpenaiProvider,
			),
		})
	}
	if config.Compiler.Llm.Model != types.LlmGpt4Model {
		panic(&errors.ConfigError{
			Message: fmt.Sprintf("invalid LLM model configuration. Supported values: '%s'",
				types.LlmGpt4Model,
			),
		})
	}
}

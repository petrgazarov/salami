package config

import (
	"log"
	"os"
	"salami/compiler/types"

	"gopkg.in/yaml.v2"
)

type CompilerConfig struct {
	Target         types.CompilerTarget `yaml:"target"`
	Llm            types.Llm            `yaml:"llm"`
	SourceDir      string               `yaml:"source_dir"`
	DestinationDir string               `yaml:"destination_dir"`
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
	if config.Compiler.Target != types.Terraform {
		log.Fatalf("Invalid target configuration. Only '%s' is currently supported.", types.Terraform)
	}
	if _, err := os.Stat(config.Compiler.SourceDir); os.IsNotExist(err) {
		log.Fatalf("Source directory does not exist: %s", config.Compiler.SourceDir)
	}
	if config.Compiler.Llm != types.Gpt4 {
		log.Fatalf("Invalid LLM configuration. Only '%s' is currently supported.", types.Gpt4)
	}
}

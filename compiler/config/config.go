package config

import (
	"log"
	"os"
	"salami/compiler/types"

	"gopkg.in/yaml.v2"
)

type CompilerConfig struct {
	Target         types.CompilerTarget `yaml:"target"`
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
	if config.Compiler.Target != types.PulumiPython {
		log.Fatalf("Invalid target configuration. Only 'pulumi_python' is currently supported.")
	}
	if _, err := os.Stat(config.Compiler.SourceDir); os.IsNotExist(err) {
		log.Fatalf("Source directory does not exist: %s", config.Compiler.SourceDir)
	}
}

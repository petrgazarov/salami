package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Target string

const (
	PulumiPython Target = "pulumi_python"
)

type Config struct {
	Target         Target `yaml:"target"`
	SourceDir      string `yaml:"source_dir"`
	DestinationDir string `yaml:"destination_dir"`
}

var config Config

func init() {
	yamlFile, _ := os.ReadFile("salami.yaml")
	// TODO: handle error from os.ReadFile

	yaml.Unmarshal(yamlFile, &config)
	// TODO: handle error from yaml.Unmarshal
}

func GetConfig() Config {
	return config
}

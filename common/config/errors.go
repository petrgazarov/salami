package config

import "fmt"

var fieldNameMap = map[string]string{
	"ConfigType.Compiler.Target":          "target configuration",
	"ConfigType.Compiler.Target.Platform": "target platform",
	"ConfigType.Compiler.Llm":             "llm configuration",
	"ConfigType.Compiler.Llm.Provider":    "llm provider",
	"ConfigType.Compiler.Llm.Model":       "llm model",
	"ConfigType.Compiler.TargetDir":       "target directory",
}

type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error: %s", e.Message)
}

func getInvalidFieldError(namespace string, value interface{}) error {
	fieldName := fieldNameMap[namespace]
	if value == nil {
		return &ConfigError{Message: fmt.Sprintf("invalid %s", fieldName)}
	} else {
		return &ConfigError{Message: fmt.Sprintf("invalid %s: %v", fieldName, value)}
	}
}

func getMissingFieldError(namespace string) error {
	fieldName := fieldNameMap[namespace]
	if fieldName == "" {
		fieldName = namespace
	}
	return &ConfigError{Message: fmt.Sprintf("%s is required", fieldName)}
}

package errors

import "fmt"

type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error: %s", e.Message)
}

type SemanticError struct {
	SourceFilePath string
	Message        string
}

func (e *SemanticError) Error() string {
	return fmt.Sprintf(
		"\n%s\n  semantic error: %s",
		e.SourceFilePath,
		e.Message,
	)
}

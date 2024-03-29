package errors

import (
	"fmt"
	"salami/frontend/types"
)

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

type LexerError struct {
	FilePath string
	Line     int
	Column   int
	Message  string
}

func (e *LexerError) Error() string {
	return fmt.Sprintf(
		"\n%s\n  lexical error on line %d, column %d: %s",
		e.FilePath,
		e.Line,
		e.Column,
		e.Message,
	)
}

type ParseError struct {
	FilePath string
	Token    *types.Token
	Message  string
}

func (e *ParseError) Error() string {
	message := e.Message
	if message == "" {
		message = fmt.Sprintf(
			"unexpected token %s of type %s",
			e.Token.Value,
			e.Token.Type,
		)
	}

	return fmt.Sprintf(
		"\n%s\n  parsing error on line %d, column %d: %s",
		e.FilePath,
		e.Token.Line,
		e.Token.Column,
		message,
	)
}

type MissingEOFTokenError struct {
	FilePath string
}

func (e *MissingEOFTokenError) Error() string {
	return fmt.Sprintf("\n%s\n  parsing error: EOF token missing", e.FilePath)
}

type LlmError struct {
	Message string
}

func (e *LlmError) Error() string {
	return fmt.Sprintf("llm error: %s", e.Message)
}

type TargetError struct {
	Message string
}

func (e *TargetError) Error() string {
	return fmt.Sprintf("target error: %s", e.Message)
}
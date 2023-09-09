package errors

import (
	"fmt"
	"salami/compiler/types"
)

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

type MissingEOFToken struct {
	FilePath string
}

func (e *MissingEOFToken) Error() string {
	return "parsing error: EOF token missing"
}

package errors

import (
	"fmt"
	"salami/compiler/types"
)

type ParsingError struct {
	Token   types.Token
	Message string
}

func (e *ParsingError) Error() string {
	message := e.Message
	if message == "" {
		message = fmt.Sprintf(
			"unexpected token %s",
			e.Token.Value,
		)
	}

	return fmt.Sprintf(
		"parsing error at line %d, column %d: %s",
		e.Token.Line,
		e.Token.Column,
		message,
	)
}

type MissingEOFToken struct{}

func (e *MissingEOFToken) Error() string {
	return "parsing error: EOF token missing"
}

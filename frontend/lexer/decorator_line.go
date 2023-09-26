package lexer

import (
	"salami/common/errors"
	"salami/frontend/types"
	"unicode"
)

func (l *Lexer) processDecoratorLine() ([]*types.Token, error) {
	decoratorNameToken := l.getDecoratorNameToken()
	decoratorArgTokens := []*types.Token{}
	var err error
	if l.current() == '(' {
		l.advance()
		decoratorArgTokens, err = l.getDecoratorArgTokens()
		if err != nil {
			return nil, err
		}
	}
	l.skipWhitespace()
	if l.current() != '\n' {
		return nil, &errors.LexerError{
			FilePath: l.filePath,
			Line:     l.line,
			Column:   l.column,
			Message:  "decorator must be followed by arguments or a newline",
		}
	}

	return append(
		[]*types.Token{decoratorNameToken},
		decoratorArgTokens...,
	), nil
}

func (l *Lexer) getDecoratorNameToken() *types.Token {
	startPosition := l.pos
	startLine := l.line
	startColumn := l.column
	l.advance()
	for unicode.IsLetter(l.current()) {
		l.advance()
	}
	value := l.source[startPosition:l.pos]
	return l.newToken(types.DecoratorName, value, startLine, startColumn, false)
}

func (l *Lexer) getDecoratorArgTokens() ([]*types.Token, error) {
	var tokens []*types.Token
	l.skipWhitespace()

	for {
		startPosition := l.pos
		startLine := l.line
		startColumn := l.column
		for l.current() != ',' && l.current() != ')' && l.current() != '\n' && l.current() != 0 {
			l.advance()
		}
		if l.current() != ',' && l.current() != ')' {
			return nil, &errors.LexerError{
				FilePath: l.filePath,
				Line:     l.line,
				Column:   l.column,
				Message:  "decorator arguments must be followed by a comma or a closing parenthesis",
			}
		}
		value := l.source[startPosition:l.pos]
		tokens = append(tokens, l.newToken(types.DecoratorArg, value, startLine, startColumn, true))

		if l.current() == ')' {
			l.advance()
			break
		} else {
			l.advance()
			l.skipWhitespace()
		}
	}
	return tokens, nil
}

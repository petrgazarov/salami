package lexer

import (
	"salami/common/errors"
	"salami/frontend/types"
	"unicode"
)

func (l *Lexer) processConstructorLine() ([]*types.Token, error) {
	constructorNameToken := l.getConstructorNameToken()
	constructorArgTokens := []*types.Token{}
	var err error
	if l.current() == '(' {
		l.advance()
		constructorArgTokens, err = l.getConstructorArgTokens()
		if err != nil {
			return nil, err
		}
	}
	l.skipWhitespace()
	if l.current() != '\n' && l.current() != 0 {
		return nil, &errors.LexerError{
			FilePath: l.filePath,
			Line:     l.line,
			Column:   l.column,
			Message:  "constructor line must end in a newline or EOF",
		}
	}

	return append(
		[]*types.Token{constructorNameToken},
		constructorArgTokens...,
	), nil
}

func (l *Lexer) getConstructorNameToken() *types.Token {
	startPosition := l.pos
	startLine := l.line
	startColumn := l.column
	l.advance()
	for unicode.IsLetter(l.current()) {
		l.advance()
	}
	value := l.source[startPosition:l.pos]
	return l.newToken(types.ConstructorName, value, startLine, startColumn, false)
}

func (l *Lexer) getConstructorArgTokens() ([]*types.Token, error) {
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
				Message:  "constructor arguments must be followed by a comma or a closing parenthesis",
			}
		}
		value := l.source[startPosition:l.pos]
		tokens = append(tokens, l.newToken(types.ConstructorArg, value, startLine, startColumn, true))

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

func (l *Lexer) processLine() *types.Token {
	startPosition := l.pos
	startColumn := l.column
	for l.current() != '\n' && l.current() != 0 {
		l.advance()
	}
	lineText := l.source[startPosition:l.pos]
	return l.newToken(types.NaturalLanguage, lineText, l.line, startColumn, false)
}

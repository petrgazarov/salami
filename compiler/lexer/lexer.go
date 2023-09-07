package lexer

import (
	"salami/compiler/errors"
	"salami/compiler/types"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	filePath string
	source   string
	pos      int
	line     int
	column   int
	tokens   []*types.Token
}

func (l *Lexer) Run() ([]*types.Token, error) {
	for {
		switch {
		case l.current() == 0:
			l.tokens = append(l.tokens, l.newToken(types.EOF, "", l.line, l.column))
			return l.tokens, nil

		case l.current() == '\n':
			l.tokens = append(l.tokens, l.newToken(types.Newline, "", l.line, l.column))
			l.advance()

		case l.current() == '@':
			decoratorTokens, err := l.getDecoratorTokens()
			if err != nil {
				return nil, err
			}
			l.tokens = append(l.tokens, decoratorTokens...)

		default:
			startPosition := l.pos
			startLine := l.line
			startColumn := l.column
			for l.current() != ':' && l.current() != '\n' && l.current() != 0 {
				l.advance()
			}
			textSinceStartOfLine := l.source[startPosition:l.pos]

			if l.current() == ':' {
				_, validFieldName := types.ValidFieldNames[textSinceStartOfLine]
				if validFieldName {
					l.tokens = append(l.tokens, l.getFieldTokens(textSinceStartOfLine)...)
				} else {
					restOfLineText := l.getLineText()
					l.tokens = append(
						l.tokens,
						l.newToken(types.NaturalLanguage, textSinceStartOfLine+restOfLineText, startLine, startColumn),
					)
				}
			} else {
				l.tokens = append(
					l.tokens,
					&types.Token{Type: types.NaturalLanguage, Value: textSinceStartOfLine, Line: l.line, Column: l.column},
				)
			}
		}
	}
}

func (l *Lexer) getDecoratorTokens() ([]*types.Token, error) {
	startPosition := l.pos
	startLine := l.line
	startColumn := l.column
	l.advance()
	for unicode.IsLetter(l.current()) {
		l.advance()
	}
	value := l.source[startPosition:l.pos]
	decoratorNameToken := l.newToken(types.DecoratorName, value, startLine, startColumn)
	decoratorArgs := []*types.Token{}
	var err error
	if l.current() == '(' {
		l.advance()
		decoratorArgs, err = l.getDecoratorArgs()
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
		decoratorArgs...,
	), nil
}

func (l *Lexer) getDecoratorArgs() ([]*types.Token, error) {
	var tokens []*types.Token
	l.skipWhitespace()
	startPosition := l.pos
	startLine := l.line
	startColumn := l.column
	for l.current() != ')' {
		if l.current() == ',' {
			value := l.source[startPosition:l.pos]
			tokens = append(tokens, l.newToken(types.DecoratorArg, value, startLine, startColumn))
			l.advance()
			l.skipWhitespace()
			startPosition = l.pos
			startLine = l.line
			startColumn = l.column
		}
		l.advance()
	}
	value := l.source[startPosition:l.pos]
	tokens = append(tokens, l.newToken(types.DecoratorArg, value, startLine, startColumn))
	l.advance()
	l.skipWhitespace()
	if l.current() != '\n' {
		return nil, &errors.LexerError{
			FilePath: l.filePath,
			Line:     l.line,
			Column:   l.column,
			Message:  "decorator arguments must be followed by a newline",
		}
	}
	return tokens, nil
}

func (l *Lexer) getFieldTokens(fieldName string) []*types.Token {
	fieldNameToken := l.newToken(types.FieldName, fieldName)
	l.advance()
	fieldValue := l.getLineText()
	return append(
		[]*types.Token{fieldNameToken},
		l.newToken(types.FieldValue, fieldValue),
	)
}

func (l *Lexer) getLineText() string {
	start := l.pos
	for l.current() != '\n' && l.current() != 0 {
		l.advance()
	}
	return l.source[start:l.pos]
}

func (l *Lexer) newToken(tokenType types.TokenType, value string, line int, column int) *types.Token {
	return &types.Token{Type: tokenType, Value: value, Line: line, Column: column}
}

func (l *Lexer) current() rune {
	if l.pos >= len(l.source) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.source[l.pos:])
	return r
}

func (l *Lexer) advance() {
	if l.pos >= len(l.source) {
		return
	}
	_, width := utf8.DecodeRuneInString(l.source[l.pos:])
	l.pos += width
	if l.current() == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
}

func (l *Lexer) skipWhitespace() {
	for l.current() == ' ' {
		l.advance()
	}
}

func NewLexer(filePath string, source string) *Lexer {
	return &Lexer{
		filePath: filePath,
		source:   source,
		pos:      0,
		line:     1,
		column:   1,
	}
}

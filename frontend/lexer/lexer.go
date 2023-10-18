package lexer

import (
	"salami/frontend/types"
	"strings"
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
			l.tokens = append(l.tokens, l.newToken(types.EOF, "", l.line, l.column, false))
			return l.tokens, nil

		case l.current() == '\n':
			l.tokens = append(l.tokens, l.newToken(types.Newline, "", l.line, l.column, false))
			l.advance()

		case l.current() == '@':
			constructorTokens, err := l.processConstructorLine()
			if err != nil {
				return nil, err
			}
			l.tokens = append(l.tokens, constructorTokens...)

		default:
			lineToken := l.processLine()
			l.tokens = append(l.tokens, lineToken)
		}
	}
}

func (l *Lexer) newToken(
	tokenType types.TokenType,
	value string,
	line int,
	column int,
	trimWhitespace bool,
) *types.Token {
	normalizedValue := value
	if trimWhitespace {
		normalizedValue = strings.TrimSpace(value)
		leadingSpaces := len(value) - len(strings.TrimLeft(value, " "))
		column += leadingSpaces
	}
	return &types.Token{Type: tokenType, Value: normalizedValue, Line: line, Column: column}
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
	if l.current() == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	_, width := utf8.DecodeRuneInString(l.source[l.pos:])
	l.pos += width
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

package lexer

import (
	"salami/compiler/types"
	"strings"
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

func (l *Lexer) Run() []*types.Token {
	for {
		switch {
		case l.current() == 0:
			l.tokens = append(l.tokens, &types.Token{Type: types.EOF, Line: l.line, Column: l.column})
			return l.tokens

		case l.current() == '\n':
			l.tokens = append(l.tokens, &types.Token{Type: types.Newline, Line: l.line, Column: l.column})
			l.advance()

		case l.current() == '@':
			l.tokens = append(l.tokens, l.getDecoratorTokens()...)

		default:
			start := l.pos
			for l.current() != ':' && l.current() != '\n' && l.current() != 0 {
				l.advance()
			}
			subText := l.source[start:l.pos]

			if l.current() == ':' {
				_, validFieldName := types.ValidFieldNames[subText]
				if validFieldName {
					l.tokens = append(l.tokens, l.getFieldTokens(subText)...)
				} else {
					restOfLineText := l.getLineText()
					l.tokens = append(
						l.tokens,
						&types.Token{Type: types.NaturalLanguage, Value: subText + restOfLineText, Line: l.line, Column: l.column},
					)
				}
			} else {
				l.tokens = append(
					l.tokens,
					&types.Token{Type: types.NaturalLanguage, Value: subText, Line: l.line, Column: l.column},
				)
			}
		}
	}
}

func (l *Lexer) getDecoratorTokens() []*types.Token {
	startLine, startColumn := l.line, l.column
	l.advance()
	start := l.pos
	for unicode.IsLetter(l.current()) {
		l.advance()
	}
	value := l.source[start:l.pos]
	decoratorArgs := []*types.Token{}
	if l.current() == '(' {
		l.advance()
		decoratorArgs = l.getDecoratorArgs()
	}
	return append(
		[]*types.Token{{Type: types.DecoratorName, Value: value, Line: startLine, Column: startColumn}},
		decoratorArgs...,
	)
}

func (l *Lexer) getDecoratorArgs() []*types.Token {
	var tokens []*types.Token
	start := l.pos
	for l.current() != ')' {
		if l.current() == ',' {
			value := strings.TrimSpace(l.source[start:l.pos])
			tokens = append(tokens, &types.Token{Type: types.DecoratorArg, Value: value, Line: l.line, Column: l.column})
			l.advance()
			start = l.pos
		}
		l.advance()
	}
	value := strings.TrimSpace(l.source[start:l.pos])
	tokens = append(tokens, &types.Token{Type: types.DecoratorArg, Value: value, Line: l.line, Column: l.column})
	for l.current() != '\n' && l.current() != 0 {
		l.advance()
	}
	return tokens
}

func (l *Lexer) getFieldTokens(fieldName string) []*types.Token {
	fieldNameLineStart, fieldNameColumnStart := l.line, l.column
	l.advance()
	fieldValue := l.getLineText()
	return append(
		[]*types.Token{{Type: types.FieldName, Value: fieldName, Line: fieldNameLineStart, Column: fieldNameColumnStart}},
		&types.Token{Type: types.FieldValue, Value: fieldValue, Line: l.line, Column: l.column},
	)
}

func (l *Lexer) getLineText() string {
	start := l.pos
	for l.current() != '\n' && l.current() != 0 {
		l.advance()
	}
	return l.source[start:l.pos]
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

func NewLexer(filePath string, source string) *Lexer {
	return &Lexer{
		filePath: filePath,
		source:   source,
		pos:      0,
		line:     1,
		column:   1,
	}
}

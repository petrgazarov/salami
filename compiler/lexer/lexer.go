package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	DecoratorName TokenType = iota
	DecoratorArg
	FieldName
	FieldValue
	FreeText
	VariableRef
	Newline
	EOF
	Error
)

var TokenTypeNames = map[TokenType]string{
	DecoratorName: "DecoratorName",
	DecoratorArg:  "DecoratorArg",
	FieldName:     "FieldName",
	FieldValue:    "FieldValue",
	FreeText:      "FreeText",
	VariableRef:   "VariableRef",
	Newline:       "Newline",
	EOF:           "EOF",
	Error:         "Error",
}

func (t TokenType) String() string {
	return TokenTypeNames[t]
}

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

type Lexer struct {
	filePath    string
	source      string
	pos         int
	line        int
	column      int
	tokensQueue []Token
}

var validFieldNames = map[string]bool{
	"Resource type": true,
	"Logical name":  true,
	"Description":   true,
	"Name":          true,
	"Value":         true,
}

func (l *Lexer) NextToken() Token {
	if len(l.tokensQueue) > 0 {
		token := l.tokensQueue[0]
		l.tokensQueue = l.tokensQueue[1:]
		return token
	}

	switch {
	case l.current() == 0:
		return Token{Type: EOF, Line: l.line, Column: l.column}

	case l.current() == '\n':
		l.advance()
		return Token{Type: Newline, Line: l.line, Column: l.column}

	case l.current() == '@':
		return l.getDecoratorTokens()

	default:
		start := l.pos
		for l.current() != ':' && l.current() != '\n' && l.current() != 0 {
			l.advance()
		}
		subText := strings.TrimSpace(l.source[start:l.pos])
		if l.current() == ':' {
			if _, ok := validFieldNames[subText]; ok {
				return l.getFieldTokens(subText)
			}
		}

		restOfLineText := l.getLineText()
		return Token{Type: FreeText, Value: subText + restOfLineText, Line: l.line, Column: l.column}
	}
}

func (l *Lexer) getDecoratorTokens() Token {
	l.advance()
	start := l.pos
	for unicode.IsLetter(l.current()) {
		l.advance()
	}
	value := l.source[start:l.pos]
	decoratorArgs := []Token{}
	if l.current() == '(' {
		l.advance()
		decoratorArgs = l.getDecoratorArgs()
	}
	l.tokensQueue = append(
		l.tokensQueue,
		Token{Type: DecoratorName, Value: value, Line: l.line, Column: l.column},
	)
	l.tokensQueue = append(l.tokensQueue, decoratorArgs...)
	return l.NextToken()
}

func (l *Lexer) getFieldTokens(fieldName string) Token {
	l.advance()
	fieldValue := l.getLineText()
	l.tokensQueue = append(
		l.tokensQueue,
		Token{Type: FieldName, Value: fieldName, Line: l.line, Column: l.column},
		Token{Type: FieldValue, Value: fieldValue, Line: l.line, Column: l.column},
	)
	return l.NextToken()
}

func (l *Lexer) getDecoratorArgs() []Token {
	var tokens []Token
	start := l.pos
	for l.current() != ')' {
		if l.current() == ',' {
			value := strings.TrimSpace(l.source[start:l.pos])
			tokens = append(tokens, Token{Type: DecoratorArg, Value: value, Line: l.line, Column: l.column})
			l.advance()
			start = l.pos
		}
		l.advance()
	}
	value := strings.TrimSpace(l.source[start:l.pos])
	tokens = append(tokens, Token{Type: DecoratorArg, Value: value, Line: l.line, Column: l.column})
	l.advance()
	return tokens
}

func (l *Lexer) getLineText() string {
	start := l.pos
	for l.current() != '\n' && l.current() != 0 {
		l.advance()
	}
	lineText := strings.TrimSpace(l.source[start:l.pos])

	if l.current() == '\n' {
		l.tokensQueue = append(l.tokensQueue, Token{Type: Newline, Line: l.line, Column: l.column})
		l.advance()
	} else if l.current() == 0 {
		l.tokensQueue = append(l.tokensQueue, Token{Type: EOF, Line: l.line, Column: l.column})
	}

	return lineText
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
		l.column = 0
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

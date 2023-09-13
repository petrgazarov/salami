package lexer

import (
	"salami/frontend/types"
)

func (l *Lexer) processLine() ([]*types.Token, error) {
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
			fieldNameToken := l.newToken(types.FieldName, textSinceStartOfLine, startLine, startColumn, false)
			l.advance()
			startLine := l.line
			startColumn := l.column
			fieldValueToken := l.newToken(types.FieldValue, l.getLineText(), startLine, startColumn, true)
			return []*types.Token{fieldNameToken, fieldValueToken}, nil
		} else {
			restOfLineText := l.getLineText()
			newToken := l.newToken(types.NaturalLanguage, textSinceStartOfLine+restOfLineText, startLine, startColumn, false)
			return []*types.Token{newToken}, nil
		}
	} else {
		newToken := l.newToken(types.NaturalLanguage, textSinceStartOfLine, startLine, startColumn, false)
		return []*types.Token{newToken}, nil
	}
}

func (l *Lexer) getLineText() string {
	start := l.pos
	for l.current() != '\n' && l.current() != 0 {
		l.advance()
	}
	return l.source[start:l.pos]
}

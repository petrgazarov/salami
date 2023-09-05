package parser

import (
	"salami/compiler/errors"
	"salami/compiler/types"
)

func (p *Parser) handleNaturalLanguageLine() error {
	freeTextToken := p.currentToken()
	p.advance()
	for p.currentToken().Type != types.Newline && p.currentToken().Type != types.EOF {
		return &errors.ParsingError{Token: p.currentToken()}
	}
	p.advance()

	if p.currentObjectTypeIs(Unset) {
		return &errors.ParsingError{
			Token:   freeTextToken,
			Message: "ambiguous object type. Ensure object has Resource type field or @variable decorator before the current line",
		}
	}
	if !p.currentObjectTypeIs(Resource) {
		return &errors.ParsingError{
			Token:   freeTextToken,
			Message: "natural language can only be used on resource",
		}
	}
	p.addLineToNaturalLanguage(freeTextToken.Value)
	return nil
}

func (p *Parser) addLineToNaturalLanguage(line string) {
	if p.currentResource().NaturalLanguage != "" {
		p.currentResource().NaturalLanguage += "\n"
	}
	p.currentResource().NaturalLanguage += (line + "\n")
}

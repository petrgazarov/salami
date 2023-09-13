package parser

import (
	"regexp"
	"salami/frontend/types"
	"strings"
)

func (p *Parser) handleNaturalLanguageLine() error {
	freeTextToken := p.currentToken()
	p.advance()
	for p.currentToken().Type != types.Newline && p.currentToken().Type != types.EOF {
		return p.parseError(p.currentToken())
	}

	if p.currentObjectTypeIs(Unset) {
		return p.parseError(
			freeTextToken,
			"ambiguous object type. Ensure object has Resource type field or @variable decorator before the current line",
		)
	}
	if !p.currentObjectTypeIs(Resource) {
		return p.parseError(freeTextToken, "natural language can only be used on resource")
	}
	err := p.addLineToNaturalLanguage(freeTextToken.Value)
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) addLineToNaturalLanguage(line string) error {
	addition := line

	if p.currentResource().NaturalLanguage != "" {
		addition = "\n" + addition
	}
	p.currentResource().NaturalLanguage += addition
	err := p.extractAndSetReferencedVariables(line)
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) extractAndSetReferencedVariables(text string) error {
	re := regexp.MustCompile(`\{([^{}]+)\}`)
	isAlphanumeric := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		variable := strings.TrimSpace(match[1])
		if strings.Contains(variable, "{") || strings.Contains(variable, "}") {
			return p.parseError(
				p.currentToken(),
				"Nested curly braces in referenced variables are not allowed",
			)
		}
		if !isAlphanumeric.MatchString(variable) {
			return p.parseError(
				p.currentToken(),
				"Variable inside curly braces must be alphanumeric",
			)
		}
		isUnique := true
		for _, v := range p.currentResource().ReferencedVariables {
			if v == variable {
				isUnique = false
				break
			}
		}
		if isUnique {
			p.currentResource().ReferencedVariables = append(p.currentResource().ReferencedVariables, variable)
		}
	}
	return nil
}

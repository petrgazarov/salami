package parser

import (
	"regexp"
	commonTypes "salami/common/types"
	"salami/frontend/types"
	"strings"
)

func (p *Parser) parseNaturalLanguage() error {
	freeTextToken := p.currentToken()
	p.advance()
	for p.currentToken().Type != types.Newline && p.currentToken().Type != types.EOF {
		return p.parseError(p.currentToken())
	}

	if p.currentObjectTypeIs(Unset) {
		return p.parseError(
			freeTextToken,
			"ambiguous object type. Object must start with a constructor",
		)
	}
	p.currentObject().AddNaturalLanguage(freeTextToken.Value)
	err := p.extractAndSetReferencedVariables(freeTextToken.Value)
	if err != nil {
		return err
	}
	err = p.extractAndSetReferencedResources(freeTextToken.Value)
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
				"nested curly braces in referenced variables are not allowed",
			)
		}
		if !isAlphanumeric.MatchString(variable) {
			return p.parseError(
				p.currentToken(),
				"variable inside curly braces must be alphanumeric. Underscores are allowed.",
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

func (p *Parser) extractAndSetReferencedResources(text string) error {
	re := regexp.MustCompile(`\$([a-zA-Z0-9_]+)`)
	matches := re.FindAllStringSubmatchIndex(text, -1)

	for _, match := range matches {
		// Ignore if the dollar sign is escaped
		if match[0] > 0 && text[match[0]-1] == '\\' {
			continue
		}

		resource := strings.TrimSpace(text[match[2]:match[3]])
		if strings.Contains(resource, "$") {
			return p.parseError(
				p.currentToken(),
				"Nested dollar signs in referenced resources are not allowed",
			)
		}
		isUnique := true
		for _, r := range p.currentResource().ReferencedResources {
			if r == commonTypes.LogicalName(resource) {
				isUnique = false
				break
			}
		}
		if isUnique {
			p.currentResource().ReferencedResources = append(
				p.currentResource().ReferencedResources,
				commonTypes.LogicalName(resource),
			)
		}
	}
	return nil
}

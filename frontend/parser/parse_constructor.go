package parser

import (
	commonTypes "salami/common/types"
	frontendTypes "salami/frontend/types"
)

func (p *Parser) parseConstructor() error {
	constructorNameToken := p.currentToken()
	var constructorArgTokens []*frontendTypes.Token
	p.advance()
	for p.currentToken().Type != frontendTypes.Newline && p.currentToken().Type != frontendTypes.EOF {
		if p.currentToken().Type == frontendTypes.ConstructorArg {
			constructorArgTokens = append(constructorArgTokens, p.currentToken())
		} else {
			return p.parseError(p.currentToken())
		}
		p.advance()
	}

	switch constructorNameToken.Value {
	case "@resource":
		err := p.parseResourceConstructor(constructorNameToken, constructorArgTokens)
		if err != nil {
			return err
		}
	case "@variable":
		err := p.parseVariableConstructor(constructorNameToken, constructorArgTokens)
		if err != nil {
			return err
		}
	default:
		return p.parseError(constructorNameToken)
	}
	if p.currentToken().Type != frontendTypes.Newline && p.currentToken().Type != frontendTypes.EOF {
		return p.parseError(p.currentToken())
	}
	return nil
}

func (p *Parser) parseResourceConstructor(
	constructorNameToken *frontendTypes.Token,
	constructorArgTokens []*frontendTypes.Token,
) error {
	if !p.currentObjectTypeIs(Unset) {
		return p.parseError(
			constructorNameToken,
			"@resource constructor must be the first line of a resource block",
		)
	}
	p.setCurrentObjectType(Resource)
	if len(constructorArgTokens) != 2 {
		return p.parseError(constructorNameToken, "exactly two arguments are expected for @resource constructor")
	}
	p.currentResource().ResourceType = commonTypes.ResourceType(constructorArgTokens[0].Value)
	p.currentResource().LogicalName = commonTypes.LogicalName(constructorArgTokens[1].Value)
	return nil
}

func (p *Parser) parseVariableConstructor(
	constructorNameToken *frontendTypes.Token,
	constructorArgTokens []*frontendTypes.Token,
) error {
	if !p.currentObjectTypeIs(Unset) {
		return p.parseError(
			constructorNameToken,
			"@variable constructor must be the first line of a variable block",
		)
	}
	p.setCurrentObjectType(Variable)

	if len(constructorArgTokens) > 3 || len(constructorArgTokens) < 2 {
		return p.parseError(constructorNameToken, "exactly two or three arguments are expected for @variable constructor")
	}
	p.currentVariable().Name = constructorArgTokens[0].Value
	p.currentVariable().Type = commonTypes.VariableType(constructorArgTokens[1].Value)
	if len(constructorArgTokens) == 3 {
		p.currentVariable().Default = constructorArgTokens[2].Value
	}
	return nil
}

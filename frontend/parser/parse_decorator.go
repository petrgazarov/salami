package parser

import (
	commonTypes "salami/common/types"
	frontendTypes "salami/frontend/types"
)

func (p *Parser) parseDecorator() error {
	decoratorNameToken := p.currentToken()
	var decoratorArgTokens []*frontendTypes.Token
	p.advance()
	for p.currentToken().Type != frontendTypes.Newline && p.currentToken().Type != frontendTypes.EOF {
		if p.currentToken().Type == frontendTypes.DecoratorArg {
			decoratorArgTokens = append(decoratorArgTokens, p.currentToken())
		} else {
			return p.parseError(p.currentToken())
		}
		p.advance()
	}

	switch decoratorNameToken.Value {
	case "@resource":
		err := p.parseResourceDecorator(decoratorNameToken, decoratorArgTokens)
		if err != nil {
			return err
		}
	case "@variable":
		err := p.parseVariableDecorator(decoratorNameToken, decoratorArgTokens)
		if err != nil {
			return err
		}
	default:
		return p.parseError(decoratorNameToken)
	}
	if p.currentToken().Type != frontendTypes.Newline && p.currentToken().Type != frontendTypes.EOF {
		return p.parseError(p.currentToken())
	}
	return nil
}

func (p *Parser) parseResourceDecorator(
	decoratorNameToken *frontendTypes.Token,
	decoratorArgTokens []*frontendTypes.Token,
) error {
	if !p.currentObjectTypeIs(Unset) {
		return p.parseError(
			decoratorNameToken,
			"@resource decorator must be the first line of a resource block",
		)
	}
	p.setCurrentObjectType(Resource)
	if len(decoratorArgTokens) != 2 {
		return p.parseError(decoratorNameToken, "exactly two arguments are expected for @resource decorator")
	}
	p.currentResource().ResourceType = commonTypes.ResourceType(decoratorArgTokens[0].Value)
	p.currentResource().LogicalName = commonTypes.LogicalName(decoratorArgTokens[1].Value)
	return nil
}

func (p *Parser) parseVariableDecorator(
	decoratorNameToken *frontendTypes.Token,
	decoratorArgTokens []*frontendTypes.Token,
) error {
	if !p.currentObjectTypeIs(Unset) {
		return p.parseError(
			decoratorNameToken,
			"@variable decorator must be the first line of a variable block",
		)
	}
	p.setCurrentObjectType(Variable)

	if len(decoratorArgTokens) > 3 || len(decoratorArgTokens) < 2 {
		return p.parseError(decoratorNameToken, "exactly two or three arguments are expected for @variable decorator")
	}
	p.currentVariable().Name = decoratorArgTokens[0].Value
	p.currentVariable().Type = commonTypes.VariableType(decoratorArgTokens[1].Value)
	if len(decoratorArgTokens) == 3 {
		p.currentVariable().Default = decoratorArgTokens[2].Value
	}
	return nil
}

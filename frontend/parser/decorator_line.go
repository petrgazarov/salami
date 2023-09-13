package parser

import (
	commonTypes "salami/common/types"
	frontendTypes "salami/frontend/types"
	"strings"
)

func (p *Parser) handleDecoratorLine() error {
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
	case "@uses":
		err := p.handleUsesDecorator(decoratorNameToken, decoratorArgTokens)
		if err != nil {
			return err
		}
	case "@exports":
		err := p.handleExportsDecorator(decoratorNameToken, decoratorArgTokens)
		if err != nil {
			return err
		}
	case "@variable":
		err := p.handleVariableDecorator(decoratorNameToken, decoratorArgTokens)
		if err != nil {
			return err
		}
	default:
		return p.parseError(decoratorNameToken)
	}
	if p.currentToken().Type != frontendTypes.Newline && p.currentToken().Type != frontendTypes.EOF {
		return p.parseError(p.currentToken())
	}
	p.advance()
	return nil
}

func (p *Parser) handleUsesDecorator(
	decoratorNameToken *frontendTypes.Token,
	decoratorArgTokens []*frontendTypes.Token,
) error {
	if p.currentObjectTypeIs(Unset) {
		p.setCurrentObjectType(Resource)
	}
	if !p.currentObjectTypeIs(Resource) {
		return p.parseError(decoratorNameToken, "@uses decorator can only be used on resource")
	}
	for _, arg := range decoratorArgTokens {
		p.currentResource().Uses = append(p.currentResource().Uses, commonTypes.LogicalName(arg.Value))
	}
	return nil
}

func (p *Parser) handleExportsDecorator(
	decoratorNameToken *frontendTypes.Token,
	decoratorArgTokens []*frontendTypes.Token,
) error {
	if p.currentObjectTypeIs(Unset) {
		p.setCurrentObjectType(Resource)
	}
	if !p.currentObjectTypeIs(Resource) {
		return p.parseError(decoratorNameToken, "@exports decorator can only be used on resource")
	}
	for _, arg := range decoratorArgTokens {
		kv := strings.Split(arg.Value, ":")
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		p.currentResource().Exports[key] = value
	}
	return nil
}

func (p *Parser) handleVariableDecorator(
	decoratorNameToken *frontendTypes.Token,
	decoratorArgTokens []*frontendTypes.Token,
) error {
	if p.currentObjectTypeIs(Unset) {
		p.setCurrentObjectType(Variable)
	}
	if !p.currentObjectTypeIs(Variable) {
		return p.parseError(decoratorNameToken, "@variable decorator can only be used on variable")
	}
	if len(decoratorArgTokens) > 1 {
		return p.parseError(decoratorNameToken, "Only one argument is expected for @variable decorator")
	}
	variableType, err := commonTypes.StringToVariableType(decoratorArgTokens[0].Value)
	if err != nil {
		return p.parseError(decoratorNameToken, err.Error())
	}
	p.currentVariable().Type = variableType
	return nil
}

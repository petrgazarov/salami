package parser

import (
	"salami/compiler/errors"
	"salami/compiler/types"
	"strings"
)

func (p *Parser) handleDecoratorLine() error {
	decoratorNameToken := p.currentToken()
	var decoratorArgTokens []types.Token
	p.advance()
	for p.currentToken().Type != types.Newline && p.currentToken().Type != types.EOF {
		if p.currentToken().Type == types.DecoratorArg {
			decoratorArgTokens = append(decoratorArgTokens, p.currentToken())
		} else {
			return &errors.ParsingError{Token: p.currentToken()}
		}
		p.advance()
	}
	p.advance()

	switch decoratorNameToken.Value {
	case "uses":
		err := p.handleUsesDecorator(decoratorNameToken, decoratorArgTokens)
		if err != nil {
			return err
		}
	case "exports":
		err := p.handleExportsDecorator(decoratorNameToken, decoratorArgTokens)
		if err != nil {
			return err
		}
	case "variable":
		err := p.handleVariableDecorator(decoratorNameToken, decoratorArgTokens)
		if err != nil {
			return err
		}
	}
	return &errors.ParsingError{Token: decoratorNameToken}
}

func (p *Parser) handleUsesDecorator(decoratorNameToken types.Token, decoratorArgTokens []types.Token) error {
	if p.currentObjectTypeIs(Unset) {
		p.setCurrentObjectType(Resource)
	}
	if !p.currentObjectTypeIs(Resource) {
		return &errors.ParsingError{
			Token:   decoratorNameToken,
			Message: "@uses decorator can only be used on resource",
		}
	}
	for _, arg := range decoratorArgTokens {
		p.currentResource().Uses = append(p.currentResource().Uses, types.LogicalName(arg.Value))
	}
	return nil
}

func (p *Parser) handleExportsDecorator(decoratorNameToken types.Token, decoratorArgTokens []types.Token) error {
	if p.currentObjectTypeIs(Unset) {
		p.setCurrentObjectType(Resource)
	}
	if !p.currentObjectTypeIs(Resource) {
		return &errors.ParsingError{
			Token:   decoratorNameToken,
			Message: "@exports decorator can only be used on resource",
		}
	}
	for _, arg := range decoratorArgTokens {
		kv := strings.Split(arg.Value, ":")
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		p.currentResource().Exports[key] = value
	}
	return nil
}

func (p *Parser) handleVariableDecorator(decoratorNameToken types.Token, decoratorArgTokens []types.Token) error {
	if p.currentObjectTypeIs(Unset) {
		p.setCurrentObjectType(Variable)
	}
	if !p.currentObjectTypeIs(Variable) {
		return &errors.ParsingError{
			Token:   decoratorNameToken,
			Message: "@variable decorator can only be used on variable",
		}
	}
	if len(decoratorArgTokens) > 1 {
		return &errors.ParsingError{
			Token:   decoratorNameToken,
			Message: "Only one argument is expected for @variable decorator",
		}
	}
	variableType, err := types.StringToVariableType(decoratorArgTokens[0].Value)
	if err != nil {
		return &errors.ParsingError{
			Token:   decoratorNameToken,
			Message: err.Error(),
		}
	}
	p.currentVariable().Type = variableType
	return nil
}

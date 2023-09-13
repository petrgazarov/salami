package parser

import (
	commonTypes "salami/common/types"
	"salami/frontend/errors"
	frontendTypes "salami/frontend/types"
)

type ObjectType int

const (
	Unset ObjectType = iota
	Resource
	Variable
)

type Parser struct {
	tokens            []*frontendTypes.Token
	resources         []*commonTypes.Resource
	variables         []*commonTypes.Variable
	index             int
	currentObjectType ObjectType
	filePath          string
}

func NewParser(tokens []*frontendTypes.Token, filePath string) *Parser {
	return &Parser{
		tokens:            tokens,
		resources:         make([]*commonTypes.Resource, 0),
		variables:         make([]*commonTypes.Variable, 0),
		index:             0,
		currentObjectType: Unset,
		filePath:          filePath,
	}
}

func (p *Parser) Parse() ([]*commonTypes.Resource, []*commonTypes.Variable, error) {
	for p.index < len(p.tokens) {
		switch p.currentToken().Type {
		case frontendTypes.EOF:
			return p.resources, p.variables, nil
		case frontendTypes.Newline:
			p.advance()
			if p.currentToken().Type == frontendTypes.Newline {
				p.setCurrentObjectType(Unset)
			}
		case frontendTypes.DecoratorName:
			err := p.handleDecoratorLine()
			if err != nil {
				return nil, nil, err
			}
		case frontendTypes.FieldName:
			err := p.handleFieldLine()
			if err != nil {
				return nil, nil, err
			}
		case frontendTypes.NaturalLanguage:
			err := p.handleNaturalLanguageLine()
			if err != nil {
				return nil, nil, err
			}
		default:
			return nil, nil, p.parseError(p.currentToken())
		}
	}
	return nil, nil, &errors.MissingEOFToken{FilePath: p.filePath}
}

func (p *Parser) currentResource() *commonTypes.Resource {
	return p.resources[len(p.resources)-1]
}

func (p *Parser) currentVariable() *commonTypes.Variable {
	return p.variables[len(p.variables)-1]
}

func (p *Parser) setCurrentObjectType(t ObjectType) {
	switch t {
	case Resource:
		newResource := commonTypes.NewResource(p.filePath, p.currentToken().Line)
		p.resources = append(p.resources, newResource)
	case Variable:
		newVariable := commonTypes.NewVariable(p.filePath, p.currentToken().Line)
		p.variables = append(p.variables, newVariable)
	}
	p.currentObjectType = t
}

func (p *Parser) currentObjectTypeIs(objectType ObjectType) bool {
	return p.currentObjectType == objectType
}

func (p *Parser) currentToken() *frontendTypes.Token {
	return p.tokens[p.index]
}

func (p *Parser) advance() {
	p.index++
}

func (p *Parser) parseError(token *frontendTypes.Token, messages ...string) error {
	error := errors.ParseError{
		Token:    token,
		FilePath: p.filePath,
	}
	if len(messages) > 0 {
		error.Message = messages[0]
	}
	return &error
}

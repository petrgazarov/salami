package parser

import (
	"salami/compiler/errors"
	"salami/compiler/types"
)

type ObjectType int

const (
	Unset ObjectType = iota
	Resource
	Variable
)

type Parser struct {
	tokens            []*types.Token
	resources         []*types.Resource
	variables         []*types.Variable
	index             int
	currentObjectType ObjectType
	filePath          string
}

func NewParser(tokens []*types.Token, filePath string) *Parser {
	return &Parser{
		tokens:            tokens,
		resources:         make([]*types.Resource, 0),
		variables:         make([]*types.Variable, 0),
		index:             0,
		currentObjectType: Unset,
		filePath:          filePath,
	}
}

func (p *Parser) Parse() ([]*types.Resource, []*types.Variable, error) {
	for p.index < len(p.tokens) {
		switch p.currentToken().Type {
		case types.EOF:
			return p.resources, p.variables, nil
		case types.Newline:
			p.advance()
			if p.currentToken().Type == types.Newline {
				p.setCurrentObjectType(Unset)
			}
		case types.DecoratorName:
			err := p.handleDecoratorLine()
			if err != nil {
				return nil, nil, err
			}
		case types.FieldName:
			err := p.handleFieldLine()
			if err != nil {
				return nil, nil, err
			}
		case types.NaturalLanguage:
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

func (p *Parser) currentResource() *types.Resource {
	return p.resources[len(p.resources)-1]
}

func (p *Parser) currentVariable() *types.Variable {
	return p.variables[len(p.variables)-1]
}

func (p *Parser) setCurrentObjectType(t ObjectType) {
	switch t {
	case Resource:
		newResource := types.NewResource(p.filePath, p.currentToken().Line)
		p.resources = append(p.resources, newResource)
	case Variable:
		newVariable := types.NewVariable(p.filePath, p.currentToken().Line)
		p.variables = append(p.variables, newVariable)
	}
	p.currentObjectType = t
}

func (p *Parser) currentObjectTypeIs(objectType ObjectType) bool {
	return p.currentObjectType == objectType
}

func (p *Parser) currentToken() *types.Token {
	return p.tokens[p.index]
}

func (p *Parser) advance() {
	p.index++
}

func (p *Parser) parseError(token *types.Token, messages ...string) error {
	error := errors.ParseError{
		Token:    token,
		FilePath: p.filePath,
	}
	if len(messages) > 0 {
		error.Message = messages[0]
	}
	return &error
}

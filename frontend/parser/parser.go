package parser

import (
	"salami/common/errors"
	commonTypes "salami/common/types"
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
	resources         []*commonTypes.ParsedResource
	variables         []*commonTypes.ParsedVariable
	index             int
	currentObjectType ObjectType
	filePath          string
}

func NewParser(tokens []*frontendTypes.Token, filePath string) *Parser {
	return &Parser{
		tokens:            tokens,
		resources:         make([]*commonTypes.ParsedResource, 0),
		variables:         make([]*commonTypes.ParsedVariable, 0),
		index:             0,
		currentObjectType: Unset,
		filePath:          filePath,
	}
}

func (p *Parser) Parse() ([]*commonTypes.ParsedResource, []*commonTypes.ParsedVariable, error) {
	for p.index < len(p.tokens) {
		switch p.currentToken().Type {
		case frontendTypes.EOF:
			return p.resources, p.variables, nil
		case frontendTypes.Newline:
			p.advance()
			if p.currentToken().Type == frontendTypes.Newline {
				p.setCurrentObjectType(Unset)
			}
		case frontendTypes.ConstructorName:
			err := p.parseConstructor()
			if err != nil {
				return nil, nil, err
			}
		case frontendTypes.NaturalLanguage:
			err := p.parseNaturalLanguage()
			if err != nil {
				return nil, nil, err
			}
		default:
			return nil, nil, p.parseError(p.currentToken())
		}
	}
	return nil, nil, &errors.MissingEOFTokenError{FilePath: p.filePath}
}

func (p *Parser) currentResource() *commonTypes.ParsedResource {
	return p.resources[len(p.resources)-1]
}

func (p *Parser) currentVariable() *commonTypes.ParsedVariable {
	return p.variables[len(p.variables)-1]
}

func (p *Parser) setCurrentObjectType(t ObjectType) {
	switch t {
	case Resource:
		newResource := commonTypes.NewParsedResource(p.filePath, p.currentToken().Line)
		p.resources = append(p.resources, newResource)
	case Variable:
		newVariable := commonTypes.NewParsedVariable(p.filePath, p.currentToken().Line)
		p.variables = append(p.variables, newVariable)
	}
	p.currentObjectType = t
}

func (p *Parser) currentObjectTypeIs(objectType ObjectType) bool {
	return p.currentObjectType == objectType
}

func (p *Parser) currentObject() frontendTypes.ParsedObject {
	switch p.currentObjectType {
	case Resource:
		return p.currentResource()
	case Variable:
		return p.currentVariable()
	default:
		return nil
	}
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

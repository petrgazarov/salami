package parser

import (
	"salami/compiler/lexer"
)

type IR interface {
	IsIR()
}

type ResourceType string
type LogicalName string

type Resource struct {
	ResourceType        ResourceType
	LogicalName         LogicalName
	FreeText            string
	Uses                []LogicalName
	Exports             map[string]string
	ReferencedVariables []string
}

func (r *Resource) IsIR() {}

type Variable struct {
	Description string
	Name        string
	Value       string
}

func (v *Variable) IsIR() {}

func NewParser(tokenChannel <-chan lexer.Token, irChannel chan<- *IR) *Parser {
	return &Parser{tokenChannel: tokenChannel, irChannel: irChannel}
}

type Parser struct {
	tokenChannel <-chan lexer.Token
	irChannel    chan<- *IR
}

func (p *Parser) Parse() {
	for token := range p.tokenChannel {
		println(token.Type.String(), token.Value)
	}
}

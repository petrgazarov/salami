package semantic_analyzer

import (
	"salami/compiler/parser"
)

type SemanticAnalyzer struct {
	irChannel   <-chan *parser.IR
	symbolTable map[string]*parser.IR
}

func NewSemanticAnalyzer(irChannel <-chan *parser.IR) *SemanticAnalyzer {
	return &SemanticAnalyzer{
		irChannel:   irChannel,
		symbolTable: make(map[string]*parser.IR),
	}
}

func (sa *SemanticAnalyzer) Analyze() {
	for ir := range sa.irChannel {
		switch (*ir).(type) {
		case *parser.Resource:
			// handle Resource
		case *parser.Variable:
			// handle Variable
		}
	}
}

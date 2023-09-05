package semantic_analyzer

import (
	"fmt"
	"salami/compiler/symbol_table"
)

type SemanticAnalyzer struct {
	symbolTable *symbol_table.SymbolTable
}

func NewSemanticAnalyzer(symbolTable *symbol_table.SymbolTable) *SemanticAnalyzer {
	return &SemanticAnalyzer{
		symbolTable: symbolTable,
	}
}

func (sa *SemanticAnalyzer) Analyze() {
	fmt.Println("Resources:", sa.symbolTable.ResourceTable)
	fmt.Println("Variables:", sa.symbolTable.VariableTable)

}

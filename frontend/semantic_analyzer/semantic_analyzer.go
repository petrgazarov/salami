package semantic_analyzer

import (
	"fmt"
	"salami/common/errors"
	"salami/common/symbol_table"
)

type SemanticAnalyzer struct {
	symbolTable *symbol_table.SymbolTable
}

func NewSemanticAnalyzer(symbolTable *symbol_table.SymbolTable) *SemanticAnalyzer {
	return &SemanticAnalyzer{
		symbolTable: symbolTable,
	}
}

func (sa *SemanticAnalyzer) Analyze() error {
	if err := sa.ensureResourcesHaveAllRequiredFields(); err != nil {
		return err
	}
	if err := sa.ensureVariablesHaveAllRequiredFields(); err != nil {
		return err
	}
	if err := sa.ensureReferencedVariablesAreDefined(); err != nil {
		return err
	}
	if err := sa.ensureUsedResourcesExist(); err != nil {
		return err
	}
	return nil
}

func (sa *SemanticAnalyzer) ensureResourcesHaveAllRequiredFields() error {
	for _, resource := range sa.symbolTable.ResourceTable {
		if resource.ResourceType == "" {
			return &errors.SemanticError{
				SourceFilePath: resource.SourceFilePath,
				Message:        "Resource type field on a resource object is missing or empty",
			}
		}
		if resource.LogicalName == "" {
			return &errors.SemanticError{
				SourceFilePath: resource.SourceFilePath,
				Message:        "Logical name field on a resource object is missing or empty",
			}
		}
	}
	return nil
}

func (sa *SemanticAnalyzer) ensureVariablesHaveAllRequiredFields() error {
	for _, variable := range sa.symbolTable.VariableTable {
		if variable.Name == "" {
			return &errors.SemanticError{
				SourceFilePath: variable.SourceFilePath,
				Message:        "Name field on a variable object is missing or empty",
			}
		}
	}
	return nil
}

func (sa *SemanticAnalyzer) ensureReferencedVariablesAreDefined() error {
	for _, resource := range sa.symbolTable.ResourceTable {
		for _, variableName := range resource.ReferencedVariables {
			if _, exists := sa.symbolTable.LookupVariable(variableName); !exists {
				return &errors.SemanticError{
					SourceFilePath: resource.SourceFilePath,
					Message:        fmt.Sprintf("Referenced variable '%s' is not defined", variableName),
				}
			}
		}
	}
	return nil
}

func (sa *SemanticAnalyzer) ensureUsedResourcesExist() error {
	for _, resource := range sa.symbolTable.ResourceTable {
		for _, logicalName := range resource.Uses {
			if _, exists := sa.symbolTable.LookupResource(logicalName); !exists {
				return &errors.SemanticError{
					SourceFilePath: resource.SourceFilePath,
					Message:        fmt.Sprintf("Used resource '%s' is not defined", logicalName),
				}
			}
		}
	}
	return nil
}

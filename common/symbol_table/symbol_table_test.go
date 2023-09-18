package symbol_table_test

import (
	"salami/common/symbol_table"
	"salami/common/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSymbolTable(t *testing.T) {
	t.Run("NewSymbolTable", func(t *testing.T) {
		t.Run("should return error if logical name is not unique", func(t *testing.T) {
			resources := []*types.ParsedResource{
				{
					ResourceType:   types.ResourceType("my-resource-type"),
					LogicalName:    types.LogicalName("my-resource"),
					SourceFilePath: "path/to/my-resource.sami",
				},
				{
					ResourceType:   types.ResourceType("my-resource-type"),
					LogicalName:    types.LogicalName("my-resource"),
					SourceFilePath: "path/to/my-resource.sami",
				},
			}
			variables := []*types.ParsedVariable{}
			_, err := symbol_table.NewSymbolTable(resources, variables)
			require.EqualError(t, err, "\npath/to/my-resource.sami\n  semantic error: my-resource logical name is not unique")
		})
		t.Run("should return error if variable name is not unique", func(t *testing.T) {
			resources := []*types.ParsedResource{}
			variables := []*types.ParsedVariable{
				{
					Name:           "my-variable",
					Type:           types.VariableType("string"),
					SourceFilePath: "path/to/my-variable.sami",
				},
				{
					Name:           "my-variable",
					Type:           types.VariableType("string"),
					SourceFilePath: "path/to/my-variable.sami",
				},
			}
			_, err := symbol_table.NewSymbolTable(resources, variables)
			require.EqualError(t, err, "\npath/to/my-variable.sami\n  semantic error: my-variable variable name is not unique")
		})
	})
}

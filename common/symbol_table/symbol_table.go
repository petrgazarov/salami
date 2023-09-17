package symbol_table

import (
	"fmt"
	"salami/common/errors"
	"salami/common/types"
)

type SymbolTable struct {
	ResourceTable map[types.LogicalName]*types.ParsedResource
	VariableTable map[string]*types.ParsedVariable
}

func NewSymbolTable(resources []*types.ParsedResource, variables []*types.ParsedVariable) (*SymbolTable, error) {
	st := &SymbolTable{
		ResourceTable: make(map[types.LogicalName]*types.ParsedResource),
		VariableTable: make(map[string]*types.ParsedVariable),
	}
	for _, r := range resources {
		if _, exists := st.ResourceTable[r.LogicalName]; exists {
			return nil, &errors.SemanticError{
				SourceFilePath: r.SourceFilePath,
				Message:        fmt.Sprintf("%s logical name is not unique", r.LogicalName),
			}
		}
		st.InsertResource(r)
	}
	for _, v := range variables {
		if _, exists := st.VariableTable[v.Name]; exists {
			return nil, &errors.SemanticError{
				SourceFilePath: v.SourceFilePath,
				Message:        fmt.Sprintf("%s variable name is not unique", v.Name),
			}
		}
		st.InsertVariable(v)
	}
	return st, nil
}

func (st *SymbolTable) InsertResource(r *types.ParsedResource) {
	st.ResourceTable[r.LogicalName] = r
}

func (st *SymbolTable) InsertVariable(v *types.ParsedVariable) {
	st.VariableTable[v.Name] = v
}

func (st *SymbolTable) LookupResource(identifier types.LogicalName) (*types.ParsedResource, bool) {
	res, exists := st.ResourceTable[identifier]
	return res, exists
}

func (st *SymbolTable) LookupVariable(identifier string) (*types.ParsedVariable, bool) {
	vari, exists := st.VariableTable[identifier]
	return vari, exists
}

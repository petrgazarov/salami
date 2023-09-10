package symbol_table

import (
	"fmt"
	"salami/compiler/errors"
	"salami/compiler/types"
)

type SymbolTable struct {
	ResourceTable map[types.LogicalName]*types.Resource
	VariableTable map[string]*types.Variable
}

func NewSymbolTable(resources []*types.Resource, variables []*types.Variable) (*SymbolTable, error) {
	st := &SymbolTable{
		ResourceTable: make(map[types.LogicalName]*types.Resource),
		VariableTable: make(map[string]*types.Variable),
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

func (st *SymbolTable) InsertResource(r *types.Resource) {
	st.ResourceTable[r.LogicalName] = r
}

func (st *SymbolTable) InsertVariable(v *types.Variable) {
	st.VariableTable[v.Name] = v
}

func (st *SymbolTable) LookupResource(identifier types.LogicalName) (*types.Resource, bool) {
	res, exists := st.ResourceTable[identifier]
	return res, exists
}

func (st *SymbolTable) LookupVariable(identifier string) (*types.Variable, bool) {
	vari, exists := st.VariableTable[identifier]
	return vari, exists
}

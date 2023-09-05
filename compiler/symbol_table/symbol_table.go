package symbol_table

import (
	"salami/compiler/types"
)

type SymbolTable struct {
	ResourceTable map[types.LogicalName]types.Resource
	VariableTable map[string]types.Variable
}

func NewSymbolTable(resources []types.Resource, variables []types.Variable) *SymbolTable {
	st := &SymbolTable{
		ResourceTable: make(map[types.LogicalName]types.Resource),
		VariableTable: make(map[string]types.Variable),
	}
	for _, r := range resources {
		st.InsertResource(r)
	}
	for _, v := range variables {
		st.InsertVariable(v)
	}
	return st
}

func (st *SymbolTable) InsertResource(r types.Resource) {
	st.ResourceTable[r.LogicalName] = r
}

func (st *SymbolTable) InsertVariable(v types.Variable) {
	st.VariableTable[v.Name] = v
}

func (st *SymbolTable) LookupResource(identifier types.LogicalName) (types.Resource, bool) {
	res, exists := st.ResourceTable[identifier]
	return res, exists
}

func (st *SymbolTable) LookupVariable(identifier string) (types.Variable, bool) {
	vari, exists := st.VariableTable[identifier]
	return vari, exists
}

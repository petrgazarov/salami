package terraform

import (
	"salami/common/symbol_table"
	"salami/common/types"
)

type TerraformTarget struct {
}

func (tt *TerraformTarget) GenerateCode(*types.ChangeSet, *symbol_table.SymbolTable) []error {
	return nil
}

func (tt *TerraformTarget) GetNewFiles([]*types.Object) []*types.TargetFile {
	return nil
}

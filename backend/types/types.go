package types

import (
	"salami/common/symbol_table"
	"salami/common/types"
)

type GenerateCodeFunc func(*types.ChangeSet, *symbol_table.SymbolTable) []error
type GetFilesFromObjectsFunc func([]*types.Object) []*types.TargetFile

type Target struct {
	GenerateCode        GenerateCodeFunc
	GetFilesFromObjects GetFilesFromObjectsFunc
}

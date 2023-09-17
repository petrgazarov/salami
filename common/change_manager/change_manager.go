package change_manager

import (
	"salami/common/symbol_table"
	"salami/common/types"
)

func GenerateChangeSet(
	previousObjects []*types.Object,
	newSymbolTable *symbol_table.SymbolTable,
) (*types.ChangeSet, error) {
	return &types.ChangeSet{}, nil
}
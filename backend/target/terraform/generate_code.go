package terraform

import (
	backendTypes "salami/backend/types"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"
)

func GenerateCode(
	changeSet *commonTypes.ChangeSet,
	symbolTable *symbol_table.SymbolTable,
	llm *backendTypes.Llm,
) []error {
	return nil
}

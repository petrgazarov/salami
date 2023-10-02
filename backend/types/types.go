package types

import (
	"salami/common/symbol_table"
	"salami/common/types"
)

type Target interface {
	GenerateCode(*types.ChangeSet, *symbol_table.SymbolTable, Llm) []error
	GetFilesFromObjects([]*types.Object) []*types.TargetFile
}

type NewTargetFunc func() Target

type Llm interface {
	GetSlug() string
	GetMaxConcurrentExecutions() int
	CreateCompletion(messages []interface{}) (string, error)
}

type NewLlmFunc func(types.LlmConfig) Llm

type GetTargetLlmMessagesFunc func(*types.ChangeSetDiff, *symbol_table.SymbolTable) ([]interface{}, error)

type TargetLlmMessages struct {
	GetMessages GetTargetLlmMessagesFunc
}

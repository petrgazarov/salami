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
	CreateCompletion(messages []*LlmMessage) (string, error)
}

type NewLlmFunc func(types.LlmConfig) Llm

type LlmMessageRole string

type LlmMessage struct {
	Role    LlmMessageRole
	Content string
}

type GetTargetLlmMessagesFunc func(*types.ChangeSetDiff, *symbol_table.SymbolTable) ([]*LlmMessage, error)

type TargetLlmMessages struct {
	GetMessages GetTargetLlmMessagesFunc
}

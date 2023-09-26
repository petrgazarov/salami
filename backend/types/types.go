package types

import (
	"salami/common/symbol_table"
	"salami/common/types"
)

type GenerateCodeFunc func(*types.ChangeSet, *symbol_table.SymbolTable, *Llm) []error
type GetFilesFromObjectsFunc func([]*types.Object) []*types.TargetFile

type Target struct {
	GenerateCode        GenerateCodeFunc
	GetFilesFromObjects GetFilesFromObjectsFunc
}

type CreateCompletionFunc func(
	messages []*LlmMessage,
	llmConfig types.LlmConfig,
) (string, error)

type Llm struct {
	CreateCompletion CreateCompletionFunc
}

type LlmMessageRole string

type LlmMessage struct {
	Role    LlmMessageRole
	Content string
}

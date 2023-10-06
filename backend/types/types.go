package types

import (
	"salami/common/change_set"
	"salami/common/symbol_table"
	"salami/common/types"
)

type Target interface {
	VerifyPeerDependencies() error
	GenerateCode(*symbol_table.SymbolTable, *change_set.ChangeSetRepository, Llm) []error
	GetFilesFromObjects([]*types.Object) []*types.TargetFile
	ValidateCode([]*types.Object, *symbol_table.SymbolTable, *change_set.ChangeSetRepository, Llm) []error
}

type Llm interface {
	GetSlug() string
	GetMaxConcurrentExecutions() int
	CreateCompletion(messages []interface{}) (string, error)
}

type CodeValidationResult struct {
	ValidatedObject   *types.Object
	ErrorMessages     []string
	ReferencedObjects []*types.Object
}

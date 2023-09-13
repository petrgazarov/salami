package lock_file_manager

import (
	backendTypes "salami/backend/types"
	"salami/common/symbol_table"
	commonTypes "salami/common/types"
)

type LockFileManager struct {
}

func NewLockFileManager() *LockFileManager {
	return &LockFileManager{}
}

func ValidateTargetFilesUnchanged() error {
	return nil
}

func GenerateSourceDiff(symbolTable *symbol_table.SymbolTable) (commonTypes.SourceDiff, error) {
	return nil, nil
}

func (bm *LockFileManager) UpdateLockFile([]*backendTypes.TargetFile, []*backendTypes.BackendObject) error {
	return nil
}

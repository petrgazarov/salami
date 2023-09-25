package lock_file_manager

import (
	"log"
	"os"
	"salami/common/types"

	"github.com/BurntSushi/toml"
)

var lockFilePath = "salami-lock.toml"
var loadedLockFile *LockFile

func SetLockFilePath(path string) {
	lockFilePath = path
	loadedLockFile = nil
}

func GetTargetFilesMeta() []types.TargetFileMeta {
	targetFilesMeta := getLockFile().TargetFilesMeta
	result := make([]types.TargetFileMeta, len(targetFilesMeta))
	for i := range targetFilesMeta {
		result[i] = types.TargetFileMeta{
			FilePath: targetFilesMeta[i].FilePath,
			Checksum: targetFilesMeta[i].Checksum,
		}
	}
	return result
}

func GetObjects() []*types.Object {
	lockFileObjects := getLockFile().Objects
	result := make([]*types.Object, len(lockFileObjects))
	for i := range lockFileObjects {
		currentObject := lockFileObjects[i]
		var parsedResource *types.ParsedResource
		var parsedVariable *types.ParsedVariable
		if currentObject.IsResource() {
			parsedResource = getCommonParsedResource(currentObject)
		} else if currentObject.IsVariable() {
			parsedVariable = getCommonParsedVariable(currentObject)
		}
		codeSegments := getCommonCodeSegments(currentObject)

		result[i] = &types.Object{
			SourceFilePath: currentObject.SourceFilePath,
			ParsedResource: parsedResource,
			ParsedVariable: parsedVariable,
			CodeSegments:   codeSegments,
		}
	}
	return result
}

func UpdateLockFile(targetFilesMeta []types.TargetFileMeta, objects []*types.Object) error {
	// Merge changeSet into lockFile
	writeLockFile()
	return nil
}

func getCommonParsedResource(lockFileObject Object) *types.ParsedResource {
	uses := make([]types.LogicalName, len(lockFileObject.ParsedResource.Uses))
	for j, use := range lockFileObject.ParsedResource.Uses {
		uses[j] = types.LogicalName(use)
	}
	return &types.ParsedResource{
		ResourceType:        types.ResourceType(lockFileObject.ParsedResource.ResourceType),
		LogicalName:         types.LogicalName(lockFileObject.ParsedResource.LogicalName),
		NaturalLanguage:     lockFileObject.ParsedResource.NaturalLanguage,
		Uses:                uses,
		Exports:             lockFileObject.ParsedResource.Exports,
		ReferencedVariables: lockFileObject.ParsedResource.ReferencedVariables,
		SourceFilePath:      lockFileObject.SourceFilePath,
	}
}

func getCommonParsedVariable(lockFileObject Object) *types.ParsedVariable {
	return &types.ParsedVariable{
		Description:    lockFileObject.ParsedVariable.Description,
		Name:           lockFileObject.ParsedVariable.Name,
		Default:        lockFileObject.ParsedVariable.DefaultValue,
		Type:           types.VariableType(lockFileObject.ParsedVariable.VariableType),
		SourceFilePath: lockFileObject.SourceFilePath,
	}
}

func getCommonCodeSegments(lockFileObject Object) []types.CodeSegment {
	codeSegments := make([]types.CodeSegment, len(lockFileObject.CodeSegments))
	for j, segment := range lockFileObject.CodeSegments {
		codeSegments[j] = types.CodeSegment{
			SegmentType:    types.CodeSegmentType(segment.SegmentType),
			Content:        segment.Content,
		}
	}
	return codeSegments
}

func writeLockFile() error {
	file, err := os.OpenFile(lockFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(getLockFile()); err != nil {
		return err
	}
	return nil
}

func getLockFile() *LockFile {
	if loadedLockFile == nil {
		log.Fatal("Lock file not loaded")
	}
	return loadedLockFile
}

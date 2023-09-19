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
	objects := getLockFile().Objects
	result := make([]*types.Object, len(objects))
	for i := range objects {
		var parsed types.ParsedObject
		switch objects[i].Parsed.getObjectType() {
		case "Resource":
			parsedResource := objects[i].Parsed.(*ParsedResource)
			uses := make([]types.LogicalName, len(parsedResource.Uses))
			for j, use := range parsedResource.Uses {
				uses[j] = types.LogicalName(use)
			}
			parsed = &types.ParsedResource{
				ResourceType:        types.ResourceType(parsedResource.ResourceType),
				LogicalName:         types.LogicalName(parsedResource.LogicalName),
				NaturalLanguage:     parsedResource.NaturalLanguage,
				Uses:                uses,
				Exports:             parsedResource.Exports,
				ReferencedVariables: parsedResource.ReferencedVariables,
				SourceFilePath:      objects[i].SourceFilePath,
			}
		case "Variable":
			parsedVariable := objects[i].Parsed.(*ParsedVariable)
			parsed = &types.ParsedVariable{
				Description:    parsedVariable.Description,
				Name:           parsedVariable.Name,
				Default:        parsedVariable.DefaultValue,
				Type:           types.VariableType(parsedVariable.VariableType),
				SourceFilePath: objects[i].SourceFilePath,
			}
		}

		codeSegments := make([]types.CodeSegment, len(objects[i].CodeSegments))
		for j, segment := range objects[i].CodeSegments {
			codeSegments[j] = types.CodeSegment{
				SegmentType:    types.CodeSegmentType(segment.SegmentType),
				Content:        segment.Content,
				TargetFilePath: segment.TargetFilePath,
			}
		}

		result[i] = &types.Object{
			SourceFilePath: objects[i].SourceFilePath,
			Parsed:         parsed,
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

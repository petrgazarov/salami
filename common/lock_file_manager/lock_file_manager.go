package lock_file_manager

import (
	"log"
	"os"
	commonTypes "salami/common/types"

	"github.com/BurntSushi/toml"
)

const lockFilePath = "salami-lock.toml"

var loadedLockFile *LockFile

func GetTargetFilesMeta() []*commonTypes.TargetFileMeta {
	targetFilesMeta := getLockFile().TargetFilesMeta
	result := make([]*commonTypes.TargetFileMeta, len(targetFilesMeta))
	for i := range targetFilesMeta {
		result[i] = &commonTypes.TargetFileMeta{
			FilePath: targetFilesMeta[i].FilePath,
			Checksum: targetFilesMeta[i].Checksum,
		}
	}
	return result
}

func GetObjects() []*commonTypes.Object {
	objects := getLockFile().Objects
	result := make([]*commonTypes.Object, len(objects))
	for i := range objects {
		var parsed commonTypes.ParsedObject
		switch objects[i].Parsed.getObjectType() {
		case "Resource":
			parsedResource := objects[i].Parsed.(*ParsedResource)
			uses := make([]commonTypes.LogicalName, len(parsedResource.Uses))
			for j, use := range parsedResource.Uses {
				uses[j] = commonTypes.LogicalName(use)
			}
			parsed = &commonTypes.ParsedResource{
				ResourceType:        commonTypes.ResourceType(parsedResource.ResourceType),
				LogicalName:         commonTypes.LogicalName(parsedResource.LogicalName),
				NaturalLanguage:     parsedResource.NaturalLanguage,
				Uses:                uses,
				Exports:             parsedResource.Exports,
				ReferencedVariables: parsedResource.ReferencedVariables,
				SourceFilePath:      objects[i].SourceFilePath,
			}
		case "Variable":
			parsedVariable := objects[i].Parsed.(*ParsedVariable)
			parsed = &commonTypes.ParsedVariable{
				Description:    parsedVariable.Description,
				Name:           parsedVariable.Name,
				Default:        parsedVariable.Default,
				Type:           commonTypes.VariableType(parsedVariable.Type),
				SourceFilePath: objects[i].SourceFilePath,
			}
		}

		codeSegments := make([]commonTypes.CodeSegment, len(objects[i].CodeSegments))
		for j, segment := range objects[i].CodeSegments {
			codeSegments[j] = commonTypes.CodeSegment{
				SegmentType:    commonTypes.CodeSegmentType(segment.SegmentType),
				Content:        segment.Content,
				TargetFilePath: segment.TargetFilePath,
			}
		}

		result[i] = &commonTypes.Object{
			SourceFilePath: objects[i].SourceFilePath,
			Parsed:         parsed,
			CodeSegments:   codeSegments,
		}
	}
	return result
}

func UpdateLockFile(changeSet *commonTypes.ChangeSet) error {
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
		if _, err := toml.DecodeFile(lockFilePath, loadedLockFile); err != nil {
			if err != nil && !os.IsNotExist(err) {
				log.Fatalf("failed to read lock file")
			} else {
				loadedLockFile = &LockFile{}
			}
		}
	}
	return loadedLockFile
}

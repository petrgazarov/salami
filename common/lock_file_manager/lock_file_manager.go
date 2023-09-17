package lock_file_manager

import (
	"log"
	"os"
	"salami/common/types"
	commonTypes "salami/common/types"

	"github.com/BurntSushi/toml"
)

const lockFilePath = "salami-lock.toml"

var loadedLockFile *lockFile

func GetTargetFilesMeta() []commonTypes.TargetFileMeta {
	targetFilesMeta := getLockFile().targetFilesMeta
	result := make([]commonTypes.TargetFileMeta, len(targetFilesMeta))
	for i := range targetFilesMeta {
		result[i] = commonTypes.TargetFileMeta{
			FilePath: targetFilesMeta[i].filePath,
			Checksum: targetFilesMeta[i].checksum,
		}
	}
	return result
}

func GetObjects() []*commonTypes.Object {
	objects := getLockFile().objects
	result := make([]*commonTypes.Object, len(objects))
	for i := range objects {
		var parsed commonTypes.ParsedObject
		switch objects[i].parsed.getObjectType() {
		case "Resource":
			parsedResource := objects[i].parsed.(*parsedResource)
			uses := make([]commonTypes.LogicalName, len(parsedResource.uses))
			for j, use := range parsedResource.uses {
				uses[j] = commonTypes.LogicalName(use)
			}
			parsed = &commonTypes.ParsedResource{
				ResourceType:        commonTypes.ResourceType(parsedResource.resourceType),
				LogicalName:         commonTypes.LogicalName(parsedResource.logicalName),
				NaturalLanguage:     parsedResource.naturalLanguage,
				Uses:                uses,
				Exports:             parsedResource.exports,
				ReferencedVariables: parsedResource.referencedVariables,
				SourceFilePath:      objects[i].sourceFilePath,
			}
		case "Variable":
			parsedVariable := objects[i].parsed.(*parsedVariable)
			parsed = &commonTypes.ParsedVariable{
				Description:    parsedVariable.description,
				Name:           parsedVariable.name,
				Default:        parsedVariable.defaultValue,
				Type:           commonTypes.VariableType(parsedVariable.variableType),
				SourceFilePath: objects[i].sourceFilePath,
			}
		}

		codeSegments := make([]commonTypes.CodeSegment, len(objects[i].codeSegments))
		for j, segment := range objects[i].codeSegments {
			codeSegments[j] = commonTypes.CodeSegment{
				SegmentType:    commonTypes.CodeSegmentType(segment.segmentType),
				Content:        segment.content,
				TargetFilePath: segment.targetFilePath,
			}
		}

		result[i] = &commonTypes.Object{
			SourceFilePath: objects[i].sourceFilePath,
			Parsed:         parsed,
			CodeSegments:   codeSegments,
		}
	}
	return result
}

func UpdateLockFile(targetFiles []types.TargetFileMeta, newObjects []*types.Object) error {
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

func getLockFile() *lockFile {
	if loadedLockFile == nil {
		if _, err := toml.DecodeFile(lockFilePath, loadedLockFile); err != nil {
			if err != nil && !os.IsNotExist(err) {
				log.Fatalf("failed to read lock file")
			} else {
				loadedLockFile = &lockFile{}
			}
		}
	}
	return loadedLockFile
}

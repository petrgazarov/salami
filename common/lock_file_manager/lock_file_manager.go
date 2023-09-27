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

func GetTargetFileMetas() []types.TargetFileMeta {
	targetFileMetas := getLockFile().TargetFileMetas
	result := make([]types.TargetFileMeta, len(targetFileMetas))
	for i := range targetFileMetas {
		result[i] = types.TargetFileMeta{
			FilePath: targetFileMetas[i].FilePath,
			Checksum: targetFileMetas[i].Checksum,
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

		result[i] = &types.Object{
			ParsedResource: parsedResource,
			ParsedVariable: parsedVariable,
			TargetCode:     currentObject.TargetCode,
		}
	}
	return result
}

func UpdateLockFile(targetFileMetas []types.TargetFileMeta, objects []*types.Object) error {
	// Merge changeSet into lockFile
	writeLockFile()
	return nil
}

func getCommonParsedResource(lockFileObject Object) *types.ParsedResource {
	referencedResources := make([]types.LogicalName, len(lockFileObject.ParsedResource.ReferencedResources))
	for j, use := range lockFileObject.ParsedResource.ReferencedResources {
		referencedResources[j] = types.LogicalName(use)
	}
	return &types.ParsedResource{
		ResourceType:        types.ResourceType(lockFileObject.ParsedResource.ResourceType),
		LogicalName:         types.LogicalName(lockFileObject.ParsedResource.LogicalName),
		NaturalLanguage:     lockFileObject.ParsedResource.NaturalLanguage,
		ReferencedResources: referencedResources,
		ReferencedVariables: lockFileObject.ParsedResource.ReferencedVariables,
		SourceFilePath:      lockFileObject.ParsedResource.SourceFilePath,
		SourceFileLine:      lockFileObject.ParsedResource.SourceFileLine,
	}
}

func getCommonParsedVariable(lockFileObject Object) *types.ParsedVariable {
	return &types.ParsedVariable{
		Name:            lockFileObject.ParsedVariable.Name,
		NaturalLanguage: lockFileObject.ParsedVariable.NaturalLanguage,
		Default:         lockFileObject.ParsedVariable.Default,
		Type:            types.VariableType(lockFileObject.ParsedVariable.VariableType),
		SourceFilePath:  lockFileObject.ParsedVariable.SourceFilePath,
		SourceFileLine:  lockFileObject.ParsedVariable.SourceFileLine,
	}
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

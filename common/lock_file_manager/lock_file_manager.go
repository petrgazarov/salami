package lock_file_manager

import (
	"bytes"
	"log"
	"salami/common/types"
	"salami/common/constants"
	"salami/common/utils/file_utils"

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
			parsedResource = lockFileToCommonResource(currentObject)
		} else if currentObject.IsVariable() {
			parsedVariable = lockFileToCommonVariable(currentObject)
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
	lockFile := getLockFile()
	lockFile.Version = constants.SalamiVersion
	lockFile.TargetFileMetas = make([]TargetFileMeta, len(targetFileMetas))
	lockFile.Objects = make([]Object, len(objects))
	for i := range targetFileMetas {
		lockFile.TargetFileMetas[i] = TargetFileMeta{
			FilePath: targetFileMetas[i].FilePath,
			Checksum: targetFileMetas[i].Checksum,
		}
	}
	for i := range objects {
		currentObject := objects[i]

		lockFile.Objects[i] = Object{
			ParsedResource: commonToLockFileResource(currentObject.ParsedResource),
			ParsedVariable: commonToLockFileVariable(currentObject.ParsedVariable),
			TargetCode:     currentObject.TargetCode,
		}
	}

	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	if err := encoder.Encode(lockFile); err != nil {
		return err
	}
	tomlString := buf.String()

	if err := file_utils.WriteFileIfChanged(lockFilePath, tomlString); err != nil {
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

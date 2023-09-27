package terraform

import (
	"salami/common/constants"
	"salami/common/types"
	"strings"
)

const terraformFileExtension = ".tf"

func GetFilesFromObjects(objects []*types.Object) []*types.TargetFile {
	targetFiles := []*types.TargetFile{}

	var currentTargetFile *types.TargetFile
	var lastSourceFilePath string
	for _, object := range objects {
		if object.GetSourceFilePath() != lastSourceFilePath {
			currentTargetFile = &types.TargetFile{
				FilePath: getTargetFilePath(object.GetSourceFilePath()),
				Content:  object.TargetCode,
			}
			targetFiles = append(targetFiles, currentTargetFile)
			lastSourceFilePath = object.GetSourceFilePath()
		} else {
			currentTargetFile.Content += ("\n\n" + object.TargetCode)
		}
	}

	return targetFiles
}

func getTargetFilePath(sourceFilePath string) string {
	if strings.HasSuffix(sourceFilePath, constants.SalamiFileExtension) {
		return strings.TrimSuffix(sourceFilePath, constants.SalamiFileExtension) + terraformFileExtension
	}
	return sourceFilePath
}

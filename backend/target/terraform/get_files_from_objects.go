package terraform

import (
	"salami/common/types"
	"strings"
)

func GetFilesFromObjects(objects []*types.Object) []*types.TargetFile {
	targetFiles := []*types.TargetFile{}

	var currentTargetFile *types.TargetFile
	var lastSourceFilePath string
	for _, object := range objects {
		currentContent := ""
		for _, codeSegment := range object.CodeSegments {
			currentContent += codeSegment.Content
		}

		if object.GetSourceFilePath() != lastSourceFilePath {
			currentTargetFile = &types.TargetFile{
				FilePath: getTargetFilePath(object.GetSourceFilePath()),
				Content:  currentContent,
			}
			targetFiles = append(targetFiles, currentTargetFile)
			lastSourceFilePath = object.GetSourceFilePath()
		} else {
			currentTargetFile.Content += ("\n\n" + currentContent)
		}
	}

	return targetFiles
}

func getTargetFilePath(sourceFilePath string) string {
	if strings.HasSuffix(sourceFilePath, ".sami") {
		return strings.TrimSuffix(sourceFilePath, ".sami") + ".tf"
	}
	return sourceFilePath
}

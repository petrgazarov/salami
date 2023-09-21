package lock_file_manager

import (
	"fmt"
	"regexp"
)

var fieldNameMap = map[string]string{
	"LockFile.Version":                                   "lock file version",
	"LockFile.TargetFilesMeta[*].FilePath":               "target file path",
	"LockFile.TargetFilesMeta[*].Checksum":               "target file checksum",
	"LockFile.Objects[*].SourceFilePath":                 "source file path",
	"LockFile.Objects[*].ParsedResource":                 "parsed resource",
	"LockFile.Objects[*].ParsedVariable":                 "parsed variable",
	"LockFile.Objects[*].ParsedResource.ResourceType":    "parsed resource type",
	"LockFile.Objects[*].ParsedResource.LogicalName":     "parsed resource logical name",
	"LockFile.Objects[*].ParsedVariable.Name":            "parsed variable name",
	"LockFile.Objects[*].ParsedVariable.VariableType":    "parsed variable type",
	"LockFile.Objects[*].CodeSegments[*].SegmentType":    "code segment type",
	"LockFile.Objects[*].CodeSegments[*].TargetFilePath": "code segment target file path",
	"LockFile.Objects[*].CodeSegments[*].Content":        "code segment content",
}

type LockFileError struct {
	Message string
}

func (e *LockFileError) Error() string {
	return fmt.Sprintf("lock file error: %s", e.Message)
}

func getMissingFieldError(field string) error {
	re := regexp.MustCompile(`\[\d+\]`)
	field = re.ReplaceAllStringFunc(field, func(s string) string {
		return "[*]"
	})

	fieldName := fieldNameMap[field]
	if fieldName == "" {
		fieldName = field
	}
	return &LockFileError{Message: fmt.Sprintf("missing or invalid %s", fieldName)}
}

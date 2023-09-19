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
	"LockFile.Objects[*].RawParsed":                      "parsed object",
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

	fmt.Println("field", field)
	fieldName := fieldNameMap[field]
	return &LockFileError{Message: fmt.Sprintf("missing %s", fieldName)}
}

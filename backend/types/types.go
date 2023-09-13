package types

import commonTypes "salami/common/types"

type CodeSegmentType string

type CodeSegment struct {
	Type      CodeSegmentType
	Content   string
	FilePath  string
	LineStart int
	LineEnd   int
}

type BackendObject struct {
	SourceFilePath string
	ParsedObject   commonTypes.Object
	CodeSegments   []CodeSegment
}

type TargetFile struct {
	FilePath string
	Checksum string
}

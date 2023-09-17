package types

type CodeSegmentType string

type CodeSegment struct {
	SegmentType    CodeSegmentType
	Content        string
	TargetFilePath string
}

type Object struct {
	SourceFilePath string
	Parsed         ParsedObject
	CodeSegments   []CodeSegment
}

type ChangeSetDiff struct {
	OldObject *Object
	NewObject *Object
}

type ChangeSet struct {
	Diffs []ChangeSetDiff
}

type TargetFileMeta struct {
	FilePath string
	Checksum string
}

type TargetConfig struct {
	Platform string
}

type LlmConfig struct {
	Provider string
	Model    string
}

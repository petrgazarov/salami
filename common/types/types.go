package types

const SalamiFileExtension = ".sami"

type ChangeSetParsedObject interface{}

type CodeSegmentType string

type CodeSegment struct {
	Type           CodeSegmentType
	Content        string
	TargetFilePath string
}

type ChangeSetObject struct {
	SourceFilePath string
	Parsed         ChangeSetParsedObject
	CodeSegments   []CodeSegment
}

type ChangeSetDiff struct {
	OldObject *ChangeSetObject
	NewObject *ChangeSetObject
}

type ChangeSet struct {
	Diffs []ChangeSetDiff
}

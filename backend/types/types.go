package types

import commonTypes "salami/common/types"

type BackendObject struct {
	SourceFilePath string
	ParsedObject   *commonTypes.ParsedObject
	CodeSegments   []*commonTypes.CodeSegment
}

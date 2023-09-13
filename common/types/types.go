package types

type SourceDiff struct {
	Added   []Object
	Removed []Object
	Changed []Object
}

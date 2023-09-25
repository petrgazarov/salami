package types

type ResourceType string
type LogicalName string

type ParsedResource struct {
	ResourceType        ResourceType
	LogicalName         LogicalName
	NaturalLanguage     string
	Uses                []LogicalName
	Exports             map[string]string
	ReferencedVariables []string
	SourceFilePath      string
	SourceFileLine      int
}

func NewParsedResource(SourceFilePath string, SourceFileLine int) *ParsedResource {
	return &ParsedResource{
		Uses:                []LogicalName{},
		Exports:             make(map[string]string),
		ReferencedVariables: []string{},
		SourceFilePath:      SourceFilePath,
		SourceFileLine:      SourceFileLine,
	}
}

type VariableType string

type ParsedVariable struct {
	Description    string
	Name           string
	Default        string
	Type           VariableType
	SourceFilePath string
	SourceFileLine int
}

func NewParsedVariable(SourceFilePath string, SourceFileLine int) *ParsedVariable {
	return &ParsedVariable{
		SourceFilePath: SourceFilePath,
		SourceFileLine: SourceFileLine,
	}
}

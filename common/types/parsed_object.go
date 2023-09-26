package types

type ParsedObject interface {
	AddNaturalLanguage(string)
}

type ResourceType string
type LogicalName string

type ParsedResource struct {
	ResourceType        ResourceType
	LogicalName         LogicalName
	NaturalLanguage     string
	ReferencedResources []LogicalName
	ReferencedVariables []string
	SourceFilePath      string
	SourceFileLine      int
}

func NewParsedResource(SourceFilePath string, SourceFileLine int) *ParsedResource {
	return &ParsedResource{
		ReferencedResources: []LogicalName{},
		ReferencedVariables: []string{},
		SourceFilePath:      SourceFilePath,
		SourceFileLine:      SourceFileLine,
	}
}

func (r *ParsedResource) AddNaturalLanguage(NaturalLanguage string) {
	if r.NaturalLanguage == "" {
		r.NaturalLanguage = NaturalLanguage
	} else {
		r.NaturalLanguage = r.NaturalLanguage + "\n" + NaturalLanguage
	}
}

type VariableType string

type ParsedVariable struct {
	Name            string
	NaturalLanguage string
	Default         string
	Type            VariableType
	SourceFilePath  string
	SourceFileLine  int
}

func NewParsedVariable(SourceFilePath string, SourceFileLine int) *ParsedVariable {
	return &ParsedVariable{
		SourceFilePath: SourceFilePath,
		SourceFileLine: SourceFileLine,
	}
}

func (v *ParsedVariable) AddNaturalLanguage(NaturalLanguage string) {
	if v.NaturalLanguage == "" {
		v.NaturalLanguage = NaturalLanguage
	} else {
		v.NaturalLanguage = v.NaturalLanguage + "\n" + NaturalLanguage
	}
}

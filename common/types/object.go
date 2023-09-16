package types

import "fmt"

type ParsedObject interface {
	GetSourceFileLine() int
	GetSourceFilePath() string
}

type ResourceType string
type LogicalName string

type Resource struct {
	ResourceType        ResourceType
	LogicalName         LogicalName
	NaturalLanguage     string
	Uses                []LogicalName
	Exports             map[string]string
	ReferencedVariables []string
	SourceFilePath      string
	SourceFileLine      int
}

func (r *Resource) GetSourceFileLine() int {
	return r.SourceFileLine
}

func (r *Resource) GetSourceFilePath() string {
	return r.SourceFilePath
}

func NewResource(SourceFilePath string, SourceFileLine int) *Resource {
	return &Resource{
		Uses:                []LogicalName{},
		Exports:             make(map[string]string),
		ReferencedVariables: []string{},
		SourceFilePath:      SourceFilePath,
		SourceFileLine:      SourceFileLine,
	}
}

type VariableType int

const (
	String VariableType = iota
	Number
	Boolean
)

var VariableTypeNames = map[VariableType]string{
	String:  "string",
	Number:  "number",
	Boolean: "boolean",
}

func (v VariableType) String() string {
	return VariableTypeNames[v]
}

func StringToVariableType(s string) (VariableType, error) {
	for k, v := range VariableTypeNames {
		if v == s {
			return k, nil
		}
	}
	return 0, fmt.Errorf("invalid variable type %s", s)
}

type Variable struct {
	Description    string
	Name           string
	Default        string
	Type           VariableType
	SourceFilePath string
	SourceFileLine int
}

func (v *Variable) GetSourceFileLine() int {
	return v.SourceFileLine
}

func (v *Variable) GetSourceFilePath() string {
	return v.SourceFilePath
}

func NewVariable(SourceFilePath string, SourceFileLine int) *Variable {
	return &Variable{
		SourceFilePath: SourceFilePath,
		SourceFileLine: SourceFileLine,
	}
}

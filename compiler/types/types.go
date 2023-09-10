package types

import "fmt"

type TokenType int

const (
	DecoratorName TokenType = iota
	DecoratorArg
	FieldName
	FieldValue
	NaturalLanguage
	VariableRef
	Newline
	EOF
	Error
)

var TokenTypeNames = map[TokenType]string{
	DecoratorName:   "DecoratorName",
	DecoratorArg:    "DecoratorArg",
	FieldName:       "FieldName",
	FieldValue:      "FieldValue",
	NaturalLanguage: "NaturalLanguage",
	VariableRef:     "VariableRef",
	Newline:         "Newline",
	EOF:             "EOF",
	Error:           "Error",
}

func (t TokenType) String() string {
	return TokenTypeNames[t]
}

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
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
}

func NewResource(SourceFilePath string) *Resource {
	return &Resource{
		Uses:                []LogicalName{},
		Exports:             make(map[string]string),
		ReferencedVariables: []string{},
		SourceFilePath:      SourceFilePath,
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
	Value          string
	Type           VariableType
	SourceFilePath string
}

func NewVariable(SourceFilePath string) *Variable {
	return &Variable{
		SourceFilePath: SourceFilePath,
	}
}

var ValidFieldNames = map[string]bool{
	"Resource type": true,
	"Logical name":  true,
	"Description":   true,
	"Name":          true,
	"Value":         true,
}

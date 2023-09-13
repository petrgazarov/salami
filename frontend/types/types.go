package types

type TokenType int

const (
	DecoratorName TokenType = iota
	DecoratorArg
	FieldName
	FieldValue
	NaturalLanguage
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

var ValidFieldNames = map[string]bool{
	"Resource type": true,
	"Logical name":  true,
	"Description":   true,
	"Name":          true,
	"Default":       true,
}

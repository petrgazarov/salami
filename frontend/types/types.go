package types

type TokenType int

const (
	ConstructorName TokenType = iota
	ConstructorArg
	NaturalLanguage
	Newline
	EOF
	Error
)

var TokenTypeNames = map[TokenType]string{
	ConstructorName: "ConstructorName",
	ConstructorArg:  "ConstructorArg",
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

type ParsedObject interface {
	AddNaturalLanguage(string)
}

package types

type TokenType int

const (
	DecoratorName TokenType = iota
	DecoratorArg
	NaturalLanguage
	Newline
	EOF
	Error
)

var TokenTypeNames = map[TokenType]string{
	DecoratorName:   "DecoratorName",
	DecoratorArg:    "DecoratorArg",
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

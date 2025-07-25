package tokentype

type TokenType uint8

const (
	None TokenType = iota
	String
	Number
	Boolean
	Null
	Symbol
)

var tokenTypeNames = [...]string{
	"None",
	"String",
	"Number",
	"Boolean",
	"Null",
	"Symbol",
}

func (t TokenType) String() string {
	if int(t) >= len(tokenTypeNames) {
		return "None"
	}
	return tokenTypeNames[t]
}

func (t TokenType) IsPrimitive() bool {
	return t != Symbol
}

package tokentype

type TokenType uint8

const (
	None TokenType = iota
	String
	Number
	Boolean
	Null
	Array
	Object
)

var tokenTypeNames = [...]string{
	"None",
	"String",
	"Number",
	"Boolean",
	"Null",
	"Array",
	"Object",
}

func (t TokenType) String() string {
	if int(t) >= len(tokenTypeNames) {
		return "None"
	}
	return tokenTypeNames[t]
}

func (t TokenType) IsPrimitive() bool {
	return t != Object && t != Array
}

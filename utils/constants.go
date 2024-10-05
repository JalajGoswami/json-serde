package utils

type TokenType uint8

var tokenTypeNames = []string{"None", "String", "Number", "Boolean", "Null", "Array", "Object"}

func (t TokenType) String() string {
	return tokenTypeNames[t]
}

const (
	None TokenType = iota
	String
	Number
	Boolean
	Null
	Array
	Object
)

func IsPrimitiveType(t TokenType) bool {
	return t != Object && t != Array
}

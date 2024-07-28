package utils

type TokenType = uint8

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

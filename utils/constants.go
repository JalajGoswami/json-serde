package utils

type DataType = uint8
type TokenType = uint8

const (
	NoType DataType = iota
	StringType
	NumberType
	BooleanType
	NullType
	ArrayType
	ObjectType
)

const (
	None TokenType = iota
	String
	Number
	Boolean
	Null
	Array
	Object
)

func IsPrimitiveType(t DataType) bool {
	return t != ObjectType && t != ArrayType
}

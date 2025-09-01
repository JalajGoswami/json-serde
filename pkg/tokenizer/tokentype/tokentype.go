package tokentype

type TokenType uint8

const (
	None TokenType = iota
	Null
	Boolean
	Number
	String
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

type SymbolType uint8

const (
	NoSymbol SymbolType = iota
	Comma
	Colon
	BracketOpen
	BracketClose
	BraceOpen
	BraceClose
)

var symbolTypeNames = [...]string{
	"No Symbol",
	"Comma",
	"Colon",
	"Opening Bracket",
	"Closing Bracket",
	"Opening Brace",
	"Closing Brace",
}

func (t SymbolType) String() string {
	if int(t) >= len(symbolTypeNames) {
		return "NoSymbol"
	}
	return symbolTypeNames[t]
}

var symbols = [...]byte{
	',',
	':',
	'[',
	']',
	'{',
	'}',
}

func SymbolFromByte(b byte) SymbolType {
	for i, s := range symbols {
		if s == b {
			return SymbolType(i + 1)
		}
	}
	return NoSymbol
}

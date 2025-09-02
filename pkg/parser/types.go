package parser

type NodeType uint

const (
	NullNode = iota
	BooleanNode
	NumberNode
	StringNode
	ArrayNode
	ObjectNode
	PropertyNode
)

type Node struct {
	Type     NodeType
	Value    any
	Key      string // present in property node
	Children []*Node
}

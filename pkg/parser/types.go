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

var nodeTypeNames = [...]string{
	"Null",
	"Boolean",
	"Number",
	"String",
	"Array",
	"Object",
	"Property",
}

func (t NodeType) String() string {
	if int(t) >= len(nodeTypeNames) {
		return "Null"
	}
	return nodeTypeNames[t]
}

func (t NodeType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

type Node struct {
	Type     NodeType
	Value    any     `json:",omitempty"`
	Key      string  `json:",omitempty"` // present in property node
	Children []*Node `json:",omitempty"`
}

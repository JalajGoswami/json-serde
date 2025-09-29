package parser

// transforms tree representation into corresponding native data type
func decodeTree(root *Node) any {
	return decodeNode(root)
}

func decodeNode(n *Node) any {
	switch n.Type {
	case NullNode, BooleanNode, NumberNode, StringNode:
		return n.Value

	case ObjectNode:
		var object = make(map[string]any, len(n.Children))
		for _, property := range n.Children {
			value, ok := property.Value.(*Node)
			if !ok {
				panic("unknown value of a property in object")
			}
			object[property.Key] = decodeNode(value)
		}
		return object
	case ArrayNode:
		var array = make([]any, len(n.Children))
		for i, element := range n.Children {
			array[i] = decodeNode(element)
		}
		return array
	}
	return nil
}

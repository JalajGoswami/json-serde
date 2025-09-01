// Package parser parses stream of json tokens to generate corresponding native type
package parser

import (
	"fmt"
	"json-serde/pkg/tokenizer"
	"reflect"
)

type Parser struct {
	tk tokenizer.Tokenizer
}

func (p *Parser) Parse(v any) error {
	if err := checkValidInput(v); err != nil {
		return err
	}

	valueType := reflect.TypeOf(v).Elem()
	switch valueType.Kind() {
	case reflect.Interface:
		err := p.buildTree()
		if err != nil {
			return err
		}
	}

	return nil
}

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

func (p *Parser) buildTree() error {
	// iterate over tokens got from tokenizer and then generate a parse tree for json
	_, err := p.tk.Next()
	if err != nil {

	}
	return nil
}

func NewParser(tk tokenizer.Tokenizer) Parser {
	return Parser{tk: tk}
}

func checkValidInput(v any) error {
	valueType := reflect.TypeOf(v)
	if valueType == nil {
		return ErrNilInput
	}
	if valueType.Kind() != reflect.Pointer {
		return fmt.Errorf("%w %v", ErrNonPointer, valueType)
	}
	valueType = valueType.Elem() // dereference the pointer
	if valueType != reflect.TypeOf((*any)(nil)).Elem() {
		return fmt.Errorf("cannot parse json into %v", valueType)
	}
	return nil
}

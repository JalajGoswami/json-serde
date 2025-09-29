// Package parser parses stream of json tokens to generate corresponding native type
package parser

import (
	"fmt"
	"json-serde/pkg/tokenizer"
	"json-serde/pkg/tokenizer/tokentype"
	"reflect"
	"slices"
	"strconv"
)

type Parser struct {
	tk tokenizer.Tokenizer
}

func (p *Parser) Parse(v any) error {
	if err := checkValidInput(v); err != nil {
		return err
	}

	value := reflect.ValueOf(v).Elem()
	switch value.Kind() {
	case reflect.Interface:
		parseTree, err := p.buildTree()
		if err != nil {
			return err
		}

		parsed := decodeTree(parseTree)
		value.Set(reflect.ValueOf(parsed))
	}

	return nil
}

func (p *Parser) buildTree() (*Node, error) {
	// iterate over tokens got from tokenizer and then generate a parse tree for json
	root, err := p.tk.Next()
	if err != nil {
		return nil, err
	}
	parseTree, err := p.parseToken(root)
	return parseTree, err
}

func (p *Parser) parseToken(t *tokenizer.Token) (*Node, error) {
	switch t.TokenType {

	case tokentype.Null:
		return &Node{Type: NullNode, Value: nil}, nil

	case tokentype.Boolean:
		return &Node{Type: BooleanNode, Value: t.Value[0] == 't'}, nil

	case tokentype.Number:
		number := &Node{Type: NumberNode}
		if isFractional(t.Value) {
			number.Value, _ = strconv.ParseFloat(string(t.Value), 64)
		} else {
			number.Value, _ = strconv.ParseInt(string(t.Value), 10, 64)
		}
		return number, nil

	case tokentype.String:
		return &Node{Type: StringNode, Value: string(t.Value)}, nil

	case tokentype.Symbol:
		if t.SymbolType == tokentype.BraceOpen {
			return p.parseObject(t)
		}
		if t.SymbolType == tokentype.BracketOpen {
			return p.parseArray(t)
		}

	}
	return nil, ErrUnexpectedToken
}

func (p *Parser) parseObject(t *tokenizer.Token) (*Node, error) {
	if t.SymbolType != tokentype.BraceOpen {
		return nil, fmt.Errorf("opening brace expected")
	}
	object := &Node{Type: ObjectNode}

	trailingComma := false
	for {
		token, err := p.tk.Next()
		if err != nil {
			return nil, fmt.Errorf("%w, expected property or closing brace", err)
		}

		if token.TokenType == tokentype.String {
			trailingComma = false
			property := &Node{Type: PropertyNode}
			property.Key = string(token.Value)

			token, err = p.tk.Next()
			if err != nil {
				return nil, fmt.Errorf("%w, expected colon", err)
			}
			if token.SymbolType != tokentype.Colon {
				return nil, fmt.Errorf("found %v, expected colon", token.TokenType)
			}

			token, err = p.tk.Next()
			if err != nil {
				return nil, fmt.Errorf("%w, expected value", err)
			}
			property.Value, err = p.parseToken(token)
			if err != nil {
				if err == ErrUnexpectedToken {
					return nil, fmt.Errorf("found %v, expected value", token.SymbolType)
				}
				return nil, err
			}
			object.Children = append(object.Children, property)

			token, err = p.tk.Next()
			if err != nil {
				return nil, fmt.Errorf("%w, expected comma or closing brace", err)
			}

			if token.SymbolType == tokentype.Comma {
				trailingComma = true
				continue
			} else if token.SymbolType == tokentype.BraceClose {
				return object, nil
			} else {
				return nil, fmt.Errorf("found %v, expected comma or closing brace", token.TokenType)
			}

		} else if token.SymbolType == tokentype.BraceClose {
			if trailingComma {
				return nil, ErrTrailingComma
			}
			return object, nil
		} else {
			return nil, fmt.Errorf("found %v, expected property or closing brace", token.TokenType)
		}
	}
}

func (p *Parser) parseArray(t *tokenizer.Token) (*Node, error) {
	if t.SymbolType != tokentype.BracketOpen {
		return nil, fmt.Errorf("opening bracket expected")
	}
	trailingComma := false
	array := &Node{Type: ArrayNode}
	for {
		token, err := p.tk.Next()
		if err != nil {
			return nil, fmt.Errorf("%w, expected value or closing bracket", err)
		}

		if token.SymbolType == tokentype.BracketClose {
			if trailingComma {
				return nil, ErrTrailingComma
			}
			return array, nil
		}
		trailingComma = false

		value, err := p.parseToken(token)
		if err != nil {
			if err == ErrUnexpectedToken {
				return nil, fmt.Errorf("found %v, expected value", token.SymbolType)
			}
			return nil, err
		}
		array.Children = append(array.Children, value)

		token, err = p.tk.Next()
		if err != nil {
			return nil, fmt.Errorf("%w, expected comma or closing bracket", err)
		}

		if token.SymbolType == tokentype.Comma {
			trailingComma = true
			continue
		} else if token.SymbolType == tokentype.BracketClose {
			return array, nil
		} else {
			return nil, fmt.Errorf("found %v, expected comma or closing bracket", token.TokenType)
		}
	}
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

func isFractional(text []byte) bool {
	return slices.ContainsFunc(text, func(b byte) bool {
		return b == '.' || b == 'e' || b == 'E'
	})
}

// Package tokenizer provides a stream based lexer for json strings
package tokenizer

import (
	"cmp"
	"fmt"
	"io"
	"json-serde/pkg/tokenizer/tokentype"
	"slices"
)

type Token struct {
	TokenType  tokentype.TokenType
	Value      []byte
	SymbolType tokentype.SymbolType
}

type Tokenizer interface {
	Next() (*Token, error)
}

type tokenizer struct {
	reader       io.Reader
	buffer       []byte
	bufferLen    int
	readIndex    int
	valueIndex   int
	valuePadding int
	prevBuffer   []byte
	token        Token
}

type TokenizerConfig struct {
	BufferLen int // Size in bytes tokenizer can hold (read from reader) at a single time
}

func NewTokenizer(rd io.Reader, configs ...TokenizerConfig) Tokenizer {
	var bufferLen = 4 * 1024
	for _, config := range configs {
		if config.BufferLen > 0 {
			bufferLen = config.BufferLen
		}
	}
	return &tokenizer{reader: rd, buffer: make([]byte, bufferLen)}
}

// Gets the next token or error if present, like a lazy irreversible iterator
func (t *tokenizer) Next() (*Token, error) {
	t.clear()
	err := t.read()
	if err != nil {
		return nil, err
	}

	var stop = false
	for !stop {
		stop, err = t.scan()
		if err != nil {
			return nil, err
		}
		t.readIndex++
	}
	paddingFactor := 1
	if len(t.prevBuffer) != 0 {
		t.prevBuffer = slices.Delete(t.prevBuffer, 0, t.valuePadding)
		paddingFactor = 0
	}
	valueStartIndex := t.valueIndex + (paddingFactor * t.valuePadding)
	valueEndIndex := t.readIndex - t.valuePadding
	t.token.Value = append(t.prevBuffer, t.buffer[valueStartIndex:valueEndIndex]...)
	return &t.token, nil
}

// Read from reader if necessary (if no new data is there in buffer)
func (t *tokenizer) read() error {
	if t.isBufferEmpty() {
		if t.valueIndex != -1 {
			t.storeValue()
		}
		n, err := t.reader.Read(t.buffer)
		if err != nil {
			return err
		}
		t.bufferLen = n
		t.readIndex = 0
		if t.valueIndex != -1 {
			t.valueIndex = 0
		}
	}
	return nil
}

// Read and returns a wrapped error, should be used where more bytes must be presnt (in buffer/reader)
func (t *tokenizer) mustRead(e error) error {
	err := t.read()
	if err == io.EOF {
		return fmt.Errorf("%w, %w", ErrUnexpectedEOF, e)
	}
	if err != nil || t.isBufferEmpty() {
		return fmt.Errorf("%w, %w", cmp.Or(err, ErrUnexpectedEOF), e)
	}
	if t.valueIndex == -1 {
		t.valueIndex = 0
	}
	return nil
}

func (t *tokenizer) isBufferEmpty() bool {
	return t.bufferLen == 0 || t.readIndex >= t.bufferLen || t.valueIndex >= t.bufferLen
}

func (t *tokenizer) scan() (stop bool, err error) {
	switch t.token.TokenType {
	case tokentype.None:
		token, err := t.predictTokenType()
		if err != nil {
			return false, err
		}
		if token == tokentype.None {
			return false, nil
		}
		t.token.TokenType = token
		t.valueIndex = t.readIndex
		if token.IsPrimitive() {
			return false, nil
		} else {
			// case of symbols
			t.token.SymbolType = tokentype.SymbolFromByte(t.buffer[t.readIndex])
			return true, nil
		}
	case tokentype.String:
		return t.readString()
	case tokentype.Number:
		return t.readNumber()
	case tokentype.Boolean, tokentype.Null:
		return t.readLiteral()
	}
	return
}

func (t *tokenizer) predictTokenType() (tokentype.TokenType, error) {
	if err := t.read(); err != nil {
		return tokentype.None, err
	}
	ch := t.buffer[t.readIndex]
	switch ch {
	case ' ', '\n', '\r', '\t':
		return tokentype.None, nil

	case '"':
		t.valuePadding = 1
		return tokentype.String, nil

	case 't', 'f':
		return tokentype.Boolean, nil

	case 'n':
		return tokentype.Null, nil

	case '[', ']', '{', '}', ',', ':':
		return tokentype.Symbol, nil
	}

	if (ch >= '0' && ch <= '9') || ch == '-' {
		return tokentype.Number, nil
	}
	return tokentype.None, fmt.Errorf("%w %c", ErrInvalidToken, ch)
}

func (t *tokenizer) storeValue() {
	value := t.buffer[t.valueIndex:min(t.readIndex+1, t.bufferLen)]
	t.prevBuffer = append(t.prevBuffer, value...)
	t.valueIndex = t.readIndex + 1
}

func (t *tokenizer) clear() {
	t.token = Token{}
	t.prevBuffer = []byte{}
	t.valueIndex = -1
	t.valuePadding = 0
}

package tokenizer

import (
	"cmp"
	"fmt"
	"io"
	"json-serde/utils"
	"slices"
)

type Token struct {
	TokenType utils.TokenType
	Value     []byte
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
	BufferLen int
}

func NewTokenizer(rd io.Reader, configs ...TokenizerConfig) tokenizer {
	var bufferLen = 4 * 1024
	for _, config := range configs {
		if config.BufferLen > 0 {
			bufferLen = config.BufferLen
		}
	}
	return tokenizer{reader: rd, buffer: make([]byte, bufferLen)}
}

func (t *tokenizer) Next() (*Token, error) {
	err := t.read()
	if err != nil {
		if err == io.EOF && t.bufferLen == 0 {
			return nil, ErrUnexpectedEOF
		} else {
			return nil, err
		}
	}

	var stop = false
	for !stop {
		stop, err = t.readCh()
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
		t.valueIndex = -1
	}
	return nil
}

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

func (t *tokenizer) readCh() (stop bool, err error) {
	switch t.token.TokenType {
	case utils.None:
		token, err := t.predictTokenType()
		if err != nil {
			return false, err
		}
		if token == utils.None {
			return false, nil
		}
		t.token.TokenType = token
		isPrimitive := utils.IsPrimitiveType(token)
		if isPrimitive {
			t.valueIndex = t.readIndex
			return false, nil
		}
		return true, nil
	case utils.String:
		return t.readString()
	}
	return
}

func (t *tokenizer) predictTokenType() (utils.TokenType, error) {
	ch := t.buffer[t.readIndex]
	switch ch {
	case ' ', '\n', '\r', '\t':
		return utils.None, nil

	case '"':
		t.valuePadding = 1
		return utils.String, nil

	case 't', 'f':
		return utils.Boolean, nil

	case '[':
		return utils.Array, nil

	case '{':
		return utils.Object, nil
	}

	if ch >= '0' && ch <= '9' {
		return utils.Number, nil
	}
	return utils.None, fmt.Errorf("%w %c", ErrInvalidToken, ch)
}

func (t *tokenizer) readString() (stop bool, err error) {
	err = t.mustRead(ErrUnterminatedString)
	if err != nil {
		return false, err
	}
	ch := t.buffer[t.readIndex]
	if ch == '\\' {
		if t.readIndex+1 >= t.bufferLen {
			t.storeValue()
			err := t.mustRead(ErrInvalidEscapeChar)
			if err != nil {
				return false, err
			}
			t.prevBuffer = slices.Delete(t.prevBuffer, len(t.prevBuffer)-1, len(t.prevBuffer))
		} else {
			// removing escape symbol from buffer
			_ = slices.Delete(t.buffer, t.readIndex, t.readIndex+1)
			t.bufferLen--
		}

		ch = t.buffer[t.readIndex]
		switch ch {
		case 'b':
			t.buffer[t.readIndex] = '\b'

		case 'f':
			t.buffer[t.readIndex] = '\f'

		case 'n':
			t.buffer[t.readIndex] = '\n'

		case 'r':
			t.buffer[t.readIndex] = '\r'

		case 't':
			t.buffer[t.readIndex] = '\t'

		case '"', '\\', '/':
			// ", \, / all will be handled automatically

		default:
			return false, fmt.Errorf("%w", ErrInvalidEscapeChar)

		}
		return false, nil
	}
	if ch == '"' {
		return true, nil
	}
	return false, nil
}

func (t *tokenizer) storeValue() {
	value := t.buffer[t.valueIndex:min(t.readIndex+1, t.bufferLen)]
	t.prevBuffer = append(t.prevBuffer, value...)
	t.valueIndex = t.readIndex + 1
}

package tokenizer

import (
	"fmt"
	"io"
	"json-serde/utils"
)

type Token struct {
	TokenType utils.TokenType
}

type tokenizer struct {
	reader     io.Reader
	buffer     []byte
	bufferLen  int
	readIndex  int
	valueIndex int
	prevBuffer []byte
	token      Token
}

func NewTokenizer(rd io.Reader) tokenizer {
	return tokenizer{reader: rd, buffer: make([]byte, 4*1024)}
}

func (t *tokenizer) Next() (*Token, error) {
	err := t.read()
	if err != nil {
		return nil, err
	}

	for t.readIndex < t.bufferLen {

		t.readIndex++
	}
	return &t.token, nil
}

func (t *tokenizer) read() error {
	if t.bufferLen == 0 || t.readIndex >= t.bufferLen {
		n, err := t.reader.Read(t.buffer)
		if err != nil {
			return err
		}
		t.bufferLen = n
		t.readIndex = 0
		t.valueIndex = 0
		t.readCh()
	}
	return nil
}

func (t *tokenizer) readCh() (stop bool, err error) {
	switch t.token.TokenType {
	case utils.None:
		token, err := predictTokenType(t.buffer[t.readIndex])
		if err != nil {
			return false, err
		}
		if token == utils.None {
			return false, nil
		}
		isPrimitive := utils.IsPrimitiveType(token)
		if isPrimitive {
			t.valueIndex = t.readIndex
			return false, nil
		}
		return true, nil
	case utils.String:
	}
	return
}

func predictTokenType(ch byte) (utils.TokenType, error) {
	switch ch {
	case ' ', '\n', '\r', '\t':
		return utils.None, nil

	case '"':
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
	return utils.None, fmt.Errorf("invalid token %c", ch)
}

func (t *tokenizer) readString() (stop bool, err error) {
	ch := t.buffer[t.readIndex]
	if ch == '\\' {
		if t.readIndex+1 >= t.bufferLen {
			t.storeValue()
			err := t.read()
			if err == io.EOF {
				return false, fmt.Errorf("invalid use of escape (\\) sequence")
			} else if err != nil {
				return false, err
			}
		}
		t.buffer = append(t.buffer[:t.readIndex], t.buffer[t.readIndex+1:]...)
		ch = t.buffer[t.readIndex]
		switch ch {
		// ", \, / all will be handled automatically
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

		}
		return false, nil
	}
	return
}

func (t *tokenizer) storeValue() {
	value := t.buffer[t.valueIndex : t.readIndex+1]
	t.prevBuffer = append(t.prevBuffer, value...)
	t.valueIndex = t.readIndex + 1
}

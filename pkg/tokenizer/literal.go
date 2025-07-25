package tokenizer

import (
	"fmt"
)

var TRUE = []byte("true")
var FALSE = []byte("false")
var NULL = []byte("null")

func (t *tokenizer) readLiteral() (stop bool, err error) {
	err = t.mustRead(fmt.Errorf("%w %s", ErrInvalidToken, t.buffer[t.valueIndex:t.bufferLen]))
	if err != nil {
		return false, err
	}

	indx := len(t.prevBuffer) + t.readIndex - t.valueIndex
	ch := t.buffer[t.readIndex]
	firstLetter := t.buffer[t.valueIndex]
	if len(t.prevBuffer) > 0 {
		firstLetter = t.prevBuffer[0]
	}

	literalMap := map[byte][]byte{
		't': TRUE,
		'f': FALSE,
		'n': NULL,
	}
	literal, ok := literalMap[firstLetter]
	if !ok {
		return false, ErrInvalidToken
	}

	if indx < len(literal) {
		if ch != literal[indx] {
			return false, ErrInvalidToken
		}
		if indx < len(literal)-1 {
			return false, nil
		}
	}
	return true, nil
}

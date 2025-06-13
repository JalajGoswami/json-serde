package tokenizer

import "fmt"

func (t *tokenizer) readBoolean() (stop bool, err error) {
	if err = t.read(); err != nil {
		return false, err
	}
	if t.isBufferEmpty() {
		return true, fmt.Errorf("%w %s", ErrInvalidToken, t.buffer[t.valueIndex:t.bufferLen])
	}
	indx := len(t.prevBuffer) + t.readIndex - t.valueIndex
	ch := t.buffer[t.readIndex]
	firstLetter := t.buffer[t.valueIndex]
	if len(t.prevBuffer) > 0 {
		firstLetter = t.prevBuffer[0]
	}
	if firstLetter == 't' {
		TRUE := []byte("true")
	}
}

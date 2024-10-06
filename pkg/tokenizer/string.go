package tokenizer

import (
	"fmt"
	"slices"
)

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

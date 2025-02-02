package tokenizer

import "io"

func (t *tokenizer) readNumber() (stop bool, err error) {
	var firstDigit byte // if not nil signifies readIndex is at second position
	if len(t.prevBuffer) == 0 && t.valueIndex == t.readIndex-1 {
		// executes only for the first time
		firstDigit = t.buffer[t.valueIndex]
		if firstDigit == '-' {
			err = t.mustRead(ErrInvalidEndOfNumber)
			if err != nil {
				return false, err
			}
		}
	}
	if err = t.read(); err != nil && err != io.EOF {
		return false, err
	}
	if t.isBufferEmpty() {
		return true, nil
	}
	ch := t.buffer[t.readIndex]

	// handling numbers like 05 or -05, by ending token after 0,
	// thus treating them as two number tokens, which will be handled by parser
	if firstDigit == '0' {
		if ch != '.' && ch != 'e' && ch != 'E' {
			return true, nil
		}
	} else if firstDigit == '-' && ch == '0' {
		nextCh := t.buffer[t.readIndex+1]
		if nextCh != '.' && nextCh != 'e' && nextCh != 'E' {
			return true, nil
		}
	}

	switch {
	case isDigit(ch):
		{
			if firstDigit == '0' {
				return false, ErrInvalidNumber
			}
			return false, nil
		}

	case ch == '.':
		{
			if firstDigit == '-' {
				return false, ErrInvalidNumber
			}
			t.readIndex++
			err := t.mustRead(ErrInvalidEndOfNumber)
			if err != nil {
				return false, err
			}
			ch = t.buffer[t.readIndex]
			if !isDigit(ch) {
				return false, ErrInvalidEndOfNumber
			}

		}

	case ch == 'e' || ch == 'E':
		{
			var lastDigit byte
			if t.readIndex > t.valueIndex {
				lastDigit = t.buffer[t.readIndex-1]
			} else {
				lastDigit = t.prevBuffer[len(t.prevBuffer)-1]
			}
			if lastDigit == '-' || lastDigit == '.' {
				return false, ErrInvalidNumber
			}

			t.readIndex++
			err := t.mustRead(ErrInvalidEndOfNumber)
			if err != nil {
				return false, err
			}
			ch = t.buffer[t.readIndex]
			if ch == '-' || ch == '+' {
				t.readIndex++
				err := t.mustRead(ErrInvalidEndOfNumber)
				if err != nil {
					return false, err
				}
				ch = t.buffer[t.readIndex]
				if !isDigit(ch) {
					return false, ErrInvalidEndOfNumber
				}
			} else if !isDigit(ch) {
				return false, ErrInvalidEndOfNumber
			}
		}

	default:
		{
			t.readIndex--
			return true, nil
		}
	}
	return
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

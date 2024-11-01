package tokenizer

func (t *tokenizer) readNumber() (stop bool, err error) {
	var firstDigit byte // if not nil signifies readIndex is at second position
	if len(t.prevBuffer) == 0 && t.valueIndex == t.readIndex-1 {
		// executes only for the first time
		firstDigit = t.buffer[t.valueIndex]
		if firstDigit == '-' {
			t.mustRead(ErrInvalidEndOfNumber)
		}
	}
	ch := t.buffer[t.readIndex]

	switch {
	case ch >= '0' && ch <= '9':
		{
			// handling invalid number like 05
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
			if ch < '0' || ch > '9' {
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
		}

	default:
		{
			t.readIndex--
			return true, nil
		}
	}
	return
}

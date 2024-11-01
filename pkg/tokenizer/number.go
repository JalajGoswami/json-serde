package tokenizer

func (t *tokenizer) readNumber() (stop bool, err error) {
	t.mustRead(ErrInvalidEndOfNumber)
	ch := t.buffer[t.readIndex]

	switch {
	case ch >= '0' && ch <= '9':
		{
			if (len(t.prevBuffer) == 0 && t.valueIndex == t.readIndex-1) ||
				(len(t.prevBuffer) == 1 && t.valueIndex == t.readIndex) {
				// handling invalid number like 05
				var lastDigit byte
				if len(t.prevBuffer) > 0 {
					lastDigit = t.prevBuffer[0]
				} else {
					lastDigit = t.buffer[t.valueIndex]
				}
				if lastDigit == '0' {
					return false, ErrInvalidNumber
				}
			}
			return false, nil
		}

	case ch == '.':
		{
			if t.readIndex+1 >= t.bufferLen {
				t.storeValue()
				err := t.mustRead(ErrInvalidEndOfNumber)
				if err != nil {
					return false, err
				}
			}
			t.readIndex++
			ch = t.buffer[t.readIndex]
			if ch < '0' || ch > '9' {
				return false, ErrInvalidEndOfNumber
			}

		}

	}
	return
}

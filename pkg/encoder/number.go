package encoder

import (
	"math"
	"slices"
	"strconv"
)

func (e *Encoder) encodeInt(value int64) error {
	if value < 0 {
		err := e.write([]byte{'-'}) // negative sign
		if err != nil {
			return err
		}
		value *= -1 // faster math.Abs
	}
	return e.write(encodeInteger(value))
}

func (e *Encoder) encodeUint(value uint64) error {
	return e.write(encodeInteger(value))
}

func (e *Encoder) encodeFloat(value float64) error {
	return e.write([]byte(strconv.FormatFloat(value, 'g', -1, 64)))
}

type Integer interface {
	~int64 | ~uint64
}

func encodeInteger[T Integer](value T) []byte {
	if value == 0 {
		return []byte{'0'}
	}
	estimatedDigits := int(math.Log10(float64(value))) + 1
	var digits = make([]byte, 0, estimatedDigits)
	for value > 0 {
		lsd := value % 10
		digits = append(digits, '0'+byte(lsd))
		value = value / 10
	}
	var encoded = make([]byte, 0, len(digits))
	for _, d := range slices.Backward(digits) {
		encoded = append(encoded, d)
	}
	return encoded
}

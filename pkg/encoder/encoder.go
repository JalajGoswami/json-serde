package encoder

import (
	"reflect"
	"slices"
)

type Encoder struct {
	data []byte
}

func (e *Encoder) Encode(v any) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	return e.encodeValue(v)
}

func (e *Encoder) encodeValue(v any) ([]byte, error) {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Bool:
		b := v.(bool)
		if b {
			e.data = append(e.data, 116, 114, 117, 101) // ascii of true
		} else {
			e.data = append(e.data, 102, 97, 108, 115, 101) // ascii of false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.encodeInt(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		e.encodeUint(value.Uint())
	}
	return e.data, nil
}
func (e *Encoder) encodeInt(value int64) {
	if value < 0 {
		e.data = append(e.data, '-') // negative sign
		value *= -1                  // faster math.Abs
	}
	e.data = append(e.data, encodeInteger(value)...)
}

func (e *Encoder) encodeUint(value uint64) {
	e.data = append(e.data, encodeInteger(value)...)
}

type Integer interface {
	~int64 | ~uint64
}

func encodeInteger[T Integer](value T) []byte {
	var encoded []byte
	if value == 0 {
		encoded = append(encoded, '0')
	}
	var digits []byte
	for value > 0 {
		lsd := value % 10
		digits = append(digits, byte(lsd+48))
		value = value / 10
	}
	encoded = slices.Grow(encoded, len(digits))
	for _, d := range slices.Backward(digits) {
		encoded = append(encoded, d)
	}
	return encoded
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

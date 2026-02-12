// Package encoder serialize native data types into their json counterparts
package encoder

import (
	"math"
	"reflect"
	"slices"
	"strconv"
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
	case reflect.Float32, reflect.Float64:
		e.encodeFloat(value.Float())
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

func (e *Encoder) encodeFloat(value float64) {
	e.data = append(e.data, strconv.FormatFloat(value, 'g', -1, 64)...)
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

func NewEncoder() *Encoder {
	return &Encoder{}
}

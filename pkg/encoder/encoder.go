// Package encoder serialize native data types into their json counterparts
package encoder

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"reflect"
)

var ErrUnserializableType = errors.New("unserializable type")

type Encoder struct {
	writer io.Writer
	n      int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{writer: w}
}

func (e *Encoder) Encode(v any) (int, error) {
	err := e.encodeValue(v)
	return e.n, err
}

func (e *Encoder) write(p []byte) error {
	n, err := e.writer.Write(p)
	e.n += n
	if err != nil || n < len(p) {
		return cmp.Or(err, io.EOF)
	}
	return nil
}

func (e *Encoder) encodeValue(v any) error {
	if v == nil {
		return e.write([]byte("null"))
	}
	value := reflect.ValueOf(v)
	var err error
	switch value.Kind() {
	case reflect.Bool:
		err = e.encodeBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = e.encodeInt(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = e.encodeUint(value.Uint())
	case reflect.Float32, reflect.Float64:
		err = e.encodeFloat(value.Float())
	case reflect.String:
		err = e.encodeString(v)
	case reflect.Array, reflect.Slice:
		err = e.encodeArray(value)

	default:
		err = fmt.Errorf("%w, type=%s", ErrUnserializableType, value.Type())
	}
	return err
}

func (e *Encoder) encodeBool(v any) error {
	if v.(bool) {
		return e.write([]byte("true"))
	}
	return e.write([]byte("false"))
}

func (e *Encoder) encodeString(v any) error {
	n, err := fmt.Fprintf(e.writer, "%q", v)
	e.n += n
	return err
}

func (e *Encoder) encodeArray(value reflect.Value) error {
	err := e.write([]byte{'['})
	if err != nil {
		return err
	}

	n := value.Len()
	for i := range n {
		elem := value.Index(i)
		err := e.encodeValue(elem.Interface())
		if err != nil {
			return err
		}
		if i < n-1 {
			err := e.write([]byte{','})
			if err != nil {
				return err
			}
		}
	}
	return e.write([]byte{']'})
}

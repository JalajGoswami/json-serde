package parser

import "errors"

var ErrNilInput = errors.New("invalid input <nil>")

var ErrNonPointer = errors.New("non pointer input")

var ErrUnexpectedToken = errors.New("unexpected token")

var ErrTrailingComma = errors.New("trailing comma")

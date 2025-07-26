package parser

import "errors"

var ErrNilInput = errors.New("invalid input <nil>")

var ErrNonPointer = errors.New("non pointer input")

package tokenizer

import (
	"errors"
	"io"
)

var ErrInvalidToken = errors.New("invalid token")

var ErrUnexpectedEOF = io.ErrUnexpectedEOF

var ErrInvalidEscapeChar = errors.New("invalid use of escape (\\) sequence")

var ErrUnterminatedString = errors.New("unterminated string")

var ErrInvalidNumber = errors.New("invalid number")

var ErrInvalidEndOfNumber = errors.New("invalid end of number")

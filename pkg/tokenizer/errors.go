package tokenizer

import (
	"errors"
	"io"
)

var ErrInvalidToken = errors.New("invalid token")

var ErrUnexpectedEOF = io.ErrUnexpectedEOF

var ErrInvalidEscapeChar = errors.New("invalid use of escape (\\) sequence")

package tokenizer

import (
	"bytes"
	"errors"
	"io"
	"json-serde/utils"
	"testing"
)

func TestString(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		err    error
		result string
	}{
		{name: "Empty Contents", input: "", err: ErrUnexpectedEOF},
		{name: "Unterminated Quote", input: "\"", err: ErrUnexpectedEOF},
		{name: "Un-quotted String", input: "hi", err: ErrInvalidToken},
		{name: "Unterminated String", input: "\"hi", err: ErrUnexpectedEOF},
		{name: "Small String", input: "\"hi\"", result: "hi"},
		{
			name:   "Long String",
			input:  "\"a long string overflowing provided buffer length\"",
			result: "a long string overflowing provided buffer length",
		},
		{name: "Multiline", input: "\"a\\n multilinee\\b\\t msg\"", result: "a\n multiline\t msg"},
		{name: "Esape Keywords", input: "\"\\\" \\\\ \\/\"", result: "\\\" \\\\ \\/"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rd := readerFrom(testCase.input)
			config := TokenizerConfig{
				BufferLen: 8, // small value to check bytes synchronisation
			}
			tokenGenerator := NewTokenizer(rd, config)
			token, err := tokenGenerator.Next()
			if testCase.err != nil && !errors.Is(err, testCase.err) {
				t.Errorf("Expected error: %v got %v", testCase.err, err)
			}
			if testCase.err == nil {
				if token == nil {
					t.Log(testCase.input)
					t.Errorf("Got error: %v with token: <nil> expected a string token", err)
					return
				}
				if token.TokenType != utils.String {
					t.Errorf("Invalid token type expected: %v got %v", utils.String, token.TokenType)
					return
				}
			}
		})
	}
}

func readerFrom(value string) io.Reader {
	return bytes.NewReader([]byte(value))
}
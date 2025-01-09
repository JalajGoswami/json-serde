package tokenizer

import (
	"bytes"
	"errors"
	"io"
	"json-serde/pkg/tokenizer/tokentype"
	"testing"
)

type testCase struct {
	name   string
	input  string
	err    error
	result Token
}

func TestString(t *testing.T) {
	testCases := []testCase{
		{name: "Empty Contents", input: "", err: io.EOF},
		{name: "Unterminated Quote", input: "\"", err: ErrUnexpectedEOF},
		{name: "Un-quoted String", input: "hi", err: ErrInvalidToken},
		{name: "Unterminated String", input: "\"hi", err: ErrUnexpectedEOF},
		{name: "Small String", input: "\"hi\"", result: stringToken("hi")},
		{
			name:   "Long String",
			input:  "\"a long string overflowing provided buffer length\"",
			result: stringToken("a long string overflowing provided buffer length"),
		},
		{
			name:   "Multiline",
			input:  "\"a\\n multiline\\b\\t msg\"",
			result: stringToken("\n multiline\b\t msg"),
		},
		{name: "Escape Keywords", input: "\"\\\" \\\\ \\/\"", result: stringToken("\" \\ /")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTestCase(t, tc)
		})
	}
}

func runTestCase(t *testing.T, tc testCase) {
	rd := bytes.NewReader([]byte(tc.result.Value))
	config := TokenizerConfig{
		BufferLen: 8, // small value to check bytes synchronization
	}
	tokenGenerator := NewTokenizer(rd, config)
	token, err := tokenGenerator.Next()
	if tc.err != nil && !errors.Is(err, tc.err) {
		t.Errorf("Expected error: %v got %v", tc.err, err)
		return
	}
	if tc.err == nil {
		if token == nil {
			t.Errorf("Got error: %v with token: <nil> expected a %v token", err, tc.result.TokenType)
			return
		}
		if token.TokenType != tc.result.TokenType {
			t.Errorf(
				"Invalid token type expected: %v got %v", tc.result.TokenType, token.TokenType,
			)
			return
		}
		switch tc.result.TokenType {
		case tokentype.String:
			if string(token.Value) != string(tc.result.Value) {
				t.Errorf("Expected value: %v, got %v", tc.result, string(token.Value))
				return
			}
		case tokentype.Number:

		}
	}
}

func stringToken(s string) Token {
	return Token{TokenType: tokentype.String, Value: []byte(s)}
}

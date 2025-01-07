package evaluation

import (
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	testCases := []struct {
		text           string
		expectedTokens []Token
	}{
		{
			text: "  foo  ",
			expectedTokens: []Token{
				{tokenType: TokenVariable, literal: "foo"},
			},
		},
		{
			text: "not(1)  ",
			expectedTokens: []Token{
				{tokenType: TokenNot, literal: "not"},
				{tokenType: TokenLparan, literal: "("},
				{tokenType: TokenValue, literal: "1"},
				{tokenType: TokenRparan, literal: ")"},
			},
		},
		{
			text: " and(not(1), not(X))  ",
			expectedTokens: []Token{
				{tokenType: TokenAnd, literal: "and"},
				{tokenType: TokenLparan, literal: "("},
				{tokenType: TokenNot, literal: "not"},
				{tokenType: TokenLparan, literal: "("},
				{tokenType: TokenValue, literal: "1"},
				{tokenType: TokenRparan, literal: ")"},
				{tokenType: TokenComma, literal: ","},
				{tokenType: TokenNot, literal: "not"},
				{tokenType: TokenLparan, literal: "("},
				{tokenType: TokenVariable, literal: "X"},
				{tokenType: TokenRparan, literal: ")"},
				{tokenType: TokenRparan, literal: ")"},
			},
		},
		{
			text: "mux(or(a,b), xor(0,1), and(c,d))",
			expectedTokens: []Token{
				{tokenType: TokenMux, literal: "mux"},
				{tokenType: TokenLparan, literal: "("},
				{tokenType: TokenOr, literal: "or"},
				{tokenType: TokenLparan, literal: "("},
				{tokenType: TokenVariable, literal: "a"},
				{tokenType: TokenComma, literal: ","},
				{tokenType: TokenVariable, literal: "b"},
				{tokenType: TokenRparan, literal: ")"},
				{tokenType: TokenComma, literal: ","},
				{tokenType: TokenXor, literal: "xor"},
				{tokenType: TokenLparan, literal: "("},
				{tokenType: TokenValue, literal: "0"},
				{tokenType: TokenComma, literal: ","},
				{tokenType: TokenValue, literal: "1"},
				{tokenType: TokenRparan, literal: ")"},
				{tokenType: TokenComma, literal: ","},
				{tokenType: TokenAnd, literal: "and"},
				{tokenType: TokenLparan, literal: "("},
				{tokenType: TokenVariable, literal: "c"},
				{tokenType: TokenComma, literal: ","},
				{tokenType: TokenVariable, literal: "d"},
				{tokenType: TokenRparan, literal: ")"},
				{tokenType: TokenRparan, literal: ")"},
			},
		},
		{
			text: " nand(foo,1 ", // doesn't have to be a correct expression
			expectedTokens: []Token{
				{tokenType: TokenNand, literal: "nand"},
				{tokenType: TokenLparan, literal: "("},
				{tokenType: TokenVariable, literal: "foo"},
				{tokenType: TokenComma, literal: ","},
				{tokenType: TokenValue, literal: "1"},
			},
		},
		{
			text: "NOT(1)", // case sensitive keywords
			expectedTokens: []Token{
				{tokenType: TokenVariable, literal: "NOT"},
				{tokenType: TokenLparan, literal: "("},
				{tokenType: TokenValue, literal: "1"},
				{tokenType: TokenRparan, literal: ")"},
			},
		},
	}

	for i, tc := range testCases {
		actual, err := ParseTokens(tc.text)
		if err != nil {
			t.Errorf("Got an unexpected error in test case %d: Error: %v", i, err)
		}
		if !reflect.DeepEqual(actual, tc.expectedTokens) {
			t.Errorf("Error in test case %d: Expected %v, but got %v", i, tc.expectedTokens, actual)
		}
	}
}

func TestTokenizerErrors(t *testing.T) {
	testCases := []string{
		" dmux(foo_bar,1,XOR) ", // underscore not allowed in variable names
		"2and3",                 // invalid value + keyword as part of identifier
		"@#$",                   // multiple invalid characters
	}

	for _, tc := range testCases {
		tok, err := ParseTokens(tc)
		if err == nil {
			t.Errorf("Expected to get an error for string \"%s\", but got %v", tc, tok)
		}
	}
}

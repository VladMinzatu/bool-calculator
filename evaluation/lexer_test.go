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
			text: " dmux(foo_bar,1,XOR) ",
			expectedTokens: []Token{
				{tokenType: TokenDmux, literal: "dmux"},
				{tokenType: TokenLparan, literal: "("},
				{tokenType: TokenVariable, literal: "foo"},
				{tokenType: TokenInvalid, literal: "_"}, // underscores in variable names not allowed
				{tokenType: TokenVariable, literal: "bar"},
				{tokenType: TokenComma, literal: ","},
				{tokenType: TokenValue, literal: "1"},
				{tokenType: TokenComma, literal: ","},
				{tokenType: TokenVariable, literal: "XOR"}, // uppercase ok for variables, but gate names are case sensitive
				{tokenType: TokenRparan, literal: ")"},
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
		{
			text: "2and3", // invalid value + keyword as part of identifier
			expectedTokens: []Token{
				{tokenType: TokenInvalid, literal: "2"},
				{tokenType: TokenAnd, literal: "and"},
				{tokenType: TokenInvalid, literal: "3"},
			},
		},
		{
			text: "@#$", // multiple invalid characters
			expectedTokens: []Token{
				{tokenType: TokenInvalid, literal: "@"},
				{tokenType: TokenInvalid, literal: "#"},
				{tokenType: TokenInvalid, literal: "$"},
			},
		},
	}

	for i, tc := range testCases {
		actual := ParseTokens(tc.text)
		if !reflect.DeepEqual(actual, tc.expectedTokens) {
			t.Errorf("Error in test case %d: Expected %v, but got %v", i, tc.expectedTokens, actual)
		}
	}
}

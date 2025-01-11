package evaluation

import (
	"reflect"
	"testing"
)

func TestParsingAndEvaluation(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectsErr     bool
		expectedResult []bool
	}{
		{
			name:           "simple literal 1",
			input:          "1",
			expectsErr:     false,
			expectedResult: []bool{true},
		},
		{
			name:           "simple literal 0",
			input:          "0",
			expectsErr:     false,
			expectedResult: []bool{false},
		},
		{
			name:           "not gate",
			input:          "not(1)",
			expectsErr:     false,
			expectedResult: []bool{false},
		},
		{
			name:           "and gate",
			input:          "and(1,1)",
			expectsErr:     false,
			expectedResult: []bool{true},
		},
		{
			name:           "and gate2",
			input:          "and(1,0)",
			expectsErr:     false,
			expectedResult: []bool{false},
		},
		{
			name:           "or gate",
			input:          "or(1,0)",
			expectsErr:     false,
			expectedResult: []bool{true},
		},
		{
			name:           "xor gate",
			input:          "xor(1,1)",
			expectsErr:     false,
			expectedResult: []bool{false},
		},
		{
			name:           "nand gate",
			input:          "nand(1,1)",
			expectsErr:     false,
			expectedResult: []bool{false},
		},
		{
			name:           "mux gate",
			input:          "mux(1,0,0)",
			expectsErr:     false,
			expectedResult: []bool{false},
		},
		{
			name:           "dmux gate",
			input:          "dmux(1,0)",
			expectsErr:     false,
			expectedResult: []bool{true, false},
		},
		{
			name:           "nested expression",
			input:          "and(not(0),or(1,0))",
			expectsErr:     false,
			expectedResult: []bool{true},
		},
		{
			name:       "invalid gate name",
			input:      "invalid(1,1)",
			expectsErr: true,
		},
		{
			name:       "missing arguments",
			input:      "and(1)",
			expectsErr: true,
		},
		{
			name:       "invalid token",
			input:      "and(1,10)",
			expectsErr: true,
		},
		{
			name:       "too many arguments",
			input:      "and(1,1,1)",
			expectsErr: true,
		},
		{
			name:       "missing closing parenthesis",
			input:      "and(1,1",
			expectsErr: true,
		},
		{
			name:       "empty expression",
			input:      "",
			expectsErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := ParseExpression(tc.input)
			if tc.expectsErr {
				if err == nil {
					t.Errorf("ParseExpression() expected error but didn't get any for input %v", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseExpression() encountered unexpected error for input %v: err = %v", tc.input, err)
				return
			}

			got := expr.Evaluate()
			if !reflect.DeepEqual(got, tc.expectedResult) {
				t.Errorf("Expression.Evaluate() = %v, expected %v", got, tc.expectedResult)
			}
		})
	}
}

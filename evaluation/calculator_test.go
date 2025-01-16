package evaluation

import (
	"reflect"
	"testing"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   *Result
		expectErr  bool
	}{
		{
			name:       "Simple constant true",
			expression: "1",
			expected: &Result{
				Variables:   nil,
				Outputs:     [][]bool{{true}},
				Assignments: nil,
			},
		},
		{
			name:       "Simple constant false",
			expression: "0",
			expected: &Result{
				Variables:   nil,
				Outputs:     [][]bool{{false}},
				Assignments: nil,
			},
		},
		{
			name:       "Single variable",
			expression: "X",
			expected: &Result{
				Variables:   []string{"X"},
				Outputs:     [][]bool{{false}, {true}},
				Assignments: [][]bool{{false}, {true}},
			},
		},
		{
			name:       "AND operation",
			expression: "and(X,Y)",
			expected: &Result{
				Variables: []string{"X", "Y"},
				Outputs:   [][]bool{{false}, {false}, {false}, {true}},
				Assignments: [][]bool{
					{false, false},
					{false, true},
					{true, false},
					{true, true},
				},
			},
		},
		{
			name:       "OR operation",
			expression: "or(X,Y)",
			expected: &Result{
				Variables: []string{"X", "Y"},
				Outputs:   [][]bool{{false}, {true}, {true}, {true}},
				Assignments: [][]bool{
					{false, false},
					{false, true},
					{true, false},
					{true, true},
				},
			},
		},
		{
			name:       "NOT operation",
			expression: "not(X)",
			expected: &Result{
				Variables:   []string{"X"},
				Outputs:     [][]bool{{true}, {false}},
				Assignments: [][]bool{{false}, {true}},
			},
		},
		{
			name:       "Invalid expression",
			expression: "and(X",
			expectErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Compute(tc.expression)
			if (err != nil) != tc.expectErr {
				t.Errorf("Compute() returned error = %v, when expectErr=%v", err, tc.expectErr)
				return
			}
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error, but didn't get any")
				}
				return
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected  %v, but Compute() returned %v, ", tc.expected, got)
			}
		})
	}
}

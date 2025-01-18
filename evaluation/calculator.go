package evaluation

import (
	"sort"
	"strings"
)

type Result struct {
	Variables   []string
	Outputs     [][]bool
	Assignments [][]bool
}

func (r Result) String() string {
	outputSpacing := "  "

	var sb strings.Builder

	if len(r.Variables) == 0 {
		// just print the result
		for i, val := range r.Outputs[0] {
			if i > 0 {
				sb.WriteString(outputSpacing)
			}
			sb.WriteString(boolToString(val))
		}
		sb.WriteString("\n")
		return sb.String()
	}

	// Print header
	for _, v := range r.Variables {
		sb.WriteString(v)
		sb.WriteString("\t")
	}
	sb.WriteString("Output\n")

	// Print assignments
	for i := 0; i < len(r.Assignments); i++ {
		for _, val := range r.Assignments[i] {
			sb.WriteString(boolToString(val))
			sb.WriteString("\t")
		}

		for idx, val := range r.Outputs[i] {
			if idx > 0 {
				sb.WriteString(outputSpacing)
			}
			sb.WriteString(boolToString(val))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func boolToString(val bool) string {
	if val {
		return "1"
	}
	return "0"
}

func Compute(expression string) (*Result, error) {
	expr, vars, err := ParseExpression(expression)
	if err != nil {
		return nil, err
	}

	variables := getVarsSlice(vars)
	if len(variables) == 0 {
		value, err := expr.Evaluate(map[string]bool{})
		if err != nil {
			return nil, err
		}
		return &Result{Variables: nil, Outputs: [][]bool{value}, Assignments: nil}, nil
	}

	result := Result{Variables: variables}
	assignments := generateCombinations(len(variables))
	for _, assignment := range assignments {
		value, err := expr.Evaluate(getArgs(variables, assignment))
		if err != nil {
			return nil, err
		}
		result.Outputs = append(result.Outputs, value)
		result.Assignments = append(result.Assignments, assignment)
	}
	return &result, nil
}

func generateCombinations(n int) [][]bool {
	total := 1 << n //2^n combinations
	result := make([][]bool, total)

	for i := 0; i < total; i++ {
		combination := make([]bool, n)
		for j := 0; j < n; j++ {
			combination[n-j-1] = (i & (1 << j)) != 0
		}
		result[i] = combination
	}

	return result
}

func getVarsSlice(vars map[string]bool) []string {
	result := []string{}
	for v, _ := range vars {
		result = append(result, v)
	}
	sort.Strings(result)
	return result
}

func getArgs(variables []string, assignment []bool) map[string]bool {
	result := map[string]bool{}
	for i := 0; i < len(variables); i++ {
		result[variables[i]] = assignment[i]
	}
	return result
}

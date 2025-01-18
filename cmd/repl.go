package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/VladMinzatu/bool-calculator/evaluation"
)

const (
	exitStr       = "exit"
	outputSpacing = "  "
)

func RunRepl() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Boolean Calculator REPL")
	fmt.Printf("Enter expressions to evaluate (or '%s' to quit)\n", exitStr)

	for {
		fmt.Print(">>> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == exitStr {
			break
		}

		if input == "" {
			continue
		}

		result, err := evaluation.Compute(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Print(result.String())

		// 	if len(result.Variables) == 0 {
		// 		// just print the result
		// 		for i, val := range result.Outputs[0] {
		// 			if i > 0 {
		// 				fmt.Print(outputSpacing)
		// 			}
		// 			fmt.Print(boolValueStr(val))
		// 		}
		// 		fmt.Println()
		// 		continue
		// 	}

		// 	// We have variables, so we'll print all possible assignments:
		// 	// Print header
		// 	for _, v := range result.Variables {
		// 		fmt.Printf("%s\t", v)
		// 	}
		// 	fmt.Printf("Output\n")

		// 	// Print assignments
		// 	for i := 0; i < len(result.Assignments); i++ {
		// 		for _, val := range result.Assignments[i] {
		// 			fmt.Printf("%s\t", boolValueStr(val))
		// 		}

		// 		for idx, val := range result.Outputs[i] {
		// 			if idx > 0 {
		// 				fmt.Print(outputSpacing)
		// 			}
		// 			fmt.Printf("%s", boolValueStr(val))
		// 		}
		// 		fmt.Println()
		// 	}
	}
}

func boolValueStr(val bool) string {
	if val {
		return "1"
	} else {
		return "0"
	}
}

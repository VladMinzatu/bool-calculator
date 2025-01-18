package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/VladMinzatu/bool-calculator/evaluation"
)

const (
	exitStr = "exit"
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
	}
}

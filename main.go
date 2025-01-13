package main

import (
	"fmt"
	"os"

	"github.com/VladMinzatu/bool-calculator/evaluation"
)

func main() {
	expression, _, err := evaluation.ParseExpression("dmux(and(1, not(0)), 0)")
	if err != nil {
		fmt.Println(err)
		return
	}
	val, err := expression.Evaluate(map[string]bool{})
	if err != nil {
		fmt.Printf("error occurred evaluating expression: %v", err)
		os.Exit(1)
	}
	fmt.Println(val)
}

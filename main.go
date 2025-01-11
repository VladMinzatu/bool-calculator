package main

import (
	"fmt"

	"github.com/VladMinzatu/bool-calculator/evaluation"
)

func main() {
	expression, _, err := evaluation.ParseExpression("dmux(and(1, not(0)), 0)")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(expression.Evaluate())
}

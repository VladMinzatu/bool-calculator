package main

import (
	"fmt"

	"github.com/VladMinzatu/bool-calculator/evaluation"
)

func main() {
	fmt.Println(evaluation.Compute("mux(and(1, not(X)), Y, 0)"))
}

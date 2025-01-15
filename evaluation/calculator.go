package evaluation

func GenerateCombinations(n int) [][]bool {
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

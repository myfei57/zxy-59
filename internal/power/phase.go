package power

func BalancePhases(requests []int) [3]int {
	var phases [3]int
	for i, request := range requests {
		phases[i%3] += request
	}
	return phases
}

func PhaseImbalance(a, b, c int) int {
	max := a
	if b > max {
		max = b
	}
	if c > max {
		max = c
	}
	min := a
	if b < min {
		min = b
	}
	if c < min {
		min = c
	}
	return max - min
}

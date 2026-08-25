package power

func TotalRequested(requests []int) int {
	total := 0
	for _, request := range requests {
		total += request
	}
	return total
}

func AverageRequested(requests []int) int {
	if len(requests) == 0 {
		return 0
	}
	return TotalRequested(requests) / len(requests)
}

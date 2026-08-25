package power

type ShedResult struct {
	Allowed []int `json:"allowed"`
	Shed    []int `json:"shed"`
}

func (s *Service) Shed(requests []int) ShedResult {
	allowed := make([]int, len(requests))
	shed := make([]int, len(requests))
	remaining := s.grid.Current()
	for i, request := range requests {
		if request <= remaining {
			allowed[i] = request
			remaining -= request
			continue
		}
		allowed[i] = remaining
		shed[i] = request - remaining
		remaining = 0
		for j := i + 1; j < len(requests); j++ {
			shed[j] = requests[j]
		}
		break
	}
	return ShedResult{Allowed: allowed, Shed: shed}
}

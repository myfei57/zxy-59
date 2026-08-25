package power

func (s *Service) Coordinate(requests []int) []int {
	remaining := s.grid.Current()
	out := make([]int, len(requests))
	for i, request := range requests {
		if request <= 0 {
			continue
		}
		if request > remaining {
			out[i] = remaining
			remaining = 0
		} else {
			out[i] = request
			remaining -= request
		}
		if remaining <= 0 {
			break
		}
	}
	return out
}

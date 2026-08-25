package power

func (s *Service) Utilization() float64 {
	capacity := float64(s.grid.Current())
	if capacity <= 0 {
		return 0
	}
	return s.quota.Total() / capacity
}

func (s *Service) PeakLoad(requests []int) int {
	peak := 0
	for _, request := range requests {
		if request > peak {
			peak = request
		}
	}
	return peak
}

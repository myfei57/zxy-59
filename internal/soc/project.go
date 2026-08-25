package soc

func (s *Service) ProjectSoc(busID string, additionalWattHours int) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	capacity := s.capacity[busID]
	if capacity <= 0 {
		return 0
	}
	projected := s.energy[busID] + additionalWattHours
	if projected > capacity {
		projected = capacity
	}
	percent := projected * 100 / capacity
	if percent > 100 {
		percent = 100
	}
	return percent
}

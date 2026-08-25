package bus

func (s *Service) Occupancy(pileNumber string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, mapped := range s.mapping {
		if mapped == pileNumber {
			count++
		}
	}
	return count
}

func (s *Service) MappedBusCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.mapping)
}

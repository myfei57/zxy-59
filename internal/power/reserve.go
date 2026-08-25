package power

func (s *Service) Reserve(pileID string, requested int) int {
	s.reserveMu.Lock()
	defer s.reserveMu.Unlock()
	capacity := s.grid.Current()
	used := 0
	for id, watts := range s.reserved {
		if id != pileID {
			used += watts
		}
	}
	available := capacity - used
	if requested > available {
		requested = available
	}
	if requested < 0 {
		requested = 0
	}
	s.reserved[pileID] = requested
	return requested
}

func (s *Service) Release(pileID string) {
	s.reserveMu.Lock()
	defer s.reserveMu.Unlock()
	delete(s.reserved, pileID)
}

func (s *Service) ReservedTotal() int {
	s.reserveMu.RLock()
	defer s.reserveMu.RUnlock()
	total := 0
	for _, watts := range s.reserved {
		total += watts
	}
	return total
}

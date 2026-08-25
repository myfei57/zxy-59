package bus

func (s *Service) SetMapping(busID, pileNumber string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mapping[busID] = pileNumber
}

func (s *Service) Mapping(busID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.mapping[busID]
}

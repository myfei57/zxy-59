package soc

func (s *Service) Refresh(busID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.estimated[busID] = s.current[busID]
	return s.estimated[busID]
}

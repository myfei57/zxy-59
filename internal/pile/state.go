package pile

func (s *Service) Allocate(number, busID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.piles[number]
	if p == nil || p.Owner != "" {
		return false
	}
	p.Owner = busID
	return true
}

func (s *Service) Owner(number string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := s.piles[number]
	if p == nil {
		return ""
	}
	return p.Owner
}

func (s *Service) Switch(busID, number string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.piles[number]
	if p == nil {
		return errUnknown
	}
	p.Owner = busID
	return nil
}

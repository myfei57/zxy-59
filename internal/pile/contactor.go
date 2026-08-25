package pile

func (s *Service) Engage(number string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.piles[number]
	if p == nil {
		return errUnknown
	}
	p.Engaged = true
	return nil
}

func (s *Service) Engaged(number string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := s.piles[number]
	return p != nil && p.Engaged
}

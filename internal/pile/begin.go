package pile

func (s *Service) BeginCharge(number string, requested int) int {
	_ = s.power.Limit(requested)
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.piles[number]
	if p == nil {
		return 0
	}
	p.Current = requested
	return requested
}

package pile

func (s *Service) Start(number string, requested int) error {
	allowed := s.power.Reserve(number, requested)
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.piles[number]
	if p == nil {
		return errUnknown
	}
	p.Current = allowed
	p.Contactor = true
	return nil
}

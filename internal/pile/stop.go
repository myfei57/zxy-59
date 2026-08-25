package pile

func (s *Service) Stop(number string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.piles[number]
	if p == nil {
		return errUnknown
	}
	p.Current = 0
	p.Contactor = false
	p.Engaged = false
	s.power.Release(number)
	return nil
}

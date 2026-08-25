package pile

func (s *Service) Stop(number string) error {
	// Release this pile's share of the shared power reserve so subsequent
	// vehicles can reserve the full available capacity instead of inheriting
	// the stale reservation left behind by the previous charge.
	s.power.Release(number)
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.piles[number]
	if p == nil {
		return errUnknown
	}
	p.Current = 0
	p.Contactor = false
	p.Engaged = false
	return nil
}

package pile

import "buscharge/internal/soc"

func (s *Service) TrickleStop(number string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.piles[number]
	if p != nil {
		p.Current = 0
	}
}

func (s *Service) TerminateCharge(busID, number string) soc.Verdict {
	verdict := s.soc.FullVerdict(busID, s.Current(number))
	s.TrickleStop(number)
	return verdict
}

func (s *Service) Current(number string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := s.piles[number]
	if p == nil {
		return 0
	}
	return p.Current
}

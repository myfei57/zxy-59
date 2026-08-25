package pile

func RatedPower(kind string) int {
	if kind == "fast" {
		return 120000
	}
	return 30000
}

func EstimateDuration(kind string, wattHours int) int {
	power := RatedPower(kind)
	if power <= 0 {
		return 0
	}
	return wattHours / power
}

func (s *Service) Rate(pileID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p := s.piles[pileID]
	if p == nil {
		return 0
	}
	return RatedPower(p.Kind)
}

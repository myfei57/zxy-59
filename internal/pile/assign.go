package pile

func (s *Service) Assign(busID string) (string, error) {
	percent, err := s.soc.Lookup(busID)
	if err != nil {
		return "", err
	}
	kind := "slow"
	if percent <= 40 {
		kind = "fast"
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for number, p := range s.piles {
		if p.Kind == kind && p.Owner == "" {
			p.Owner = busID
			return number, nil
		}
	}
	return "", errNoPile
}

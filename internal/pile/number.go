package pile

func (s *Service) Renumber(busID, oldNumber, newNumber string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.piles[oldNumber]
	if p == nil {
		return errUnknown
	}
	if _, exists := s.piles[newNumber]; exists {
		return errNumberUsed
	}
	delete(s.piles, oldNumber)
	p.Number = newNumber
	s.piles[newNumber] = p
	s.bus.SetMapping(busID, newNumber)
	return nil
}

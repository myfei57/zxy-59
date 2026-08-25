package bus

func (s *Service) SetDeparture(busID, slot string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.departures[busID] = slot
}

func (s *Service) Departure(busID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.departures[busID]
}

package soc

func (s *Service) AverageSoc(busIDs []string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(busIDs) == 0 {
		return 0
	}
	total := 0
	for _, id := range busIDs {
		total += s.current[id]
	}
	return total / len(busIDs)
}

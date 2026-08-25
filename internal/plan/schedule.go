package plan

func (s *Service) Schedule(busIDs []string) map[string]int {
	out := make(map[string]int, len(busIDs))
	for _, id := range busIDs {
		out[id] = s.Read(id)
	}
	return out
}

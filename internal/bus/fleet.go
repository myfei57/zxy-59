package bus

type FleetStats struct {
	Total      int `json:"total"`
	Authorized int `json:"authorized"`
	OnRoute    int `json:"on_route"`
}

func (s *Service) FleetStats() FleetStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	authorized := 0
	for id := range s.authorized {
		if s.authorized[id] {
			authorized++
		}
	}
	seen := make(map[string]bool)
	for _, ids := range s.routes {
		for _, id := range ids {
			seen[id] = true
		}
	}
	return FleetStats{
		Total:      len(s.vehicles),
		Authorized: authorized,
		OnRoute:    len(seen),
	}
}

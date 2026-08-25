package bus

type RosterEntry struct {
	BusID     string `json:"bus_id"`
	Route     string `json:"route"`
	Pile      string `json:"pile"`
	Departure string `json:"departure"`
}

func (s *Service) Roster() []RosterEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]string, 0, len(s.vehicles))
	for id := range s.vehicles {
		ids = append(ids, id)
	}
	out := make([]RosterEntry, 0, len(ids))
	for _, id := range ids {
		out = append(out, RosterEntry{
			BusID:     id,
			Route:     s.vehicles[id].Route,
			Pile:      s.mapping[id],
			Departure: s.departures[id],
		})
	}
	return out
}

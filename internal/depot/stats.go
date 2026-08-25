package depot

type Stats struct {
	Zones    int `json:"zones"`
	Piles    int `json:"piles"`
	Buses    int `json:"buses"`
	Capacity int `json:"capacity"`
}

func (s *Service) Stats() Stats {
	return Stats{
		Zones:    len(s.zones),
		Piles:    len(s.piles.List()),
		Buses:    len(s.buses.List()),
		Capacity: s.grid.Current(),
	}
}

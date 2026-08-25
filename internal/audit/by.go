package audit

func (s *Service) ByPile(pileID string) ([]Record, error) {
	records, err := s.List()
	if err != nil {
		return nil, err
	}
	out := make([]Record, 0)
	for _, rec := range records {
		if rec.PileID == pileID {
			out = append(out, rec)
		}
	}
	return out, nil
}

func (s *Service) ByBus(busID string) ([]Record, error) {
	records, err := s.List()
	if err != nil {
		return nil, err
	}
	out := make([]Record, 0)
	for _, rec := range records {
		if rec.BusID == busID {
			out = append(out, rec)
		}
	}
	return out, nil
}

func (s *Service) ByDepot(depot string) ([]Record, error) {
	records, err := s.List()
	if err != nil {
		return nil, err
	}
	out := make([]Record, 0)
	for _, rec := range records {
		if rec.Depot == depot {
			out = append(out, rec)
		}
	}
	return out, nil
}

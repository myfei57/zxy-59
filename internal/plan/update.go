package plan

func (s *Service) Update(busID, pileID string, sequence int) error {
	p := Plan{BusID: busID, PileID: pileID, Sequence: sequence}
	if err := s.store.SavePlan(busID, p); err != nil {
		return err
	}
	return s.apply(p)
}

func (s *Service) apply(p Plan) error {
	return s.piles.Switch(p.BusID, p.PileID)
}

func (s *Service) Load(busID string) (Plan, error) {
	var p Plan
	if err := s.store.LoadPlan(busID, &p); err != nil {
		return Plan{}, err
	}
	return p, nil
}

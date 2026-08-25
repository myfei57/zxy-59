package console

func (s *Server) SeedDemo() error {
	s.depot.AddZone("zone-a")
	s.depot.AddZone("zone-b")
	s.depot.AssignPileToZone("zone-a", "p1")
	s.depot.AssignPileToZone("zone-a", "p2")
	s.depot.AssignPileToZone("zone-b", "p3")

	if err := s.piles.Register("p1", "fast"); err != nil {
		return err
	}
	if err := s.piles.Register("p2", "slow"); err != nil {
		return err
	}
	if err := s.piles.Register("p3", "fast"); err != nil {
		return err
	}

	s.buses.AddVehicle("b1", "route-1")
	s.buses.AddVehicle("b2", "route-2")
	s.buses.Authorize("b1")
	s.buses.Authorize("b2")
	s.buses.SetDeparture("b1", "06:00")
	s.buses.SetDeparture("b2", "06:30")
	s.buses.SetRoute("route-1", []string{"b1"})
	s.buses.SetRoute("route-2", []string{"b2"})

	s.soc.SetCapacityWh("b1", 200000)
	s.soc.SetCapacityWh("b2", 200000)
	s.soc.SetEnergy("b1", 30000)
	s.soc.SetEnergy("b2", 180000)
	s.soc.Refresh("b1")
	s.soc.Refresh("b2")

	if _, err := s.piles.Assign("b1"); err != nil {
		return err
	}
	if _, err := s.piles.Assign("b2"); err != nil {
		return err
	}
	if err := s.piles.BeginSession("p1", "b1"); err != nil {
		return err
	}
	if err := s.piles.BeginSession("p2", "b2"); err != nil {
		return err
	}
	if err := s.piles.Start("p1", 90000); err != nil {
		return err
	}
	if err := s.piles.Start("p2", 30000); err != nil {
		return err
	}

	s.power.Meter("p1", 12.5)
	s.power.Meter("p2", 3.2)
	if _, err := s.audit.Record("depot-main", "charge", "b1", "p1", 12.5); err != nil {
		return err
	}
	if _, err := s.audit.Record("depot-main", "charge", "b2", "p2", 3.2); err != nil {
		return err
	}
	if err := s.plans.Update("b1", "p1", 1); err != nil {
		return err
	}
	if err := s.plans.Update("b2", "p2", 1); err != nil {
		return err
	}

	s.grid.AddFeeder("feeder-a", 80000)
	s.grid.AddFeeder("feeder-b", 40000)
	s.limits.SetDaily("p1", 500)
	s.limits.SetDaily("p2", 300)
	return nil
}

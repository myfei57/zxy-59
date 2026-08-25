package plan

func (s *Service) Read(busID string) int {
	s.soc.Refresh(busID)
	return s.soc.Estimate(busID)
}

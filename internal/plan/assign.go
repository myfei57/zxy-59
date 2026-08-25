package plan

func (s *Service) Assign(busID, number string) bool {
	return s.piles.Allocate(number, busID)
}

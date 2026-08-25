package depot

func (s *Service) RenumberPile(busID, oldNumber, newNumber string) error {
	return s.piles.Renumber(busID, oldNumber, newNumber)
}

func (s *Service) ExpandGrid(watts int) int {
	return s.grid.Expand(watts)
}

func (s *Service) GridCapacity() int {
	return s.grid.Current()
}

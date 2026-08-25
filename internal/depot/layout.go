package depot

func (s *Service) AssignPileToZone(zone, number string) {
	s.zonePiles[zone] = append(s.zonePiles[zone], number)
}

func (s *Service) ZonePiles(zone string) []string {
	return s.zonePiles[zone]
}

func (s *Service) ZoneCount() int {
	return len(s.zonePiles)
}

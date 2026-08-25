package pile

func (s *Service) Plug(busID, number string) error {
	if !s.bus.Authorized(busID) {
		return errUnauthorized
	}
	return s.Engage(number)
}

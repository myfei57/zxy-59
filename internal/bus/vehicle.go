package bus

func (s *Service) AddVehicle(id, route string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.vehicles[id] = Vehicle{ID: id, Route: route}
}

func (s *Service) Vehicle(id string) (Vehicle, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.vehicles[id]
	return v, ok
}

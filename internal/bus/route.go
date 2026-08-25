package bus

func (s *Service) SetRoute(route string, busIDs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes[route] = append([]string(nil), busIDs...)
}

func (s *Service) Route(route string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.routes[route]...)
}

func (s *Service) RouteCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.routes)
}

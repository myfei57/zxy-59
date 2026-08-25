package bus

func (s *Service) Authorize(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.authorized[id] = true
}

func (s *Service) Authorized(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.authorized[id]
}

package soc

func (s *Service) SetCapacityWh(busID string, wattHours int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.capacity[busID] = wattHours
	s.current[busID] = s.socFromEnergy(busID)
}

func (s *Service) CapacityWh(busID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.capacity[busID]
}

func (s *Service) SetEnergy(busID string, wattHours int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.energy[busID] = wattHours
	s.current[busID] = s.socFromEnergy(busID)
}

func (s *Service) Energy(busID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.energy[busID]
}

func (s *Service) Charge(busID string, wattHours int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.energy[busID] += wattHours
	if cap := s.capacity[busID]; cap > 0 && s.energy[busID] > cap {
		s.energy[busID] = cap
	}
	s.current[busID] = s.socFromEnergy(busID)
	return s.current[busID]
}

func (s *Service) Discharge(busID string, wattHours int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.energy[busID] -= wattHours
	if s.energy[busID] < 0 {
		s.energy[busID] = 0
	}
	s.current[busID] = s.socFromEnergy(busID)
	return s.current[busID]
}

func (s *Service) socFromEnergy(busID string) int {
	capacity := s.capacity[busID]
	if capacity <= 0 {
		return 0
	}
	percent := s.energy[busID] * 100 / capacity
	if percent > 100 {
		percent = 100
	}
	return percent
}

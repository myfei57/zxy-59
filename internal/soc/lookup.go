package soc

import "errors"

var ErrNoTelemetry = errors.New("no telemetry for bus")

func (s *Service) Lookup(busID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	percent, ok := s.current[busID]
	if !ok {
		return 0, ErrNoTelemetry
	}
	return percent, nil
}

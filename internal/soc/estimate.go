package soc

import "sync"

type Service struct {
	mu        sync.RWMutex
	current   map[string]int
	estimated map[string]int
	capacity  map[string]int
	energy    map[string]int
}

func NewService() *Service {
	return &Service{
		current:   make(map[string]int),
		estimated: make(map[string]int),
		capacity:  make(map[string]int),
		energy:    make(map[string]int),
	}
}

func (s *Service) SetCurrent(busID string, percent int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current[busID] = percent
}

func (s *Service) Current(busID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current[busID]
}

func (s *Service) Estimate(busID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.estimated[busID]
}

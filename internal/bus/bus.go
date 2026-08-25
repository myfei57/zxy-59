package bus

import "sync"

type Vehicle struct {
	ID    string `json:"id"`
	Route string `json:"route"`
}

type Service struct {
	mu         sync.RWMutex
	vehicles   map[string]Vehicle
	authorized map[string]bool
	mapping    map[string]string
	departures map[string]string
	routes     map[string][]string
}

func NewService() *Service {
	return &Service{
		vehicles:   make(map[string]Vehicle),
		authorized: make(map[string]bool),
		mapping:    make(map[string]string),
		departures: make(map[string]string),
		routes:     make(map[string][]string),
	}
}

func (s *Service) List() []Vehicle {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Vehicle, 0, len(s.vehicles))
	for _, v := range s.vehicles {
		out = append(out, v)
	}
	return out
}

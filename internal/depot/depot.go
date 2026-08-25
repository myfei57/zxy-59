package depot

import (
	"buscharge/internal/bus"
	"buscharge/internal/grid"
	"buscharge/internal/ns"
	"buscharge/internal/pile"
)

type Service struct {
	name      string
	zones     map[string]ns.Namespace
	zonePiles map[string][]string
	piles     *pile.Service
	buses     *bus.Service
	grid      *grid.Capacity
}

func NewService(name string, p *pile.Service, b *bus.Service, g *grid.Capacity) *Service {
	return &Service{
		name:      name,
		zones:     make(map[string]ns.Namespace),
		zonePiles: make(map[string][]string),
		piles:     p,
		buses:     b,
		grid:      g,
	}
}

func (s *Service) AddZone(zone string) ns.Namespace {
	n := ns.New(s.name, zone)
	s.zones[zone] = n
	return n
}

func (s *Service) Zone(zone string) (ns.Namespace, bool) {
	n, ok := s.zones[zone]
	return n, ok
}

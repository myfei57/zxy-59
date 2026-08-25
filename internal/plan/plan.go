package plan

import (
	"buscharge/internal/pile"
	"buscharge/internal/soc"
	"buscharge/internal/store"
)

type Plan struct {
	BusID    string `json:"bus_id"`
	PileID   string `json:"pile_id"`
	Sequence int    `json:"sequence"`
}

type Service struct {
	piles   *pile.Service
	soc     *soc.Service
	store   *store.Store
	windows map[string]Window
}

func NewService(p *pile.Service, s *soc.Service, st *store.Store) *Service {
	return &Service{
		piles:   p,
		soc:     s,
		store:   st,
		windows: make(map[string]Window),
	}
}

func (s *Service) WindowCount() int {
	return len(s.windows)
}

package power

import (
	"sync"

	"buscharge/internal/grid"
	"buscharge/internal/quota"
	"buscharge/internal/store"
)

type Service struct {
	grid      *grid.Capacity
	quota     *quota.Accumulator
	store     *store.Store
	reserveMu sync.RWMutex
	reserved  map[string]int
}

func NewService(g *grid.Capacity, q *quota.Accumulator, st *store.Store) *Service {
	return &Service{grid: g, quota: q, store: st, reserved: make(map[string]int)}
}

func (s *Service) Limit(requested int) int {
	capacity := s.grid.Current()
	if requested > capacity {
		return capacity
	}
	return requested
}

func (s *Service) QuotaTotal() float64 {
	return s.quota.Total()
}

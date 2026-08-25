package pile

import (
	"sync"

	"buscharge/internal/bus"
	"buscharge/internal/power"
	"buscharge/internal/soc"
	"buscharge/internal/store"
)

type Pile struct {
	Number    string `json:"number"`
	Kind      string `json:"kind"`
	Owner     string `json:"owner"`
	Current   int    `json:"current"`
	Contactor bool   `json:"contactor"`
	Engaged   bool   `json:"engaged"`
}

type Service struct {
	mu       sync.RWMutex
	piles    map[string]*Pile
	sessions map[string]Session
	soc      *soc.Service
	bus      *bus.Service
	power    *power.Service
	store    *store.Store
}

func NewService(s *soc.Service, b *bus.Service, p *power.Service, st *store.Store) *Service {
	return &Service{
		piles:    make(map[string]*Pile),
		sessions: make(map[string]Session),
		soc:      s,
		bus:      b,
		power:    p,
		store:    st,
	}
}

func (s *Service) Register(number, kind string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.piles[number]; exists {
		return errDuplicate
	}
	s.piles[number] = &Pile{Number: number, Kind: kind}
	_ = s.store.SavePile(number, *s.piles[number])
	return nil
}

func (s *Service) List() []Pile {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Pile, 0, len(s.piles))
	for _, p := range s.piles {
		out = append(out, *p)
	}
	return out
}

func (s *Service) Snapshot(number string) (Pile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.piles[number]
	if !ok {
		return Pile{}, false
	}
	return *p, true
}

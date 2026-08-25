package audit

import (
	"time"

	"github.com/google/uuid"

	"buscharge/internal/store"
)

type Service struct {
	store *store.Store
}

func NewService(st *store.Store) *Service {
	return &Service{store: st}
}

func (s *Service) Record(depot, event, busID, pileID string, amount float64) (Record, error) {
	rec := Record{
		ID:     uuid.NewString(),
		Depot:  depot,
		Event:  event,
		BusID:  busID,
		PileID: pileID,
		Amount: amount,
		At:     time.Now().UTC(),
	}
	if err := s.store.Put("audit/"+rec.ID, rec); err != nil {
		return Record{}, err
	}
	return rec, nil
}

func (s *Service) Get(id string) (Record, error) {
	var rec Record
	if err := s.store.Get("audit/"+id, &rec); err != nil {
		return Record{}, err
	}
	return rec, nil
}

package pile

import "time"

type Session struct {
	PileID  string    `json:"pile_id"`
	BusID   string    `json:"bus_id"`
	Started time.Time `json:"started"`
	Ended   time.Time `json:"ended"`
	Active  bool      `json:"active"`
}

func (s *Service) BeginSession(pileID, busID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.piles[pileID]
	if p == nil {
		return errUnknown
	}
	s.sessions[pileID] = Session{
		PileID:  pileID,
		BusID:   busID,
		Started: time.Now().UTC(),
		Active:  true,
	}
	return nil
}

func (s *Service) EndSession(pileID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[pileID]
	if !ok {
		return errUnknown
	}
	session.Ended = time.Now().UTC()
	session.Active = false
	s.sessions[pileID] = session
	return nil
}

func (s *Service) Session(pileID string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[pileID]
	return session, ok
}

func (s *Service) ActiveSessions() []Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Session, 0, len(s.sessions))
	for _, session := range s.sessions {
		if session.Active {
			out = append(out, session)
		}
	}
	return out
}

func (s *Service) ActiveCount() int {
	return len(s.ActiveSessions())
}

package console

import "net/http"

type ChargeCycleResult struct {
	BusID   string  `json:"bus_id"`
	PileID  string  `json:"pile_id"`
	Allowed int     `json:"allowed"`
	Reading float64 `json:"reading"`
}

func (s *Server) RunChargeCycle(busID string, requested int, reading float64) (ChargeCycleResult, error) {
	pileID, err := s.piles.Assign(busID)
	if err != nil {
		return ChargeCycleResult{}, err
	}
	if err := s.piles.Plug(busID, pileID); err != nil {
		return ChargeCycleResult{}, err
	}
	if err := s.piles.Start(pileID, requested); err != nil {
		return ChargeCycleResult{}, err
	}
	if err := s.piles.BeginSession(pileID, busID); err != nil {
		return ChargeCycleResult{}, err
	}
	meter := s.power.Meter(pileID, reading)
	if _, err := s.audit.Record("depot-main", "charge", busID, pileID, reading); err != nil {
		return ChargeCycleResult{}, err
	}
	snapshot, _ := s.piles.Snapshot(pileID)
	return ChargeCycleResult{
		BusID:   busID,
		PileID:  pileID,
		Allowed: snapshot.Current,
		Reading: meter.Reading,
	}, nil
}

func (s *Server) RunTermination(busID, pileID string) error {
	if err := s.piles.EndSession(pileID); err != nil {
		return err
	}
	s.piles.TerminateCharge(busID, pileID)
	return nil
}

func (s *Server) handleChargeCycle(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID     string  `json:"bus_id"`
		Requested int     `json:"requested"`
		Reading   float64 `json:"reading"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id, requested and reading are required")
		return
	}
	result, err := s.RunChargeCycle(body.BusID, body.Requested, body.Reading)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleTerminationCycle(w http.ResponseWriter, req *http.Request) {
	var body struct {
		BusID  string `json:"bus_id"`
		PileID string `json:"pile_id"`
	}
	if err := readJSON(req, &body); err != nil {
		writeError(w, http.StatusBadRequest, "bus_id and pile_id are required")
		return
	}
	if err := s.RunTermination(body.BusID, body.PileID); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"terminated": true})
}

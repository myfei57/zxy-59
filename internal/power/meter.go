package power

type MeterReading struct {
	PileID  string  `json:"pile_id"`
	Reading float64 `json:"reading"`
	Total   float64 `json:"total"`
}

func (s *Service) Meter(pileID string, reading float64) MeterReading {
	total := s.quota.Add(reading)
	rec := MeterReading{PileID: pileID, Reading: reading, Total: total}
	_ = s.store.SaveMeter(pileID, rec)
	return rec
}

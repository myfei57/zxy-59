package audit

type Summary struct {
	Records int     `json:"records"`
	Amount  float64 `json:"amount"`
}

func (s *Service) List() ([]Record, error) {
	ids := s.store.List("audit")
	out := make([]Record, 0, len(ids))
	for _, id := range ids {
		rec, err := s.Get(id)
		if err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, nil
}

func (s *Service) Summary() (Summary, error) {
	records, err := s.List()
	if err != nil {
		return Summary{}, err
	}
	summary := Summary{Records: len(records)}
	for _, rec := range records {
		summary.Amount += rec.Amount
	}
	return summary, nil
}

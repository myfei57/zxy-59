package audit

func (s *Service) Filter(event string) ([]Record, error) {
	records, err := s.List()
	if err != nil {
		return nil, err
	}
	out := make([]Record, 0)
	for _, rec := range records {
		if rec.Event == event {
			out = append(out, rec)
		}
	}
	return out, nil
}

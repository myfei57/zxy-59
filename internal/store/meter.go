package store

func (s *Store) MeterKey(id string) string {
	return "meter/" + id
}

func (s *Store) SaveMeter(id string, value any) error {
	return s.Put(s.MeterKey(id), value)
}

func (s *Store) LoadMeter(id string, out any) error {
	return s.Get(s.MeterKey(id), out)
}

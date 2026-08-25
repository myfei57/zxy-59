package store

func (s *Store) PileKey(id string) string {
	return "pile/" + id
}

func (s *Store) SavePile(id string, value any) error {
	return s.Put(s.PileKey(id), value)
}

func (s *Store) LoadPile(id string, out any) error {
	return s.Get(s.PileKey(id), out)
}

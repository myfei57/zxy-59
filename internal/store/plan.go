package store

func (s *Store) PlanKey(id string) string {
	return "plan/" + id
}

func (s *Store) SavePlan(id string, value any) error {
	return s.Put(s.PlanKey(id), value)
}

func (s *Store) LoadPlan(id string, out any) error {
	return s.Get(s.PlanKey(id), out)
}

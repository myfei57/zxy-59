package plan

import "time"

type Window struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (s *Service) SetWindow(busID string, start, end time.Time) {
	s.windows[busID] = Window{Start: start, End: end}
}

func (s *Service) Window(busID string) (Window, bool) {
	w, ok := s.windows[busID]
	return w, ok
}

func Overlaps(a, b Window) bool {
	return a.Start.Before(b.End) && b.Start.Before(a.End)
}

func (s *Service) Conflicts() []string {
	ids := make([]string, 0, len(s.windows))
	for id := range s.windows {
		ids = append(ids, id)
	}
	out := make([]string, 0)
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			if Overlaps(s.windows[ids[i]], s.windows[ids[j]]) {
				out = append(out, ids[i]+"/"+ids[j])
			}
		}
	}
	return out
}

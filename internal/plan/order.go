package plan

import "sort"

type ChargeCandidate struct {
	BusID   string `json:"bus_id"`
	Percent int    `json:"percent"`
}

func (s *Service) ChargeOrder(busIDs []string) []ChargeCandidate {
	out := make([]ChargeCandidate, 0, len(busIDs))
	for _, id := range busIDs {
		out = append(out, ChargeCandidate{BusID: id, Percent: s.Read(id)})
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Percent < out[j].Percent
	})
	return out
}

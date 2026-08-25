package console

import "net/http"

type PileInsight struct {
	Number    string `json:"number"`
	Kind      string `json:"kind"`
	Owner     string `json:"owner"`
	Occupancy int    `json:"occupancy"`
	Rate      int    `json:"rate"`
}

func (s *Server) pileInsights() []PileInsight {
	out := make([]PileInsight, 0)
	for _, p := range s.piles.List() {
		out = append(out, PileInsight{
			Number:    p.Number,
			Kind:      p.Kind,
			Owner:     p.Owner,
			Occupancy: s.buses.Occupancy(p.Number),
			Rate:      s.piles.Rate(p.Number),
		})
	}
	return out
}

func (s *Server) handleInsights(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"piles":           s.pileInsights(),
		"feeder_total":    s.grid.FeederTotal(),
		"feeder_headroom": s.grid.FeederHeadroom(),
		"feeder_count":    s.grid.FeederCount(),
		"mapped_buses":    s.buses.MappedBusCount(),
	})
}

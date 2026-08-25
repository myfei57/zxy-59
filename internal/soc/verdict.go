package soc

type Verdict struct {
	Full           bool `json:"full"`
	TrickleStopped bool `json:"trickle_stopped"`
}

func (s *Service) FullVerdict(busID string, trickleCurrent int) Verdict {
	return Verdict{
		Full:           s.Current(busID) >= 100,
		TrickleStopped: trickleCurrent == 0,
	}
}

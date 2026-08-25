package depot

import "buscharge/internal/ns"

func (s *Service) ListZones() []ns.Namespace {
	out := make([]ns.Namespace, 0, len(s.zones))
	for _, n := range s.zones {
		out = append(out, n)
	}
	return out
}

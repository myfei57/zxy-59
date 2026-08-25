package plan

import "context"

func (s *Service) ReadWithContext(busID string, ctx context.Context) (int, error) {
	s.soc.Refresh(busID)
	return s.soc.Estimate(busID), nil
}

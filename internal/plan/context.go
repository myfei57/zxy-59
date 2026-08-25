package plan

import "context"

func (s *Service) ReadWithContext(busID string, ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	s.soc.Refresh(busID)
	return s.soc.Estimate(busID), nil
}

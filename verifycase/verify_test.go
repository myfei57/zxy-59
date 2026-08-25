package verifycase

import (
	"context"
	"testing"

	"buscharge/internal/bus"
	"buscharge/internal/grid"
	"buscharge/internal/pile"
	"buscharge/internal/plan"
	"buscharge/internal/power"
	"buscharge/internal/quota"
	"buscharge/internal/soc"
	"buscharge/internal/store"
)

func TestBcPlanContextCancel(t *testing.T) {
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	s := soc.NewService()
	b := bus.NewService()
	g := grid.NewCapacity(80000)
	q := quota.NewAccumulator()
	pw := power.NewService(g, q, st)
	ps := pile.NewService(s, b, pw, st)
	plans := plan.NewService(ps, s, st)
	s.SetCurrent("b1", 70)
	s.Refresh("b1")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if percent, err := plans.ReadWithContext("b1", ctx); err == nil {
		t.Fatalf("read must respect cancellation, got %d", percent)
	}
}

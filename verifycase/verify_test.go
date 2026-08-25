package verifycase

import (
	"testing"

	"buscharge/internal/bus"
	"buscharge/internal/grid"
	"buscharge/internal/pile"
	"buscharge/internal/power"
	"buscharge/internal/quota"
	"buscharge/internal/soc"
	"buscharge/internal/store"
)

func TestBcPileCurrentDrop(t *testing.T) {
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	s := soc.NewService()
	b := bus.NewService()
	g := grid.NewCapacity(50)
	q := quota.NewAccumulator()
	pw := power.NewService(g, q, st)
	ps := pile.NewService(s, b, pw, st)
	if err := ps.Register("p1", "fast"); err != nil {
		t.Fatal(err)
	}
	allowed := ps.BeginCharge("p1", 120)
	if allowed != 50 {
		t.Fatalf("charge current must be limited to grid capacity, got %d", allowed)
	}
}

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

func TestBcSocLookupError(t *testing.T) {
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
	if err := ps.Register("p1", "fast"); err != nil {
		t.Fatal(err)
	}
	if err := ps.Register("p2", "slow"); err != nil {
		t.Fatal(err)
	}
	number, err := ps.Assign("ghost")
	if err == nil {
		t.Fatalf("assign must fail for a bus with no telemetry, got pile %s", number)
	}
	if owner := ps.Owner("p1"); owner != "" {
		t.Fatalf("no pile should be assigned without telemetry, got %s", owner)
	}
}

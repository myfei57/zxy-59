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

func TestBcPowerReserveRelease(t *testing.T) {
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	s := soc.NewService()
	b := bus.NewService()
	g := grid.NewCapacity(100)
	q := quota.NewAccumulator()
	pw := power.NewService(g, q, st)
	ps := pile.NewService(s, b, pw, st)
	if err := ps.Register("p1", "fast"); err != nil {
		t.Fatal(err)
	}
	if err := ps.Register("p2", "fast"); err != nil {
		t.Fatal(err)
	}
	if err := ps.Start("p1", 80); err != nil {
		t.Fatal(err)
	}
	if err := ps.Stop("p1"); err != nil {
		t.Fatal(err)
	}
	if err := ps.Start("p2", 80); err != nil {
		t.Fatal(err)
	}
	snap, _ := ps.Snapshot("p2")
	if snap.Current != 80 {
		t.Fatalf("second pile must get full capacity after the first stops, got %d", snap.Current)
	}
}
